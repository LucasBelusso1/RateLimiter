package middleware

import (
	"net"
	"net/http"
	"strconv"

	config "github.com/LucasBelusso1/go-ratelimiter/configs"
	"github.com/LucasBelusso1/go-ratelimiter/internal/dbstrategy"
)

func AccessLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		token := r.Header.Get("API_KEY")

		if err != nil {
			http.Error(w, "IP not found", http.StatusUnauthorized)
			return
		}

		redisStrategy := dbstrategy.NewRedisStrategy()
		dbContext := dbstrategy.NewDbContext(redisStrategy)

		config, err := config.LoadConfig(".")

		if err != nil {
			http.Error(w, "Error trying to load the server configurations.", http.StatusInternalServerError)
			return
		}

		var clientCanAccess bool

		if token != "" {
			clientCanAccess, err = validateLimitations(dbContext, token, config.TimeForToken, config.TokenLimit)
		} else {
			clientCanAccess, err = validateLimitations(dbContext, ip, config.TimeForIp, config.IpLimit)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !clientCanAccess {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func validateLimitations(dbContext *dbstrategy.DbContext, field string, timeLimit, requestLimit int) (bool, error) {
	value := dbContext.GetKey(field)

	if value == "" {
		dbContext.SetNewKeyWithTimeLimit(field, timeLimit)
		return true, nil
	}

	currentRequests, err := strconv.Atoi(value)

	if err != nil {
		return false, err
	}

	if currentRequests >= requestLimit {
		return false, nil
	} else {
		currentRequests++
		dbContext.IncrementExistingkey(field, currentRequests)
	}

	return true, nil
}
