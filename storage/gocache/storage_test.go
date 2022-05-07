package gocache

import (
	"github.com/hedzr/assert"
	"log"
	"testing"
)

func TestGoCacheKVStorage_GetData(t *testing.T) {
	s := New()

	key := "testId"
	testData := &GoCacaheData{
		id:        key,
		raw:       "raw",
		signature: "signature",
		timestamp: "time",
	}

	_, err := s.GetData(key)

	// assertion
	assert.Error(t, err)
	assert.Equal(t, "Data not exists.", err.Error())

	list, err := s.GetAllData()
	assert.Equal(t, 0, len(list))

	err = s.SaveData(testData)
	assert.NoError(t, err)

	d, err := s.GetData(key)
	assert.NoError(t, err)
	log.Print(d)
	assert.Equal(t, key, d.GetId())
	assert.Equal(t, "raw", d.GetRawData())
	assert.Equal(t, "signature", d.GetSignature())
	assert.Equal(t, "time", d.GetTimeStamp())

	testData.id = "newTestId"
	err = s.SaveData(testData)
	assert.NoError(t, err)

	list, err = s.GetAllData()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(list))
}
