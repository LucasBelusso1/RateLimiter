package main

import (
	"log"
	"net/http"

	config "github.com/LucasBelusso1/go-ratelimiter/configs"
	"github.com/LucasBelusso1/go-ratelimiter/internal/dbstrategy"
	"github.com/LucasBelusso1/go-ratelimiter/pkg/limiter"
	"github.com/LucasBelusso1/go-ratelimiter/pkg/middleware"
)

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatalln("Couldn't read configurations", err)
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	middleware := setNewRateLimiter(config)
	err = http.ListenAndServe(":8080", middleware.RateLimitMiddleware(mux))

	if err != nil {
		log.Fatalf("Could not listen on :8080: %v\n", err)
	}

	log.Println("Listening on port :8080")
}

func setNewRateLimiter(config *config.Conf) *middleware.Middleware {
	strategy := dbstrategy.NewRedisStrategy(config.RedisAddress, config.RedisPassword, config.RedisPort)
	context := dbstrategy.NewDbContext(strategy)

	limiters := []limiter.Limiter{
		{
			DbContext:    context,
			Field:        config.TokenAName,
			TimeLimit:    config.TimeForTokenA,
			RequestLimit: config.TokenALimit,
		},
		{
			DbContext:    context,
			Field:        config.TokenBName,
			TimeLimit:    config.TimeForTokenB,
			RequestLimit: config.TokenBLimit,
		},
	}

	return middleware.NewMiddleware(context, limiters, config.IpLimit, config.TimeForIp)
}
