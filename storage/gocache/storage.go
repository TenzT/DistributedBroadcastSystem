package gocache

import (
	"BroadcastService/storage"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"time"
)

const EXPIRED_TIME_HOUR = 24
const PURGE_INTERVAL_MIN = 5

// GoCache implementation of Storage
type GoCacheKVStorage struct {
	kvcache *cache.Cache
}

func (s *GoCacheKVStorage) CheckExist(key string) (bool, error) {
	_, found := s.kvcache.Get(key)
	return found, nil
}

func (s *GoCacheKVStorage) GetData(key string) (data storage.StoredData, err error) {
	d, found := s.kvcache.Get(key)

	if !found {
		return nil, errors.New("Data not exists.")
	}

	transfered, ok := d.(*GoCacaheData)
	if !ok {
		return nil, errors.New("Error parsing results from go cache")
	}
	return transfered, nil
}

func (s *GoCacheKVStorage) SaveData(data storage.StoredData) error {
	d := &GoCacaheData{
		id: data.GetId(),
		raw: data.GetRawData(),
		signature: data.GetSignature(),
		timestamp: data.GetTimeStamp(),
	}
	return s.kvcache.Add(data.GetId(), d, EXPIRED_TIME_HOUR * time.Hour)
}

func (s *GoCacheKVStorage) GetAllData() (dataList []storage.StoredData, err error) {
	itemMap := s.kvcache.Items()
	list := make([]storage.StoredData, 0)
	for _, v := range itemMap {
		transfered, ok := v.Object.(*GoCacaheData)
		if !ok {
			return nil, errors.New("Error parsing results from go cache")
		}
		list = append(list, transfered)
	}

	return list, nil
}


func New() storage.Storage {
	return &GoCacheKVStorage{
		kvcache: cache.New(EXPIRED_TIME_HOUR * time.Hour, PURGE_INTERVAL_MIN * time.Minute),
	}
}
