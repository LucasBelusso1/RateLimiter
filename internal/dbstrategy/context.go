package dbstrategy

type DbContext struct {
	strategy DbStrategy
}

func NewDbContext(strategy DbStrategy) *DbContext {
	return &DbContext{strategy: strategy}
}

func (dbc *DbContext) GetKey(field string) string {
	return dbc.strategy.GetKey(field)
}

func (dbc *DbContext) SetNewKeyWithTimeLimit(field string, timeLimit int) error {
	return dbc.strategy.SetNewKeyWithTimeLimit(field, timeLimit)
}

func (dbc *DbContext) IncrementExistingkey(field string, currentRequests int) error {
	return dbc.strategy.IncrementExistingkey(field, currentRequests)
}
