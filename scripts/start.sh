#!/bin/bash
set -e

# Set secrets (avoid ENV warning)
NEO4J_PASSWORD=neo4j123

echo "Starting MongoDB..."
mongod --fork --config /etc/mongodb.conf

echo "Starting Neo4j..."
neo4j start

until curl -s http://localhost:7474 > /dev/null; do
  echo "Waiting for Neo4j..."
  sleep 2
done

echo "Setting Neo4j password..."
echo "ALTER CURRENT USER SET PASSWORD FROM 'neo4j' TO '$NEO4J_PASSWORD';" | cypher-shell -u neo4j -p 'neo4j' -d system

until mongosh --eval "db.version()" > /dev/null 2>&1; do
  echo "Waiting for MongoDB..."
  sleep 2
done

echo "Running CRUD service tests..."
cd /app/design/crud-api
crud-test -test.v && mongo-test -test.v && neo4j-test -test.v

echo "Starting CRUD server..."
./crud-service &
CRUD_PID=$!
sleep 5

echo "Running update-api tests..."
cd /app/design/update-api
bal dist update
bal test

echo "Stopping CRUD server..."
kill $CRUD_PID

tail -f /dev/null 