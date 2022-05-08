package server

import (
	"BroadcastService/common"
	"BroadcastService/eventbus"
	"BroadcastService/storage"
	"BroadcastService/validator"
	"context"
	"github.com/pkg/errors"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

const AUTO_BROADCAST_INTERVAL_SECOND = 5

// Core handles new data from outside and sync data to cluster
type Core struct {
	kvstore   storage.Storage
	events    chan eventbus.Event
	validator validator.Validator
	ctx       context.Context
	rwMutex   sync.RWMutex
}

// Deliver trys to save data from outside
func (c *Core) Deliver(data *common.Data) (err error) {
	log.Printf("core receives data: {%v}", data)
	isValid := c.validator.Validate(data)

	if !isValid {
		return errors.New("Invalid data")
	}

	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	dataInCache, err := c.kvstore.GetData(data.Id)
	if err != nil {
		if strings.Contains(err.Error(), "not exists.") {
			err = c.kvstore.SaveData(data)
			if err != nil {
				return nil
			}
			// try broadcast new data
			go c.tryBroadcastData(*data)
			return
		} else {
			return err
		}
	}

	// try save older data
	inputTime, err := strconv.Atoi(data.TimeStamp)
	if err != nil {
		return err
	}
	cachedTime, err := strconv.Atoi(dataInCache.GetTimeStamp())
	if err != nil {
		return err
	}

	if inputTime < cachedTime {
		return c.kvstore.SaveData(data)
	} else {
		return errors.New("Data exists")
	}
}

func (c *Core) tryBroadcastData(data common.Data) {
	event := eventbus.Event{
		EventType: eventbus.EVENT_TYPE_BROADCAST_DATA_BATCH,
		Payload:   []common.Data{data},
	}

	c.events <- event
}

// GetAllData fetch all data
func (c *Core) GetAllData() (dataList []*common.Data, err error) {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()
	list, err := c.kvstore.GetAllData()
	if err != nil {
		return nil, err
	}

	l := make([]*common.Data, 0)
	for _, data := range list {
		l = append(l, &common.Data{
			Id:        data.GetId(),
			Raw:       data.GetRawData(),
			Signature: data.GetSignature(),
			TimeStamp: data.GetTimeStamp(),
		})
	}

	return l, nil
}

// Events shots events generated inside core
func (c *Core) Events() <-chan eventbus.Event {
	return c.events
}

func (c *Core) autoBroadcast() {
	l, err := c.GetAllData()
	if err != nil {
		log.Println("Error getting data list in auto broadcast.")
		return
	}

	if len(l) == 0 {
		return
	}

	list := make([]common.Data, 0)
	for _, item := range l {
		list = append(list, *item)
	}

	event := eventbus.Event{
		EventType: eventbus.EVENT_TYPE_BROADCAST_DATA_BATCH,
		Payload:   list,
	}

	c.events <- event
}

// Start main loop
func (c *Core) Run() {
	ticker := time.NewTicker(AUTO_BROADCAST_INTERVAL_SECOND * time.Second)
	for {
		select {
		case <-ticker.C:
			log.Println("auto broadcast periodically")
			c.autoBroadcast()
		case <-c.ctx.Done():
			log.Println("Stop core")
			return
		}
	}
}

func NewCore(kvstore storage.Storage, ctx context.Context, validator validator.Validator) *Core {
	return &Core{
		kvstore:   kvstore,
		ctx:       ctx,
		validator: validator,
		events:    make(chan eventbus.Event, 10000),
	}
}
