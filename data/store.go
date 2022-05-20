package data

import (
	"fmt"
	"github.com/melaurent/kafero"
	tickstore_go_client "gitlab.com/alphaticks/tickstore-go-client"
	types "gitlab.com/alphaticks/tickstore-types"
	"gitlab.com/alphaticks/tickstore/storage"
	"gitlab.com/alphaticks/tickstore/store"
	"gitlab.com/alphaticks/tickstore/utils"
	"google.golang.org/grpc"
	"math"
	"time"
)

const (
	DATA_CLIENT_LIVE int64 = -1
	DATA_CLIENT_WEB3 int64 = -2
	DATA_CLIENT_1S   int64 = 1000
	DATA_CLIENT_1M   int64 = DATA_CLIENT_1S * 60
	DATA_CLIENT_1H   int64 = DATA_CLIENT_1M * 60
	DATA_CLIENT_1D   int64 = DATA_CLIENT_1H * 24
)

type DataClient interface {
	GetClient(freq int64) (types.TickstoreClient, int64, error)
	Close() error
}

var ports = map[int64]string{
	DATA_CLIENT_LIVE: "4550",
	DATA_CLIENT_WEB3: "4554",
	DATA_CLIENT_1S:   "4551",
	DATA_CLIENT_1M:   "4552",
	DATA_CLIENT_1H:   "4553",
}

var names = map[int64]string{
	DATA_CLIENT_LIVE: "live",
	DATA_CLIENT_WEB3: "web3",
	DATA_CLIENT_1S:   "1s",
	DATA_CLIENT_1M:   "1m",
	DATA_CLIENT_1H:   "1h",
}

var shardDurations = map[int64]uint64{
	DATA_CLIENT_LIVE: 10000000,
	DATA_CLIENT_WEB3: 1000000000,
	DATA_CLIENT_1S:   10000000,
	DATA_CLIENT_1M:   1000000000,
	DATA_CLIENT_1H:   10000000000,
	//DATA_CLIENT_1D:   100000000000,
}

type StorageClient struct {
	stores       map[int64]*store.Store
	address      string
	opts         []grpc.DialOption
	measurements map[string]string
	cacheDir     string
}

// Lazy loading
func NewStorageClient(cacheDir string, address string, opts ...grpc.DialOption) (*StorageClient, error) {
	s := &StorageClient{
		stores:       make(map[int64]*store.Store),
		address:      address,
		opts:         opts,
		measurements: make(map[string]string),
		cacheDir:     cacheDir,
	}
	s.stores[DATA_CLIENT_LIVE] = nil
	s.stores[DATA_CLIENT_WEB3] = nil
	s.stores[DATA_CLIENT_1S] = nil
	s.stores[DATA_CLIENT_1M] = nil
	s.stores[DATA_CLIENT_1H] = nil

	return s, nil
}

func (s *StorageClient) GetStore(freq int64) (*store.Store, int64, error) {
	var minScore int64 = math.MaxInt64
	var cfreq int64
	if freq == DATA_CLIENT_LIVE || freq == DATA_CLIENT_WEB3 {
		cfreq = freq
	} else {
		for f := range s.stores {
			if f <= freq {
				score := freq - f
				if score < minScore {
					minScore = score
					cfreq = f
				}
			}
		}
	}

	if s.stores[cfreq] == nil {
		// Construct store
		// Construct cache storage
		cacheStrg, err := storage.NewFsStorage(
			kafero.NewBasePathFs(kafero.NewOsFs(), s.cacheDir),
			names[cfreq],
			3,
			shardDurations[cfreq],
			10*time.Minute,
			500,
			false,
			true,
			true,
			true,
			true)
		if err != nil {
			return nil, 0, fmt.Errorf("error starting cache storage: %v", err)
		}
		strg, err := storage.NewClientStorage(cacheStrg, s.address+":"+ports[cfreq], s.opts...)
		if err != nil {
			return nil, 0, err
		}

		// Build store
		str, err := store.NewStore(strg)
		if err != nil {
			return nil, 0, fmt.Errorf("error building store: %v", err)
		}
		s.stores[cfreq] = str
	}
	return s.stores[cfreq], cfreq, nil
}

func (s *StorageClient) GetClient(freq int64) (types.TickstoreClient, int64, error) {
	str, cfreq, err := s.GetStore(freq)
	if err != nil {
		return nil, cfreq, err
	}
	return store.NewLocalClient(str), cfreq, nil
}

func (s *StorageClient) Close() error {
	errs := utils.NewErrorSet()
	for freq, str := range s.stores {
		if str != nil {
			if err := str.Close(); err != nil {
				errs.Push(err)
			}
			s.stores[freq] = nil
		}
	}
	if errs.Len() > 0 {
		return errs
	} else {
		return nil
	}
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
	s.stores[DATA_CLIENT_WEB3] = nil
	s.stores[DATA_CLIENT_LIVE] = nil
	s.stores[DATA_CLIENT_1S] = nil
	s.stores[DATA_CLIENT_1M] = nil
	s.stores[DATA_CLIENT_1H] = nil
	//s.stores[DATA_CLIENT_1D] = nil

	return s, nil
}

func (s *StoreClient) GetClient(freq int64) (types.TickstoreClient, int64, error) {
	var minScore int64 = math.MaxInt64
	var cfreq int64
	if freq == DATA_CLIENT_LIVE || freq == DATA_CLIENT_WEB3 {
		cfreq = freq
	} else {
		for f := range s.stores {
			if f <= freq {
				score := freq - f
				if score < minScore {
					minScore = score
					cfreq = f
				}
			}
		}
	}
	if s.stores[cfreq] == nil {
		// Construct store
		str, err := tickstore_go_client.NewRemoteClient(s.address+":"+ports[cfreq], s.opts...)
		if err != nil {
			return nil, 0, err
		}

		s.stores[cfreq] = str
	}
	return s.stores[cfreq], cfreq, nil
}

func (s *StoreClient) Close() error {
	return nil
}
