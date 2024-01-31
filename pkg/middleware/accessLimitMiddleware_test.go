package middleware

import (
	"os"
	"path/filepath"
	"testing"

	config "github.com/LucasBelusso1/go-ratelimiter/configs"
	"github.com/LucasBelusso1/go-ratelimiter/internal/dbstrategy"
	"github.com/LucasBelusso1/go-ratelimiter/pkg/limiter"
	"github.com/stretchr/testify/suite"
)

var (
	middleware *Middleware
	cfg        *config.Conf
	context    *dbstrategy.DbContext
	limiters   []limiter.Limiter
)

type MiddlewareTestSuite struct {
	suite.Suite
}

func (suite *MiddlewareTestSuite) SetupSuite() {
	cmdDir, err := filepath.Abs("../../cmd")

	suite.NoError(err, "Error while trying to get the absolute file path to /cmd.")

	err = os.Chdir(cmdDir)

	suite.NoError(err, "Error while trying to change the file path.")

	cfg, err = config.LoadConfig(".")

	suite.NoError(err, "Error returned trying to open configuration file.")

	strategy := dbstrategy.NewRedisStrategy(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisPort)
	context = dbstrategy.NewDbContext(strategy)

	limiters = []limiter.Limiter{
		{
			DbContext:    context,
			Field:        cfg.TokenAName,
			TimeLimit:    cfg.TimeForTokenA,
			RequestLimit: cfg.TokenALimit,
		},
		{
			DbContext:    context,
			Field:        cfg.TokenBName,
			TimeLimit:    cfg.TimeForTokenB,
			RequestLimit: cfg.TokenBLimit,
		},
	}
}

func (suite *MiddlewareTestSuite) TestNewMiddleware() {
	middleware = NewMiddleware(context, limiters, cfg.IpLimit, cfg.TimeForIp)
	suite.IsType(&Middleware{}, middleware, "Wrong instace of middleware")
}

func (suite *MiddlewareTestSuite) TearDownSuite() {
	cmdDir, err := filepath.Abs("../../cmd")

	suite.NoError(err, "Error trying to get the abolute path to /cmd.")

	err = os.Chdir(filepath.Dir(cmdDir))

	suite.NoError(err, "Error trying to change the path.")
}

func Run(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}
