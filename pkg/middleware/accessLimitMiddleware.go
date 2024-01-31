package middleware

import (
	"net"
	"net/http"

	"github.com/LucasBelusso1/go-ratelimiter/internal/dbstrategy"
	"github.com/LucasBelusso1/go-ratelimiter/pkg/limiter"
)

type Middleware struct {
	Context       *dbstrategy.DbContext
	Limiters      []limiter.Limiter
	IpLimit       int
	IpTimeBlocked int
}

func NewMiddleware(context *dbstrategy.DbContext, limiters []limiter.Limiter, ipLimit, ipTimeBlocked int) *Middleware {
	return &Middleware{Context: context, Limiters: limiters, IpLimit: ipLimit, IpTimeBlocked: ipTimeBlocked}
}

func (m *Middleware) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var clientCanAccess bool
		var err error
		var limitValidator limiter.Limiter
		var validateToken bool

		limitValidator, validateToken = m.ValidateToken(r)

		if validateToken {
			clientCanAccess, err = limitValidator.ValidateLimit()
		} else {
			var ip string

			ip, _, err = net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "IP not found", http.StatusUnauthorized)
				return
			}

			limitValidator := limiter.NewLimiter(m.Context, ip, m.IpTimeBlocked, m.IpLimit)
			clientCanAccess, err = limitValidator.ValidateLimit()
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

func (m *Middleware) ValidateToken(r *http.Request) (limiter.Limiter, bool) {
	token := r.Header.Get("API_KEY")
	for _, limiter := range m.Limiters {
		if limiter.Field == token {
			return limiter, true
		}
	}

	return limiter.Limiter{}, false
}
