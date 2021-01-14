package persister

import (
	"github.com/KyberNetwork/cache/ethereum"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"path"
	"testing"
	"time"
)

type LeveldbPersisterTestSuite struct {
	suite.Suite
	persister     *LeveldbPersister
	persisterPath string
}

func TestLeveldbPersisterTestSuite(t *testing.T) {
	suite.Run(t, new(LeveldbPersisterTestSuite))
}

func (ts *LeveldbPersisterTestSuite) SetupSuite() {
	assert := ts.Assert()

	currentDir, err := os.Getwd()
	assert.NoError(err)

	err = os.Mkdir("leveldb-storage", os.ModePerm)
	assert.NoError(err)

	persisterPath := path.Join(currentDir, "leveldb-storage")
	ts.persisterPath = persisterPath
	os.Setenv("LEVELDB_PERSISTER_PATH", persisterPath)

	persister, err := NewLeveldbPersister()
	assert.NoError(err)
	ts.persister = persister
}

func (ts *LeveldbPersisterTestSuite) TearDownSuite() {
	if err := os.RemoveAll(ts.persisterPath); err != nil {
		log.Println(err)
	}
}

func (ts *LeveldbPersisterTestSuite) TestSaveWeeklyGasPrice() {
	assert := ts.Assert()
	var (
		gasPrices = make(map[int64]ethereum.GasPrice)
	)
	gasPrices[time.Now().UnixNano()/1000] = ethereum.GasPrice{
		Fast:     "100",
		Standard: "80",
		Low:      "30",
		Default:  "80",
	}
	gasPrices[time.Now().UnixNano()/1000+1] = ethereum.GasPrice{
		Fast:     "60",
		Standard: "40",
		Low:      "20",
		Default:  "40",
	}

	err := ts.persister.SaveWeeklyGasPrice(gasPrices)
	assert.NoError(err)

	result, err := ts.persister.GetWeeklyGasPrice()
	assert.NoError(err)
	assert.NotNil(result)
	assert.Equal(2, len(result))

	// test save again, value with existed key will be updated
	err = ts.persister.SaveWeeklyGasPrice(gasPrices)
	assert.NoError(err)

	result, err = ts.persister.GetWeeklyGasPrice()
	assert.NoError(err)
	assert.NotNil(result)
	assert.Equal(2, len(result))

	// teardown
	err = ts.persister.cleanGasPriceStorage()
	assert.NoError(err)
}
