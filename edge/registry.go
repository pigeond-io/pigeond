package edge

import (
	"github.com/pigeond-io/pigeond/common/docid"
	"sync"
)

var (
	globalRegistry *WsRegistry
	onceInit       sync.Once
)

type WsRegistry struct {
	index docid.EdgeSet
}

func GetWsRegistry() *WsRegistry {
	onceInit.Do(func() {
		globalRegistry = &WsRegistry{
			index: docid.MakeHashEdgeSet(),
		}
	})
	return globalRegistry
}

func (registry *WsRegistry) Register(client *WsClient) {
	if !docid.IsNil(client.SessionId) {
		registry.index.Add(client.SessionId, client)
	}
	if !docid.IsNil(client.UserId) {
		registry.index.Add(client.UserId, client)
	}
}

func (registry *WsRegistry) Deregister(client *WsClient) {
	if !docid.IsNil(client.SessionId) {
		registry.index.Remove(client.SessionId, client)
	}
	if !docid.IsNil(client.UserId) {
		registry.index.Remove(client.UserId, client)
	}
}

func (registry *WsRegistry) Subscribers(docId docid.DocId) docid.Publisher {
	return registry.index.Targets(docId)
}
