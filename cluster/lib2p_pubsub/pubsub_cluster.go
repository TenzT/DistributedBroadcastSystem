package lib2p_pubsub

import (
	"BroadcastService/common"
	"context"
	"encoding/json"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"log"

	"github.com/libp2p/go-libp2p-pubsub"
)

const PUBSUB_TOPIC = "distributed-broadcast-system"
const DISCOVER_SERVICE_TAG = "dbs-network"

const (
	MESSAGE_TYPE_UNKNOW = iota
	MESSAGE_TYPE_DATA_SINGLE
	MESSAGE_TYPE_DATA_BATCH
)

// PubSubCluster implements cluster with libp2p GossipSub
type PubSubCluster struct {
	ctx           context.Context
	messageBuffer chan *Message
	sendBuffer    chan *common.Data
	receiveBuffer chan common.Data

	host  host.Host
	ps    *pubsub.PubSub
	topic *pubsub.Topic
	sub   *pubsub.Subscription
}

func (c *PubSubCluster) Broadcast(dataList []common.Data) {
	for _, item := range dataList {
		c.sendBuffer <- &item
	}

}

func (c *PubSubCluster) Receive() <-chan common.Data {
	return c.receiveBuffer
}

func (c *PubSubCluster) init() error {
	log.Println("Init pubsub cluster ")

	// create a new libp2p Host that listens on a random TCP port
	h, err := libp2p.New(c.ctx, libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		log.Panic(err.Error())
	}
	c.host = h

	ps, err := pubsub.NewGossipSub(c.ctx, h)
	if err != nil {
		log.Panic(err.Error())
	}
	c.ps = ps

	// setup local mDNS discovery
	err = setupDiscovery(c.ctx, h)
	if err != nil {
		log.Panic(err.Error())
	}

	// join topic
	topic, err := ps.Join(PUBSUB_TOPIC)
	if err != nil {
		log.Panic(err.Error())
	}
	c.topic = topic

	// create subscribtion
	sub, err := topic.Subscribe()
	if err != nil {
		log.Panic(err.Error())
	}
	c.sub = sub

	go c.receiveLoop()
	go c.sendLoop()

	return nil
}

func (c *PubSubCluster) sendLoop() {
	for {
		data := <-c.sendBuffer
		b, err := json.Marshal(data)
		msg := Message{
			Payload: b,
			Type:    MESSAGE_TYPE_DATA_SINGLE,
		}

		b, err = json.Marshal(msg)
		if err != nil {
			log.Println("Error marshalling data")
			continue
		}
		err = c.topic.Publish(c.ctx, b)
		if err != nil {
			log.Println("Error publishing data")
		}
	}
}

func (c *PubSubCluster) receiveLoop() {
	for {
		msg, _ := c.sub.Next(c.ctx)

		// only forward messages delivered by others
		if msg.ReceivedFrom == c.host.ID() {
			continue
		}

		cm := &Message{}
		err := json.Unmarshal(msg.Data, cm)
		if err != nil {
			log.Println("Error unmarshalling cm")
			continue
		}
		switch cm.Type {
		case MESSAGE_TYPE_DATA_SINGLE:

			data := &common.Data{}
			err := json.Unmarshal(cm.Payload, data)
			if err != nil {
				log.Println("Error unmarshalling data from message")
				continue
			}
			// send valid messages onto the Messages channel
			c.receiveBuffer <- *data
		}

	}
}

// discoveryNotifee gets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
	h host.Host
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	log.Printf("discovered new peer %s\n", pi.ID.Pretty())
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		log.Printf("error connecting to peer %s: %s\n", pi.ID.Pretty(), err)
	}
}

// setupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func setupDiscovery(ctx context.Context, h host.Host) error {
	// setup mDNS discovery to find local peers
	disc := mdns.NewMdnsService(h, DISCOVER_SERVICE_TAG)

	n := discoveryNotifee{h: h}
	disc.RegisterNotifee(&n)
	return nil
}

func NewPubSubCluster(ctx context.Context) *PubSubCluster {
	cluster := &PubSubCluster{
		ctx:           ctx,
		sendBuffer:    make(chan *common.Data, 100000),
		receiveBuffer: make(chan common.Data, 100000),
	}

	cluster.init()

	return cluster
}

type Message struct {
	Type    int
	Payload []byte
}
