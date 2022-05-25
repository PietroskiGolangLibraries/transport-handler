package worker

import (
	"sync"
)

type (
	ChanMonitor struct {
		isClosed bool
		mtx      *sync.Mutex
		tl       ThreadLocker
	}
)

func NewChanMonitor() *ChanMonitor {
	return &ChanMonitor{
		isClosed: false,
		mtx:      &sync.Mutex{},
		tl:       NewThreadLocker(),
	}
}

func (cm *ChanMonitor) IsChanClosed() bool {
	cm.tl.Lock()
	defer cm.tl.Unlock()
	return cm.isClosed
}

func (cm *ChanMonitor) SetChanToClose() {
	cm.tl.Lock()
	defer cm.tl.Unlock()
	cm.isClosed = true
}
