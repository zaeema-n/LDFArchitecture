# Dockerfile.test
# ===============
#
# This Dockerfile creates a self-contained environment for testing the CRUD API service.
# It builds a single container that includes:
#   1. The Go CRUD service
#   2. MongoDB database
#   3. Neo4j graph database
#   4. All test binaries
#
# The testing process:
# --------------------
# 1. First stage (builder):
#    - Uses golang:1.24 as the base image
#    - Downloads all Go dependencies
#    - Builds the CRUD service binary
#    - Compiles test binaries for each package
#
# 2. Second stage (final):
#    - Uses Ubuntu 22.04 as the base image
#    - Installs system dependencies
#    - Sets up MongoDB and Neo4j databases
#    - Copies the compiled binaries from the builder stage
#    - Configures the environment for testing
#    - Runs tests in an isolated environment
#
# Usage:
# ------
# Build:  docker build -t crud-service-test-standalone -f Dockerfile.test .
# Run:    docker run --rm crud-service-test-standalone
#
# This approach ensures consistent test results regardless of the host environment
# and eliminates the need for external database services during testing.


# -------------------
# Stage 1: Build Go binaries
# -------------------
FROM golang:1.24 AS builder

# Set working directory
WORKDIR /app

# Copy the source code
COPY . .

# Download all dependencies
RUN cd design/crud-api && go mod download

# Build the crud-api test binary
RUN cd design/crud-api && go build ./...
RUN cd design/crud-api && go build -o crud-service cmd/server/service.go cmd/server/utils.go

# Create a directory for test binaries
RUN mkdir -p /app/testbin

# Compile tests for each package
RUN cd design/crud-api/cmd/server && go test -c -o /app/testbin/crud-test .
RUN cd design/crud-api/db/repository/mongo && go test -c -o /app/testbin/mongo-test .
RUN cd design/crud-api/db/repository/neo4j && go test -c -o /app/testbin/neo4j-test .

# -------------------
# Stage 2: Final image
# -------------------
FROM ubuntu:22.04

