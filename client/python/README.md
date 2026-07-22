# Rendis Python Client

A Python client for Rendis, an in-memory datastore.

## Installation

You can install the client using pip:

```bash
pip install rendis
```

## Usage

```python
from rendis import Client

# Initialize the client
client = Client("ws://localhost:8080", "your-auth-key")

try:
    # Ping the server
    client.ping()

    # Set a key
    client.set("mykey", "myvalue")

    # Get a key
    val = client.get("mykey")
    print(f"Got value: {val}")
    
    # Set an expiration (e.g., 60 seconds)
    client.expire("mykey", 60)
    
    # Check Time To Live (TTL)
    ttl = client.ttl("mykey")
    print(f"TTL: {ttl}")

    # Check if a key exists
    exists = client.exists("mykey")
    
    # Delete a key
    client.delete("mykey")

finally:
    # Always ensure the connection is closed
    client.close()
```

## API Methods
- `Client(url: str, key: str)`
- `ping() -> None`
- `set(key: str, value: str) -> None`
- `get(key: str) -> str`
- `delete(key: str) -> int`
- `ttl(key: str) -> int`
- `expire(key: str, seconds: int) -> bool`
- `exists(key: str) -> bool`
- `close() -> None`
