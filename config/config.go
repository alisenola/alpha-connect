package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

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
		if err != nil {
			return nil, fmt.Errorf("error opening provided config file: %v", err)
		}
		if err := viper.ReadConfig(f); err != nil {
			return nil, fmt.Errorf("error reading provided config file: %v", err)
		}
		_ = f.Close()
	} else {
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			return nil, fmt.Errorf("fatal error config file: %s \n", err)
		}
	}

	C := &Config{}
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

	return C, nil
}

type Config struct {
	ActorAddress           string
	ActorAdvertisedAddress string
	RegistryAddress        string
	DataStoreAddress       string
	DataServerAddress      string
	StrictExchange         bool
	DialerPoolInterface    string
	DialerPoolIPs          []string
	Accounts               []Account
	Exchanges              []string
	Protocols              []string
	Chains                 []string
	DB                     *DataBase
}

type Account struct {
	Name      string
	Exchange  string
	ID        string
	ApiKey    string
	ApiSecret string
	Reconcile bool
	Listen    bool
	ReadOnly  bool
}

type DataBase struct {
	Migrate          bool
	PostgresHost     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresPort     string
}