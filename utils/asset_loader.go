package utils

import (
	"cloud.google.com/go/storage"
	goContext "context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/models"
	"io/ioutil"
	"reflect"
	"strings"
	"time"
)

type checkAsset struct{}
type Ready struct{}

type AssetLoader struct {
	assetFile string
	assetMD5  string
	logger    *log.Logger
	ticker    *time.Ticker
}

func NewAssetLoaderProducer(assetFile string) actor.Producer {
	return func() actor.Actor {
		return NewAssetLoader(assetFile)
	}
}

func NewAssetLoader(assetFile string) actor.Actor {
	return &AssetLoader{
		assetFile: assetFile,
	}
}

func (state *AssetLoader) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.logger.Error("error initializing actor", log.Error(err))
			panic(err)
		}
		state.logger.Info("actor started")

	case *actor.Stopping:
		if err := state.Clean(context); err != nil {
			state.logger.Error("error stopping", log.Error(err))
			panic(err)
		}
		state.logger.Info("actor stopping")

	case *actor.Stopped:
		state.logger.Info("actor stopped")

	case *actor.Restarting:
		if err := state.Clean(context); err != nil {
			state.logger.Error("error restarting", log.Error(err))
			// Attention, no panic in restarting or infinite loop
		}
		state.logger.Info("actor restarting")

	case *checkAsset:
		if err := state.checkAsset(context); err != nil {
			state.logger.Error("error checkAsset", log.Error(err))
			panic(err)
		}

	case *Ready:
		if err := state.onReady(context); err != nil {
			state.logger.Error("error onReady", log.Error(err))
			panic(err)
		}
	}
}

func (state *AssetLoader) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	ticker := time.NewTicker(20 * time.Second)
	state.ticker = ticker
	go func(pid *actor.PID) {
		for {
			select {
			case _ = <-ticker.C:
				context.Send(pid, &checkAsset{})
			case <-time.After(2 * time.Minute):
				// timer stopped, we leave
				return
			}
		}
	}(context.Self())

	return state.checkAsset(context)
}

func (state *AssetLoader) Clean(context actor.Context) error {
	if state.ticker != nil {
		state.ticker.Stop()
		state.ticker = nil
	}

	return nil
}

func (state *AssetLoader) onReady(context actor.Context) error {
	context.Respond(&Ready{})
	return nil
}

func (state *AssetLoader) checkAsset(context actor.Context) error {
	fmt.Println("CHECK")
	var data []byte
	if state.assetFile[0:5] == "gs://" {
		ctx := goContext.Background()
		client, err := storage.NewClient(ctx)
		if err != nil {
			return err
		}

		splits := strings.Split(state.assetFile[5:], "/")
		bucketName := splits[0]
		if bucketName == "" {
			return fmt.Errorf("configurator bucket name not set")
		}
		bucket := client.Bucket(bucketName)

		configHandle := bucket.Object(splits[1])
		attrs, err := configHandle.Attrs(ctx)
		if err != nil {
			return fmt.Errorf("error getting attribute: %v", err)
		}
		if string(attrs.MD5) == state.assetMD5 {
			return nil
		}

		rc, err := configHandle.NewReader(ctx)
		if err != nil {
			return err
		}
		bytes, err := ioutil.ReadAll(rc)
		rc.Close()
		if err != nil {
			return err
		}
		data = bytes
	} else {
		file, err := ioutil.ReadFile(state.assetFile)
		if err != nil {
			return fmt.Errorf("failed to read configuration: %v", err)
		}
		data = file
	}
	assets := make(map[uint32]models.Asset)
	err := json.Unmarshal(data, &assets)
	if err != nil {
		return fmt.Errorf("error parsing config: %v", err)
	}

	if err := constants.LoadAssets(assets); err != nil {
		return err
	}

	hash := md5.New()
	hash.Write(data)
	state.assetMD5 = string(hash.Sum(nil))

	return nil
}
