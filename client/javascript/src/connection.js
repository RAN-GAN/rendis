// Browser version using native WebSocket API
class RendisClientBrowser {
  constructor(url, key) {
    this.url = url;
    this.key = key;
    this.conn = null;
    this.messageQueue = [];
    this.readyPromise = null;
  }

  async connect() {
    return new Promise((resolve, reject) => {
      try {
        this.conn = new WebSocket(this.url + '/connect');
        
        // Set custom header before connecting (not directly supported in browser)
        // Note: Custom headers in WebSocket handshake require server-side handling
        // Some servers allow passing auth via URL query string instead
        this.conn.binaryType = 'arraybuffer';

        this.conn.onopen = () => {
          // Send auth key as first message if needed
          if (this.key) {
            this.conn.send(this.key);
          }
          resolve();
        };

        this.conn.onerror = (error) => {
          reject(new Error('WebSocket error: ' + error));
        };

        this.conn.onclose = () => {
          this.conn = null;
        };
      } catch (error) {
        reject(error);
      }
    });
  }

  close() {
    if (this.conn) {
      this.conn.close();
      this.conn = null;
    }
  }

  send(data) {
    if (!this.conn || this.conn.readyState !== WebSocket.OPEN) {
      throw new Error('WebSocket is not connected');
    }

    // Convert to ArrayBuffer if needed
    if (typeof data === 'string') {
      data = new TextEncoder().encode(data);
    }

    this.conn.send(data);
  }

  receive() {
    return new Promise((resolve, reject) => {
      if (!this.conn) {
        reject(new Error('WebSocket is not connected'));
        return;
      }

      const messageHandler = (event) => {
        this.conn.removeEventListener('message', messageHandler);
        this.conn.removeEventListener('error', errorHandler);
        resolve(new Uint8Array(event.data));
      };

      const errorHandler = (error) => {
        this.conn.removeEventListener('message', messageHandler);
        this.conn.removeEventListener('error', errorHandler);
        reject(new Error('WebSocket error: ' + error));
      };

      this.conn.addEventListener('message', messageHandler);
      this.conn.addEventListener('error', errorHandler);
    });
  }
}

// Node.js version using ws package
class RendisClientNode {
  constructor(url, key) {
    this.url = url;
    this.key = key;
    this.conn = null;
  }

  async connect() {
    return new Promise(async (resolve, reject) => {
      try {
        const wsModule = await import('ws');
        const WebSocket = wsModule.default || wsModule;
        const http = await import('http');

        const headers = {
          'X-RENDIS-Key': this.key
        };

        this.conn = new WebSocket(this.url + '/connect', {
          headers: headers
        });

        this.conn.binaryType = 'arraybuffer';

        this.conn.on('open', () => {
          resolve();
        });

        this.conn.on('error', (error) => {
          reject(new Error('WebSocket error: ' + error.message));
        });

        this.conn.on('close', () => {
          this.conn = null;
        });
      } catch (error) {
        reject(error);
      }
    });
  }

  close() {
    if (this.conn) {
      this.conn.close();
      this.conn = null;
    }
  }

  send(data) {
    if (!this.conn || this.conn.readyState !== 1) { // 1 = OPEN
      throw new Error('WebSocket is not connected');
    }

    // Convert to Buffer if needed
    if (typeof data === 'string') {
      data = Buffer.from(data);
    }

    this.conn.send(data);
  }

  receive() {
    return new Promise((resolve, reject) => {
      if (!this.conn) {
        reject(new Error('WebSocket is not connected'));
        return;
      }

      const messageHandler = (data) => {
        this.conn.removeListener('message', messageHandler);
        this.conn.removeListener('error', errorHandler);
        resolve(Buffer.from(data));
      };

      const errorHandler = (error) => {
        this.conn.removeListener('message', messageHandler);
        this.conn.removeListener('error', errorHandler);
        reject(new Error('WebSocket error: ' + error.message));
      };

      this.conn.once('message', messageHandler);
      this.conn.once('error', errorHandler);
    });
  }
}

// Export based on environment
export const RendisClient = typeof window !== 'undefined' ? RendisClientBrowser : RendisClientNode;

// Also export both for flexibility
export { RendisClientBrowser as Browser, RendisClientNode as Node };
export default RendisClient;