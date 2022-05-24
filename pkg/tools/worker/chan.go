package worker

import (
	"sync"
)

type (
	ChanMonitor struct {
		isClosed bool
		//done     int64
		//old      int64
		//new      int64
		mtx       *sync.Mutex
		elemCount int
	}
)

func NewChanMonitor() *ChanMonitor {
	return &ChanMonitor{
		mtx: &sync.Mutex{},
	}
}

func (cm *ChanMonitor) CountElem() {
	cm.mtx.Lock()
	defer cm.mtx.Unlock()
	cm.elemCount++
}

func (cm *ChanMonitor) PopElem() {
	cm.mtx.Lock()
	defer cm.mtx.Unlock()
	cm.elemCount--
}

func (cm *ChanMonitor) ElemNum() int {
	cm.mtx.Lock()
	defer cm.mtx.Unlock()
	return cm.elemCount
}

func (cm *ChanMonitor) IsChanClosed() bool {
	cm.mtx.Lock()
	defer cm.mtx.Unlock()
	return cm.isClosed
}

func (cm *ChanMonitor) SetChanToClose() {
	cm.mtx.Lock()
	defer cm.mtx.Unlock()
	cm.isClosed = true
}
