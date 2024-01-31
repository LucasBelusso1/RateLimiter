package limiter

import (
	"os"
	"path/filepath"
	"testing"

	config "github.com/LucasBelusso1/go-ratelimiter/configs"
	"github.com/LucasBelusso1/go-ratelimiter/internal/dbstrategy"
	"github.com/stretchr/testify/suite"
)

type LimiterTestSuite struct {
	suite.Suite

	cfg          *config.Conf
	context      *dbstrategy.DbContext
	strategy     *dbstrategy.RedisStrategy
	limiters     []Limiter
	limitersData []mockLimiterValues
}

type mockLimiterValues struct {
	field          string
	fieldLimit     int
	fieldTimeLimit int
}

func (suite *LimiterTestSuite) SetupSuite() {
	cmdDir, err := filepath.Abs("../../cmd")

	suite.NoError(err, "Error while trying to get the absolute file path to /cmd.")

	err = os.Chdir(cmdDir)

	suite.NoError(err, "Error while trying to change the file path.")

	suite.cfg, err = config.LoadConfig(".")

	suite.NoError(err, "Error returned trying to open configuration file.")

	suite.strategy = dbstrategy.NewRedisStrategy(suite.cfg.RedisAddress, suite.cfg.RedisPassword, suite.cfg.RedisPort)
	suite.context = dbstrategy.NewDbContext(suite.strategy)

	suite.limitersData = []mockLimiterValues{
		{field: "token", fieldLimit: 10, fieldTimeLimit: 100},
		{field: "ip", fieldLimit: 100, fieldTimeLimit: 1000},
	}
}

func (suite *LimiterTestSuite) TestNewLimiter() {
	dataLimiter1 := suite.limitersData[0]
	dataLimiter2 := suite.limitersData[1]

	limiter1 := NewLimiter(suite.context, dataLimiter1.field, dataLimiter1.fieldTimeLimit, dataLimiter2.fieldLimit)
	suite.IsType(&Limiter{}, limiter1, "Wrong limiter type returned from constructor function of limiter.")
	suite.limiters = append(suite.limiters, *limiter1)

	limiter2 := NewLimiter(suite.context, dataLimiter2.field, dataLimiter2.fieldTimeLimit, dataLimiter2.fieldLimit)
	suite.IsType(&Limiter{}, limiter2, "Wrong limiter type returned from constructor function of limiter.")
	suite.limiters = append(suite.limiters, *limiter2)
}

func (suite *LimiterTestSuite) TestValidateLimit() {
	for _, limiter := range suite.limiters {
		for i := 1; i <= limiter.RequestLimit*100; i++ {
			result, err := limiter.ValidateLimit()

			if i > limiter.RequestLimit {
				suite.False(result, "Rate limit validation failed. Success when limit reached expected")
				suite.NoError(err, "An error ocurred in the rate limit validation.")
			} else {
				suite.True(result, "Rate limit validation failed. Error when the limit was not reached.")
				suite.NoError(err, "An error ocurred in the rate limit validation.")
			}
		}
	}
}

func (suite *LimiterTestSuite) TearDownSuite() {
	cmdDir, err := filepath.Abs("../../cmd")

	suite.NoError(err, "Error trying to get the abolute path to /cmd.")

	err = os.Chdir(filepath.Dir(cmdDir))

	suite.NoError(err, "Error trying to change the path.")

	suite.strategy.Client.FlushAll(suite.strategy.Ctx)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(LimiterTestSuite))
}
