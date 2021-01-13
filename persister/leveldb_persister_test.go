package persister

import (
	"github.com/KyberNetwork/cache/ethereum"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"path"
	"testing"
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

func (ts *LeveldbPersisterTestSuite) TestSaveGasPrice() {
	assert := ts.Assert()
	gasOracle := ethereum.GasPrice{
		Fast:     "100",
		Standard: "80",
		Low:      "30",
		Default:  "80",
	}

	err := ts.persister.SaveGasPrice(&gasOracle)
	assert.NoError(err)

	gasPrices, err := ts.persister.GetPassWeekGasPrice()
	assert.NoError(err)
	assert.NotNil(gasPrices)
	assert.Equal(1, len(gasPrices))

	log.Println(gasPrices)
	// teardown
	err = ts.persister.cleanGasPriceStorage()
	assert.NoError(err)
}

func (ts *LeveldbPersisterTestSuite) TestGetWeeklyAverageGasPrice() {
	assert := ts.Assert()
	var err error
	err = ts.persister.SaveGasPrice(&ethereum.GasPrice{
		Fast:     "30",
		Standard: "20",
		Low:      "10",
		Default:  "20",
	})
	assert.NoError(err)
	err = ts.persister.SaveGasPrice(&ethereum.GasPrice{
		Fast:     "100",
		Standard: "80",
		Low:      "60",
		Default:  "80",
	})
	assert.NoError(err)

	averageGas, err := ts.persister.GetWeeklyAverageGasPrice()
	assert.NoError(err)
	assert.Equal(float64(50), averageGas)
}
