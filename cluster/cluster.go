package cluster

import (
	"BroadcastService/common"
)

// Cluster responsible for nodes management in cluster and data transmission
type Cluster interface {
	// Broadcast broadcasts message to cluster
	Broadcast(dataList []common.Data)

	// Receive receives messgage from cluster
	Receive() <-chan common.Data
}
