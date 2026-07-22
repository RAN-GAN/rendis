# Rendis JavaScript Client

A JavaScript/Node.js client for Rendis, an in-memory datastore.

## Installation

Install the package via npm:

```bash
npm install rendis
```

## Usage

```javascript
const { Client } = require('rendis'); 

async function run() {
  // Initialize the client connection
  const client = new Client("ws://localhost:8080", "your-auth-key");
  
  try {
    // Ping the server
    await client.ping();

    // Set a key
    await client.set("mykey", "myvalue");
    
    // Get a key
    const val = await client.get("mykey");
    console.log(`Got value: ${val}`);

    // Set an expiration
    await client.expire("mykey", 60);

    // Check if key exists
    const exists = await client.exists("mykey");

    // Delete a key
    await client.del("mykey");
    
  } catch (err) {
    console.error("Error:", err);
  } finally {
    // Close the connection
    client.close();
  }
}

run();
```

## API Methods
The Javascript client utilizes Promises (async/await) for the underlying WebSocket operations:
- `ping()`
- `set(key, value)`
- `get(key)`
- `del(key)`
- `ttl(key)`
- `expire(key, seconds)`
- `exists(key)`
- `close()`
