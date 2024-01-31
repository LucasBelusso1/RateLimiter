package limiter

import (
	"strconv"

	"github.com/LucasBelusso1/go-ratelimiter/internal/dbstrategy"
)

type Limiter struct {
	DbContext    *dbstrategy.DbContext
	Field        string
	TimeLimit    int
	RequestLimit int
}

func NewLimiter(Db *dbstrategy.DbContext, field string, timeLimit, requestLimit int) *Limiter {
	return &Limiter{DbContext: Db, Field: field, TimeLimit: timeLimit, RequestLimit: requestLimit}
}

func (l *Limiter) ValidateLimit() (bool, error) {
	value := l.DbContext.GetKey(l.Field)

	if value == "" {
		l.DbContext.SetNewKeyWithTimeLimit(l.Field, l.TimeLimit)
		return true, nil
	}

	currentRequests, err := strconv.Atoi(value)

	if err != nil {
		return false, err
	}

	if currentRequests >= l.RequestLimit {
		return false, nil
	} else {
		currentRequests++
		l.DbContext.IncrementExistingkey(l.Field, currentRequests)
	}

	return true, nil
}
