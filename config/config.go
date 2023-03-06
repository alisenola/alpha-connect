package config

import (
	"fmt"
	"github.com/spf13/viper"
	"gitlab.com/alphaticks/tickstore-go-client/config"
	"os"
	"strings"
)

type Config struct {
	ActorAddress           string
	ActorAdvertisedAddress string
	RegistryAddress        string
	DataStoreAddress       string
	MonitorStoreAddress    string
	DataServerAddress      string
	OpenseaAPIKey          string
	FBinanceWhitelistedIPs []string
	StrictExchange         bool
	StaticLoader           bool
	DialerPoolInterface    string
	DialerPoolIPs          []string
	Accounts               []Account
	Exchanges              []string
	Protocols              []string
	ChainRPCs              []ChainRPC
	DB                     *DataBase
	Store                  *config.StoreClient
}

type Account struct {
	Portfolio        string
	Name             string
	Exchange         string
	ID               string
	ApiKey           string
	ApiSecret        string
	Reconcile        bool
	Listen           bool
	ReadOnly         bool
	MonitorPortfolio bool
	MonitorOrders    bool
	OpeningDate      string
	SOCKS5           string
	FillCollector    bool
	MakerFees        *float64
}

type DataBase struct {
	Migrate          bool
	PostgresHost     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresPort     string
}

type ChainRPC struct {
	Chain    uint32
	Endpoint string
}

func LoadConfig() (*Config, error) {
	// Allow overwrite
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/alpha-connect/")
	configFile := os.Getenv("CONFIG_FILE")
	if configFile != "" {
		f, err := os.Open(configFile)
		if err != nil {
			return nil, fmt.Errorf("error opening provided config file: %v", err)
		}
		if err := viper.ReadConfig(f); err != nil {
			return nil, fmt.Errorf("error reading provided config file: %v", err)
		}
		_ = f.Close()
	} else if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err == nil {
			_ = viper.ReadConfig(f)
			_ = f.Close()
		}
	} else {
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("error reading config file: %v", err)
			}
		} // Find and read the config file
		//if err != nil {             // Handle errors reading the config file
		//	return nil, fmt.Errorf("fatal error config file: %s \n", err)
		//}
	}

	C := &Config{
		ChainRPCs: []ChainRPC{
			{Chain: 1, Endpoint: "wss://eth-mainnet.alchemyapi.io/v2/4kjvftiD6NzHc6kkD1ih3-5wilV--3mz"},
			{Chain: 5, Endpoint: "wss://eth-goerli.g.alchemy.com/v2/YYQHtZhi0iJfppMa41Xj3sUmWuKMxhOf"},
			{Chain: 6, Endpoint: "wss://zksync2-testnet.zksync.dev/ws"},
			{Chain: 10, Endpoint: "wss://opt-mainnet.g.alchemy.com/v2/MAdidiXxtFnW5b4q9pTmBLcTW73SHoMN"},
			{Chain: 147, Endpoint: "wss://polygon-mainnet.g.alchemy.com/v2/PYNN12EJrMrmlxWjy9KZyrYK6GHrErCM"},
			{Chain: 42161, Endpoint: "wss://arb-mainnet.g.alchemy.com/v2/aFDxH7pilsr6I5Ifqp59wVNu_0pd9PaC"},
			{Chain: 1337, Endpoint: "http://127.0.0.1:7545"}, // Ganache
		},
		RegistryAddress: "registry.alphaticks.io:8021",
		StaticLoader:    true,
	}
	if err := viper.Unmarshal(C); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %v", err)
	}

	// Overwrite
	if os.Getenv("DIALER_POOL_INTERFACE") != "" {
		C.DialerPoolInterface = os.Getenv("DIALER_POOL_INTERFACE")
	}
	if os.Getenv("DIALER_POOL_IPS") != "" {
		C.DialerPoolIPs = strings.Split(os.Getenv("DIALER_POOL_IPS"), ",")
	}
	if os.Getenv("EXCHANGES") != "" {
		C.Exchanges = strings.Split(os.Getenv("EXCHANGES"), ",")
	}
	if os.Getenv("PROTOCOLS") != "" {
		C.Protocols = strings.Split(os.Getenv("PROTOCOLS"), ",")
	}
	if os.Getenv("FBINANCE_WHITELISTED_IPS") != "" {
		C.FBinanceWhitelistedIPs = strings.Split(os.Getenv("FBINANCE_WHITELISTED_IPS"), ",")
	}
	fmt.Println(C.FBinanceWhitelistedIPs)

	return C, nil
}
