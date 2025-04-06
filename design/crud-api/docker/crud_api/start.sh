#!/bin/bash
set -e

# Start Neo4j
neo4j start

# Wait for MongoDB to be ready (already running in the base image)
until mongosh --eval "db.version()" >/dev/null 2>&1; do
  echo "Waiting for MongoDB to be ready..."
  sleep 2
done

# Start the CRUD service
exec crud-service 