package echo

import (
	"BroadcastService/common"
	"github.com/hedzr/assert"
	"testing"
	"time"
)

func TestEchoCluster_Broadcast(t *testing.T) {
	cluster := NewEchoCluster()

	data := common.Data{
		Id:        "id",
		Raw:       "raw",
		Signature: "bdd166af3a63f7be696dd17a218a6ffb",
		TimeStamp: "1651939872",
	}

	go func() {
		time.Sleep(time.Second)
		cluster.Broadcast([]common.Data{data})
	}()

	receivedData := <-cluster.Receive()

	assert.Equal(t, "id", receivedData.Id)
	assert.Equal(t, "raw", receivedData.Raw)
}
