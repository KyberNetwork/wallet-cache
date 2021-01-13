package persister

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/KyberNetwork/cache/ethereum"
	"github.com/syndtr/goleveldb/leveldb"
	leveldbUtil "github.com/syndtr/goleveldb/leveldb/util"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	cleanInterval = 6 * time.Hour
)

type LeveldbPersister struct {
	db *leveldb.DB
}

func NewLeveldbPersister() (*LeveldbPersister, error) {
	leveldbPath := os.Getenv("LEVELDB_PERSISTER_PATH")
	db, err := leveldb.OpenFile(leveldbPath, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	persister := &LeveldbPersister{
		db: db,
	}

	go persister.intervalCleanGasPriceStorage()
	return persister, nil
}

func (p *LeveldbPersister) intervalCleanGasPriceStorage() {
	ticker := time.NewTicker(cleanInterval)
	defer ticker.Stop()
	for {
		p.cleanGasPriceStorage()
		<-ticker.C
	}
}

func (p *LeveldbPersister) SaveGasPrice(gasOracle ethereum.GasPrice) error {
	key := fmt.Sprintf("gas_oracle_%v", time.Now().UnixNano()/1000)

	encodedKey, err := Encode(key)
	if err != nil {
		log.Println(err)
		return err
	}

	record, err := toGasOracleModel(gasOracle)
	if err != nil {
		log.Println(err)
		return err
	}

	encodedValue, err := Encode(record)
	if err != nil {
		log.Println(err)
		return err
	}
	if err := p.db.Put(encodedKey, encodedValue, nil); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (p *LeveldbPersister) cleanGasPriceStorage() error {
	now := time.Now()
	end := now.Add(-168 * time.Hour) // pass 1 week

	key := fmt.Sprintf("gas_oracle_%v", end)
	encodedKey, err := Encode(key)
	if err != nil {
		log.Println(err)
		return err
	}

	iterator := p.db.NewIterator(&leveldbUtil.Range{
		Limit: encodedKey,
	}, nil)
	defer func() {
		iterator.Release()
		if err := iterator.Error(); err != nil {
			log.Println(err)
		}
	}()

	batch := new(leveldb.Batch)
	for iterator.Next() {
		batch.Delete(iterator.Key())
	}

	if err := p.db.Write(batch, nil); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (p *LeveldbPersister) GetPassWeekGasPrice() ([]GasOracleModel, error) {
	now := time.Now()
	start := now.Add(-168 * time.Hour) // pass 1 week

	key := fmt.Sprintf("gas_oracle_")
	encodedKey, err := Encode(key)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	iterator := p.db.NewIterator(&leveldbUtil.Range{
		Start: encodedKey,
	}, nil)
	defer func() {
		iterator.Release()
		if err := iterator.Error(); err != nil {
			log.Println(err)
		}
	}()

	var result = make([]GasOracleModel, 0)
	for iterator.Next() {
		var decodedKey string
		if err := Decode(iterator.Key(), &decodedKey); err != nil {
			log.Println(err)
			return nil, err
		}

		timestamp, err := parseTimestamp(strings.TrimPrefix(decodedKey, "gas_oracle_"))
		if err != nil {
			log.Println(err)
			return nil, err
		}
		if timestamp.Before(start) {
			continue
		}

		var gasOracle GasOracleModel
		if err := Decode(iterator.Value(), &gasOracle); err != nil {
			log.Println(err)
			return nil, err
		}

		result = append(result, gasOracle)
	}
	return result, nil
}

func (p *LeveldbPersister) GetWeeklyAverageGasPrice() (float64, error) {
	gasPrices, err := p.GetPassWeekGasPrice()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	var (
		sum float64
	)
	if len(gasPrices) == 0 {
		return 0, nil
	}
	for _, gasPrice := range gasPrices {
		sum += (gasPrice.Fast + gasPrice.Standard + gasPrice.Low) / float64(3)
	}
	average := sum / float64(len(gasPrices))
	return math.Floor(average), nil
}

func Encode(data interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func Decode(in []byte, out interface{}) error {
	buffer := bytes.NewBuffer(in)
	dec := gob.NewDecoder(buffer)
	return dec.Decode(out)
}

func parseTimestamp(timeStr string) (time.Time, error) {
	i, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return time.Unix(0, 0), err
	}

	return time.Unix(0, i*1000), nil
}

type GasOracleModel struct {
	Fast     float64
	Standard float64
	Low      float64
	Default  float64
}

func (m *GasOracleModel) ToEthereumGasPrice() ethereum.GasPrice {
	return ethereum.GasPrice{
		Fast:     strconv.FormatFloat(m.Fast, 'f', -1, 64),
		Standard: strconv.FormatFloat(m.Standard, 'f', -1, 64),
		Low:      strconv.FormatFloat(m.Low, 'f', -1, 64),
		Default:  strconv.FormatFloat(m.Default, 'f', -1, 64),
	}
}

func toGasOracleModel(gasOracleRaw ethereum.GasPrice) (GasOracleModel, error) {
	fast, err := strconv.ParseFloat(gasOracleRaw.Fast, 64)
	if err != nil {
		return GasOracleModel{}, errors.New(fmt.Sprintf("error parse `fast` gas price: %v", err.Error()))
	}
	standard, err := strconv.ParseFloat(gasOracleRaw.Standard, 64)
	if err != nil {
		return GasOracleModel{}, errors.New(fmt.Sprintf("error parse `standard` gas price: %v", err.Error()))
	}
	low, err := strconv.ParseFloat(gasOracleRaw.Low, 64)
	if err != nil {
		return GasOracleModel{}, errors.New(fmt.Sprintf("error parse `low` gas price: %v", err.Error()))
	}
	defaultGas, err := strconv.ParseFloat(gasOracleRaw.Default, 64)
	if err != nil {
		return GasOracleModel{}, errors.New(fmt.Sprintf("error parse `default` gas price: %v", err.Error()))
	}
	return GasOracleModel{
		Fast:     fast,
		Standard: standard,
		Low:      low,
		Default:  defaultGas,
	}, nil
}
