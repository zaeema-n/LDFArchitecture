package config

type MongoConfig struct {
	URI        string `env:"MONGO_URI"`
	DBName     string `env:"MONGO_DB_NAME"`
	Collection string `env:"MONGO_COLLECTION"`
}
