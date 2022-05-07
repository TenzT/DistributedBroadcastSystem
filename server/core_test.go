package server

import (
	"BroadcastService/common"
	"BroadcastService/eventbus"
	"BroadcastService/storage/gocache"
	"BroadcastService/validator/md5_validator"
	"context"
	"github.com/hedzr/assert"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestCore_Deliver(t *testing.T) {
	kvStore := gocache.New()
	ctx, _ := context.WithCancel(context.Background())
	validator := md5_validator.New()
	c := NewCore(kvStore, ctx, validator)

	timeStamp := time.Now().Unix()

	data := &common.Data{
		Id:        "id",
		Raw:       "raw",
		Signature: "bdd166af3a63f7be696dd17a218a6ffb",
		TimeStamp: strconv.Itoa(int(timeStamp)),
	}

	// Deliver
	err := c.Deliver(data)
	assert.NoError(t, err)

	// check storage
	dataList, err := c.GetAllData()
	assert.NoError(t, err)
	assert.EqualTrue(t, len(dataList) > 0)

	event := <-c.Events()

	assert.Equal(t, eventbus.EVENT_TYPE_BROADCAST_DATA_BATCH, event.EventType)

	// save older data
	data.TimeStamp = strconv.Itoa(int(timeStamp - 10000000))

	err = c.Deliver(data)
	assert.NoError(t, err)
	dataList, err = c.GetAllData()
	assert.NoError(t, err)
	assert.EqualTrue(t, len(dataList) > 0)

	cacheTimestamp, err := strconv.Atoi(dataList[0].TimeStamp)
	assert.NoError(t, err)

	assert.EqualTrue(t, cacheTimestamp < int(timeStamp))
}

func TestCore_Run(t *testing.T) {
	kvStore := gocache.New()
	ctx, _ := context.WithCancel(context.Background())
	validator := md5_validator.New()
	c := NewCore(kvStore, ctx, validator)

	timeStamp := time.Now().Unix()

	data := &common.Data{
		Id:        "id",
		Raw:       "raw",
		Signature: "bdd166af3a63f7be696dd17a218a6ffb",
		TimeStamp: strconv.Itoa(int(timeStamp)),
	}

	// Deliver
	err := c.Deliver(data)
	assert.NoError(t, err)
	<-c.Events()

	data.Id = "id2"
	err = c.Deliver(data)
	assert.NoError(t, err)
	<-c.Events()

	go c.Run()

	for i := 0; i < 5; i++ {
		event := <-c.Events()
		l, _ := event.Payload.([]common.Data)
		log.Printf("%d: event_type:{%d}, payload{%v}", i, event.EventType, l)
		assert.Equal(t, 2, len(l))
		assert.Equal(t, "raw", l[0].Raw)
		assert.Equal(t, "bdd166af3a63f7be696dd17a218a6ffb", l[1].Signature)
	}

}
