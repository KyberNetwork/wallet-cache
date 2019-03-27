package persister

import (
	"encoding/json"
	"log"

	"github.com/KyberNetwork/server-go/ethereum"
	"github.com/boltdb/bolt"
)

const (
	path   = "./persister/db/market.db"
	bucket = "market_info"
)

// BoltStorage storage for cache
type BoltStorage struct {
	marketDB *bolt.DB
}

// NewBoltStorage make bolt instance
func NewBoltStorage() (*BoltStorage, error) {
	marketDB, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	err = marketDB.Update(func(tx *bolt.Tx) error {
		if _, cErr := tx.CreateBucketIfNotExists([]byte(bucket)); cErr != nil {
			return cErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &BoltStorage{
		marketDB: marketDB,
	}, nil
}

// StoreGeneralInfo store market info
func (bs *BoltStorage) StoreGeneralInfo(mapInfo map[string]*ethereum.TokenGeneralInfo) error {
	var err error
	err = bs.marketDB.Update(func(tx *bolt.Tx) error {
		var errS error
		b := tx.Bucket([]byte(bucket))
		for k, v := range mapInfo {
			var dataJSON []byte
			dataJSON, errS = json.Marshal(*v)
			if errS != nil {
				return errS
			}
			errS = b.Put([]byte(k), dataJSON)
			if errS != nil {
				return errS
			}
		}
		return nil
	})
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// GetGeneralInfo store market info
func (bs *BoltStorage) GetGeneralInfo(mapToken map[string]ethereum.Token) (map[string]*ethereum.TokenGeneralInfo, error) {
	var err error
	result := make(map[string]*ethereum.TokenGeneralInfo)
	err = bs.marketDB.View(func(tx *bolt.Tx) error {
		var errV error
		b := tx.Bucket([]byte(bucket))
		if errV = b.ForEach(func(k, v []byte) error {
			var tokenInfo ethereum.TokenGeneralInfo
			errLoop := json.Unmarshal(v, &tokenInfo)
			if errLoop != nil {
				return errLoop
			}
			result[string(k)] = &tokenInfo
			return nil
		}); errV != nil {
			log.Println(errV.Error())
			return errV
		}
		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}
