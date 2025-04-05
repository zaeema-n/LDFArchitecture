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

### Start the Neo4j Server with Docker

**Build**

```bash
docker build -t neo4j-custom -f Dockerfile.neo4j .
```

**Run**

```bash
docker run -d \
  --name neo4j-local \
  -p 7474:7474 \
  -p 7687:7687 \
  -v $(pwd)/data:/data \
  -v $(pwd)/logs:/logs \
  -v $(pwd)/plugins:/plugins \
  -v $(pwd)/import:/var/lib/neo4j/import \
  neo4j-custom
```

### Start the Neo4j Server with Docker Composer

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
