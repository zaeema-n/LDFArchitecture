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

**Build (Please Prefer Docker compose version for now)**

Build the image (optional since you're using an official base image with no extra steps)

```bash
docker build -t neo4j-service -f Dockerfile.neo4j .
```
Run the container with mounted volumes and env file

```bash
docker run -d \
  --name neo4j-local \
  --platform linux/arm64 \
  -p 7474:7474 \
  -p 7687:7687 \
  --env-file ./config.env \
  -e NEO4J_dbms_memory_pagecache_size=2G \
  -e NEO4J_dbms_memory_heap_initial__size=2G \
  -e NEO4J_dbms_memory_heap_max__size=2G \
  -e NEO4J_dbms_memory_offheap_max__size=1G \
  -v $(pwd)/data:/data \
  -v $(pwd)/logs:/logs \
  -v $(pwd)/plugins:/plugins \
  -v $(pwd)/import:/var/lib/neo4j/import \
  custom-neo4j
```

**Run V2 (with Network)**

```bash
docker run -d \
  --name neo4j-local-v1 \
  --network crud-network \
  -p 7474:7474 \
  -p 7687:7687 \
  -v $(pwd)/data:/data \
  -v $(pwd)/logs:/logs \
  -v $(pwd)/plugins:/plugins \
  -v $(pwd)/import:/var/lib/neo4j/import \
  neo4j-service
```

### Start the Neo4j Server with Docker Composer

Note that we have added `crud-network` as the preferred network for all services
to run as we need each service to be accessible by each one. 

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

