package deprecated_worker

import (
	"sync"

	handlers_model "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/models/handlers"
)

type handler struct {
	srvMap *sync.Map
}

func (h *handler) storeInSyncMap(servers map[string]handlers_model.Server) {
	for srvName, srv := range servers {
		h.srvMap.Store(srvName, srv)
	}
}

func (h *handler) loadFromSyncMap(srvName string) handlers_model.Server {
	if rawSrv, loaded := h.srvMap.Load(srvName); loaded {
		if srv, ok := rawSrv.(handlers_model.Server); ok {
			return srv
		}
	}

	return nil
}

func (h *handler) deleteFromSyncMap(srvName string) {
	h.srvMap.Delete(srvName)
}

func (h *handler) rangeOverSyncMap(f func(key any, value any) bool) {
	h.srvMap.Range(f)
}

var SyncMapRanger = func(key any, value any) bool {
	if srv, ok := value.(handlers_model.Server); ok {
		srv.Stop()
	}

	return true
}
