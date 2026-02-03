# Signalling Server API

This service provides peer discovery and basic client registration for SyncMesh.

## Endpoints

### GET /ping
Health check with server timestamp.

Response:
```json
{
	"status": "available",
	"systemInfo": {
		"serverTimestamp": "2026-02-03T20:03:11Z"
	}
}
```

### POST /register
Register the caller and return a `clientId`.

Request body:
```json
{
	"localIp": "192.168.1.50",
	"localPort": 4242
}
```

Response:
```json
{
	"status": "success",
	"clientId": "a7c4fce7b9b74c8b5f1b0a7db5e2f5bb"
}
```

Notes:
- `publicIp` and `publicPort` are captured from the connection's `RemoteAddr`.
- If the body is empty, local fields are omitted.

### POST /heartbeat?clientId=...
Refresh the client's `LastSeen` to avoid TTL pruning.

Response:
```json
{
	"status": "success"
}
```

Errors:
- `400` if `clientId` is missing.
- `404` if the client is not registered (or expired).

### POST /unregister?clientId=...
Remove a client from the registry.

Response:
```json
{
	"status": "success"
}
```

### GET /discover
List known clients and the connection info needed to contact them.

Response:
```json
{
	"status": "success",
	"clients": [
		{
			"clientId": "a7c4fce7b9b74c8b5f1b0a7db5e2f5bb",
			"publicIp": "203.0.113.10",
			"publicPort": 51234,
			"localIp": "192.168.1.50",
			"localPort": 4242
		}
	]
}
```

## TTL Behavior
Clients are removed if they have not sent a heartbeat within 5 minutes. The registry is pruned on register, discover, heartbeat, and unregister.
