# αGraph 

> 💡 **Note (α)**  
> Name needs to be proposed, voted and finalized. 

## Running Services

### Run CRUD API Service

Read about running the [CRUD Service](design/crud-api/README.md)


### Run Update API Service

Read about running the [Update API](design/update-api/README.md)

### Run a sample query

**Create**

```bash
curl -X POST http://localhost:8080/entities \
  -H "Content-Type: application/json" \
  -d '{
    "id": "124",
    "kind": {
      "major": "example",
      "minor": "test"
    },
    "created": "2024-03-15T00:00:00Z",
    "metadata": {
      "type": "entity",
      "status": "active",
      "version": "1.0",
      "description": "test entity"
    }
  }'
```

**Update**

```bash
curl -X PUT http://localhost:8080/entities/123 \
  -H "Content-Type: application/json" \
  -d '{
    "id": "123",
    "kind": {
      "major": "example",
      "minor": "updated"
    },
    "created": "2024-03-15T00:00:00Z",
    "metadata": {
      "type": "entity",
      "status": "updated",
      "version": "2.0",
      "description": "updated entity"
    }
  }'
```

**Delete**

```bash
curl -X DELETE http://localhost:8080/entities/123
```