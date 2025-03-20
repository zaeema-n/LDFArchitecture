package config

type MongoConfig struct {
	URI        string `env:"MONGO_URI"`
	DBName     string `env:"MONGO_DB_NAME"`
	Collection string `env:"MONGO_COLLECTION"`
}

type Neo4jConfig struct {
	URI      string `env:"NEO4J_URI"`
	Username string `env:"NEO4J_USER"`
	Password string `env:"NEO4J_PASSWORD"`
}
