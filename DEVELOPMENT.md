## Development

For the development mode make sure you `source` the file containing the secrets. For instance 
you can keep a secret file like `ldf.testing.profile`

```bash
export MONGO_TESTING_DB_URI=""
export MONGO_TESTING_DB=""
export MONGO_TESTING_COLLECTION=""

export NEO4J_TESTING_DB_URI=""
export NEO4J_TESTING_USERNAME=""
export NEO4J_TESTING_PASSWORD=""
```

`config.env` or secrets in Github would make up `NEO4J_AUTH=${NEO4J_TESTING_USERNAME}/${NEO4J_TESTING_PASSWORD}`.

In the same terminal or ssh session, do the following;

This will start an instance of the neo4j database server. 

### Start the Neo4j Server

```bash
docker compose up --build
```

Go to `http://localhost:7474/browser/` and you can access the neo4j browser. 

### Shutdown the Neo4j Server

```bash
docker compose down -v
```

### BackUp Server Data (TODO)


### Restore Server Data (TODO)


### Docker (Standalone): WIP

```bash
docker build -t all-services-test-standalone -f Dockerfile .
```

```bash
docker run --rm all-services-test-standalone
```

### Docker Compose

Use the `docker compose` to up the services to run tests and to check the current version of the software is working. 

#### Up the Services

`docker compose up` 

#### Down the Services

`docker compose down` 

#### Get services up independently 

MongoDB Service

`docker compose up -d mongodb`

Neo4j Service 

`docker compose up -d neo4j` 

Build CRUD Service

`docker compose build crud` 

And to up it `docker compose up crud`
