package sub

import (
	"sync"

	pb "github.com/rpcpool/yellowstone-grpc/examples/golang/proto"
)

const (
	SUCCESS = "success"
	FAILED  = "failed"
)

type Subscription struct {
	Stream pb.Geyser_SubscribeClient
	Done   chan int
}

type SubscriptionManager struct {
	Mu            sync.Mutex
	Subscriptions map[string]*Subscription
}

func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		Subscriptions: make(map[string]*Subscription),
	}
}

func (sm *SubscriptionManager) Add(id string, stream pb.Geyser_SubscribeClient, failed bool) {
	sm.Mu.Lock()
	defer sm.Mu.Unlock()

	sm.Subscriptions[id] = &Subscription{
		Stream: stream,
		Done:   make(chan int),
	}
}

func (sm *SubscriptionManager) Remove(id string) {
	sm.Mu.Lock()
	defer sm.Mu.Unlock()

	if sub, exists := sm.Subscriptions[id]; exists {
		close(sub.Done)
	}

}

var TxManager map[string]*SubscriptionManager

func init() {

	TxManager = make(map[string]*SubscriptionManager)
	TxManager["success"] = NewSubscriptionManager()
	TxManager["failed"] = NewSubscriptionManager()
}
