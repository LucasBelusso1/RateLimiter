package middleware

import (
	"net"
	"net/http"

	"github.com/LucasBelusso1/go-ratelimiter/internal/dbstrategy"
)

type Middleware struct {
	Context       *dbstrategy.DbContext
	Limiters      []Limiter
	IpLimit       int
	IpTimeBlocked int
}

func NewMiddleware(context *dbstrategy.DbContext, limiters []Limiter, ipLimit, ipTimeBlocked int) *Middleware {
	return &Middleware{Context: context, Limiters: limiters, IpLimit: ipLimit, IpTimeBlocked: ipTimeBlocked}
}

func (m *Middleware) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var clientCanAccess bool
		var err error
		var limiter Limiter
		var validateToken bool

		limiter, validateToken = m.ValidateToken(r)

		if validateToken {
			clientCanAccess, err = limiter.validateLimit()
		} else {
			var ip string

			ip, _, err = net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "IP not found", http.StatusUnauthorized)
				return
			}

			limiter := NewLimiter(m.Context, ip, m.IpTimeBlocked, m.IpLimit)
			clientCanAccess, err = limiter.validateLimit()
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

func (m *Middleware) ValidateToken(r *http.Request) (Limiter, bool) {
	token := r.Header.Get("API_KEY")
	for _, limiter := range m.Limiters {
		if limiter.Field == token {
			return limiter, true
		}
	}

	return Limiter{}, false
}
