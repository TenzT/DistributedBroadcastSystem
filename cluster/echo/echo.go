package echo

import (
	"BroadcastService/common"
	"log"
)

type EchoCluster struct {
	buffer chan common.Data
}

func (c *EchoCluster) Broadcast(dataList []common.Data) {
	for _, item := range dataList {
		log.Printf("receive: {%v}", item)
		c.buffer <- item
	}
}

func (c *EchoCluster) Receive() <-chan common.Data {
	return c.buffer
}

func NewEchoCluster() *EchoCluster {
	return &EchoCluster{
		buffer: make(chan common.Data, 1000),
	}
}
