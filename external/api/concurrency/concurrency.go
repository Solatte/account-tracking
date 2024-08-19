package concurrency

import (
	"sync"

	"github.com/iqbalbaharum/sol-stalker/internal/generators"
)

var SubscribeWg sync.WaitGroup
var TxChannel = make(chan generators.GeyserResponse)