# Install required packages
RUN apt-get update && apt-get install -y \
    curl \
    gnupg \
    lsb-release \
    wget \
    openjdk-11-jre-headless \
    net-tools \
    nano \
    apt-transport-https \
    software-properties-common \
    && rm -rf /var/lib/apt/lists/*

# Install MongoDB
RUN wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | apt-key add - \
    && echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/6.0 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-6.0.list \
    && apt-get update \
    && apt-get install -y mongodb-org \
    && mkdir -p /data/db

# Install Neo4j
RUN wget -O - https://debian.neo4j.com/neotechnology.gpg.key | apt-key add - \
    && echo 'deb https://debian.neo4j.com stable 4.4' | tee /etc/apt/sources.list.d/neo4j.list \
    && apt-get update \
    && apt-get install -y neo4j=1:4.4.28 cypher-shell \
    && mkdir -p /var/lib/neo4j/data \
    && mkdir -p /var/log/neo4j

# Configure Neo4j logging
RUN sed -i 's/#server.default_listen_address=0.0.0.0/server.default_listen_address=0.0.0.0/' /etc/neo4j/neo4j.conf \
    && sed -i 's/#server.bolt.enabled=true/server.bolt.enabled=true/' /etc/neo4j/neo4j.conf \
    && sed -i 's/#server.bolt.address=0.0.0.0:7687/server.bolt.address=0.0.0.0:7687/' /etc/neo4j/neo4j.conf \
    && sed -i 's/#server.http.enabled=true/server.http.enabled=true/' /etc/neo4j/neo4j.conf \
    && sed -i 's/#server.http.address=0.0.0.0:7474/server.http.address=0.0.0.0:7474/' /etc/neo4j/neo4j.conf \
    && sed -i 's/#dbms.security.auth_enabled=true/dbms.security.auth_enabled=true/' /etc/neo4j/neo4j.conf \
    && echo "dbms.security.procedures.unrestricted=apoc.*" >> /etc/neo4j/neo4j.conf \
    && echo "dbms.logs.debug.level=DEBUG" >> /etc/neo4j/neo4j.conf \
    && echo "dbms.logs.query.enabled=true" >> /etc/neo4j/neo4j.conf \
    && echo "dbms.logs.query.rotation.keep_number=5" >> /etc/neo4j/neo4j.conf \
    && echo "dbms.logs.query.rotation.size=100m" >> /etc/neo4j/neo4j.conf

# Copy test binaries
COPY --from=builder /app/design/crud-api/crud-service /usr/local/bin/
COPY --from=builder /app/testbin/* /usr/local/bin/

# Copy source code
COPY --from=builder /app/design/crud-api /app/design/crud-api
COPY --from=builder /app/design/update-api /app/design/update-api

# Set working directory
WORKDIR /app

# Create log directories and configure logging
RUN mkdir -p /var/log/mongodb && \
    touch /var/log/mongodb/mongod.log && \
    chmod 777 /var/log/mongodb/mongod.log && \
    echo "systemLog:" > /etc/mongod.conf && \
    echo "  destination: file" >> /etc/mongod.conf && \
    echo "  logAppend: true" >> /etc/mongod.conf && \
    echo "  path: /var/log/mongodb/mongod.log" >> /etc/mongod.conf && \
    echo "  logRotate: reopen" >> /etc/mongod.conf && \
    echo "  verbosity: 2" >> /etc/mongod.conf && \
    echo "storage:" >> /etc/mongod.conf && \
    echo "  dbPath: /data/db" >> /etc/mongod.conf && \
    echo "  journal:" >> /etc/mongod.conf && \
    echo "    enabled: true" >> /etc/mongod.conf

# Environment variables
ENV NEO4J_URI=bolt://localhost:7687
ENV NEO4J_USER=neo4j
ENV NEO4J_PASSWORD=neo4j123
ENV MONGO_URI=mongodb://localhost:27017
ENV MONGO_DB_NAME=testdb
ENV MONGO_COLLECTION=metadata

# Expose ports
EXPOSE 7474 7687 27017

# Create startup script with logging
RUN echo '#!/bin/bash\n\
# Start MongoDB with logging\n\
mongod --config /etc/mongod.conf --fork\n\
\n\
# Start Neo4j with logging\n\
neo4j start\n\
\n\
# Wait for Neo4j\n\
until curl -s http://localhost:7474 > /dev/null; do\n\
  echo "Waiting for Neo4j..."\n\
  sleep 2\n\
done\n\
\n\
# Set Neo4j password\n\
echo "Changing initial Neo4j password..."\n\
echo "ALTER CURRENT USER SET PASSWORD FROM '\''neo4j'\'' TO '\''$NEO4J_PASSWORD'\'';" | cypher-shell -u neo4j -p '\''neo4j'\'' -d system\n\
\n\
# Wait for MongoDB\n\
until mongosh --eval "db.version()" > /dev/null 2>&1; do\n\
  echo "Waiting for MongoDB..."\n\
  sleep 2\n\
done\n\
\n\
# Run CRUD service tests first\n\
echo "Running CRUD service tests..."\n\
cd /app/design/crud-api\n\
echo "Running crud-test..."\n\
crud-test -test.v\n\
echo "Running mongo-test..."\n\
mongo-test -test.v\n\
echo "Running neo4j-test..."\n\
neo4j-test -test.v\n\
\n\
# Check if any test failed\n\
if [ $? -ne 0 ]; then\n\
  echo "CRUD tests failed. Exiting."\n\
  exit 1\n\
fi\n\
\n\
# Start CRUD server in background\n\
echo "Starting CRUD server..."\n\
cd /app/design/crud-api\n\
./crud-service &\n\
CRUD_PID=$!\n\
\n\
# Wait for CRUD server to be ready\n\
echo "Waiting for CRUD server..."\n\
# TODO: Uncomment when health endpoint is implemented
# until curl -s http://localhost:50051/health > /dev/null; do\n\
#   echo "Waiting for CRUD server..."\n\
#   sleep 2\n\
# done\n\
# For now, just wait a few seconds for the server to start
echo "Waiting for CRUD server to start..."\n\
sleep 5\n\
\n\
# Run update-api tests\n\
echo "Running update-api tests..."\n\
cd /app/design/update-api\n\
bal test\n\
\n\
# Cleanup\n\
kill $CRUD_PID\n\
\n\
# Keep container running to view logs\n\
tail -f /dev/null' > /start.sh && chmod +x /start.sh

# Start everything
CMD ["/start.sh"]
    