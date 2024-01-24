package dbstrategy

type DbStrategy interface {
	GetKey(field string) string
	SetNewKeyWithTimeLimit(field string, timeLimit int) error
	IncrementExistingkey(field string, currentRequests int) error
}
