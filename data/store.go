package data

import (
	"fmt"
	tickstore_go_client "gitlab.com/tachikoma.ai/tickstore-go-client"
	"gitlab.com/tachikoma.ai/tickstore/storage"
	"gitlab.com/tachikoma.ai/tickstore/store"
	"google.golang.org/grpc"
	"math"
)

const (
	DATA_CLIENT_LIVE int64 = 0
	DATA_CLIENT_1S   int64 = 1000
	DATA_CLIENT_1M   int64 = DATA_CLIENT_1S * 60
	DATA_CLIENT_1H   int64 = DATA_CLIENT_1M * 60
	DATA_CLIENT_1D   int64 = DATA_CLIENT_1H * 24
)

type DataClient interface {
	GetClient(freq int64) (tickstore_go_client.TickstoreClient, error)
}

var ports = map[int64]string{
	DATA_CLIENT_LIVE: "4550",
	DATA_CLIENT_1S:   "4551",
	DATA_CLIENT_1M:   "4552",
	DATA_CLIENT_1H:   "4553",
	DATA_CLIENT_1D:   "4554",
}

type StorageClient struct {
	stores       map[int64]*store.Store
	address      string
	opts         []grpc.DialOption
	measurements map[string]string
}

// Lazy loading
func NewStorageClient(address string, opts ...grpc.DialOption) (*StorageClient, error) {
	s := &StorageClient{
		stores:       make(map[int64]*store.Store),
		address:      address,
		opts:         opts,
		measurements: make(map[string]string),
	}
	s.stores[DATA_CLIENT_LIVE] = nil
	s.stores[DATA_CLIENT_1S] = nil
	s.stores[DATA_CLIENT_1M] = nil
	s.stores[DATA_CLIENT_1H] = nil
	s.stores[DATA_CLIENT_1D] = nil

	return s, nil
}

func (s *StorageClient) GetStore(freq int64) (*store.Store, error) {
	var minScore int64 = math.MaxInt64
	var cfreq int64
	for f := range s.stores {
		if f <= freq {
			score := freq - f
			if score < minScore {
				minScore = score
				cfreq = f
			}
		}
	}
	if s.stores[cfreq] == nil {
		// Construct store
		strg, err := storage.NewClientStorage(s.address+":"+ports[cfreq], s.opts...)
		if err != nil {
			return nil, err
		}

		// Build store
		str, err := store.NewStore(strg)
		if err != nil {
			return nil, fmt.Errorf("error building store: %v", err)
		}
		s.stores[cfreq] = str
	}
	return s.stores[cfreq], nil
}

func (s *StorageClient) GetClient(freq int64) (tickstore_go_client.TickstoreClient, error) {
	var minScore int64 = math.MaxInt64
	var cfreq int64
	for f := range s.stores {
		if f <= freq {
			score := freq - f
			if score < minScore {
				minScore = score
				cfreq = f
			}
		}
	}
	if s.stores[cfreq] == nil {
		// Construct store
		strg, err := storage.NewClientStorage(s.address+":"+ports[cfreq], s.opts...)
		if err != nil {
			return nil, err
		}

		// Build store
		str, err := store.NewStore(strg)
		if err != nil {
			return nil, fmt.Errorf("error building store: %v", err)
		}
		s.stores[cfreq] = str
	}
	return store.NewLocalClient(s.stores[cfreq]), nil
}

type StoreClient struct {
	stores       map[int64]*tickstore_go_client.RemoteClient
	address      string
	opts         []grpc.DialOption
	measurements map[string]string
}

// Lazy loading
func NewStoreClient(address string, opts ...grpc.DialOption) (*StoreClient, error) {
	s := &StoreClient{
		stores:       make(map[int64]*tickstore_go_client.RemoteClient),
		address:      address,
		opts:         opts,
		measurements: make(map[string]string),
	}
	s.stores[DATA_CLIENT_LIVE] = nil
	s.stores[DATA_CLIENT_1S] = nil
	s.stores[DATA_CLIENT_1M] = nil
	s.stores[DATA_CLIENT_1H] = nil
	s.stores[DATA_CLIENT_1D] = nil

	return s, nil
}

func (s *StoreClient) GetClient(freq int64) (tickstore_go_client.TickstoreClient, error) {
	var minScore int64 = math.MaxInt64
	var cfreq int64
	for f := range s.stores {
		if f <= freq {
			score := freq - f
			if score < minScore {
				minScore = score
				cfreq = f
			}
		}
	}
	if s.stores[cfreq] == nil {
		// Construct store
		str, err := tickstore_go_client.NewRemoteClient(s.address+":"+ports[cfreq], s.opts...)
		if err != nil {
			return nil, err
		}

		s.stores[cfreq] = str
	}
	return s.stores[cfreq], nil
}
