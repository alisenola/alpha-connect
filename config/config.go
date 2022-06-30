package config

type Config struct {
	ActorAddress           string
	ActorAdvertisedAddress string
	RegistryAddress        string
	DataStoreAddress       string
	DataServerAddress      string
	StrictExchange         bool
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
