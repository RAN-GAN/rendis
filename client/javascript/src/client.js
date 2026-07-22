import { RendisClient } from './connection.js';
import { encode, decode } from './protocol.js';
export class Client {
    constructor(url, key) {
        this.url = url;
        this.key = key;
        this.conn = new RendisClient(url, key);
    }

    // Send command to Redis server
    async send(data) {
        this.conn.send(data);
    }

    // Receive response from Redis server
    async receive() {
        return this.conn.receive();
    }



    // PING - Test connection
    async ping() {
        try {
            const data = encode('PING');
            await this.send(data);

            const resp = await this.receive();
            const value = decode(resp);

            if (typeof value !== 'string') {
                throw new Error('Invalid RESP response');
            }

            if (value !== 'PONG') {
                throw new Error(`Unexpected response: ${value}`);
            }

            return true;
        } catch (error) {
            throw error;
        }
    }

    // SET - Set a key-value pair
    async set(key, value) {
        try {
            const data = encode('SET', key, value);
            await this.send(data);

            const resp = await this.receive();
            decode(resp);
            return true;
        } catch (error) {
            throw error;
        }
    }

    // GET - Get value by key
    async get(key) {
        try {
            const data = encode('GET', key);
            await this.send(data);

            const resp = await this.receive();
            const value = decode(resp);

            if (typeof value !== 'string' && value !== null) {
                throw new Error('Invalid RESP response');
            }

            return value;
        } catch (error) {
            throw error;
        }
    }

    // DEL - Delete a key
    async del(key) {
        try {
            const data = encode('DEL', key);
            await this.send(data);

            const resp = await this.receive();
            const value = decode(resp);

            if (typeof value !== 'bigint') {
                throw new Error('Invalid RESP response');
            }

            return Number(value);
        } catch (error) {
            throw error;
        }
    }

    // TTL - Get time to live for a key
    async ttl(key) {
        try {
            const data = encode('TTL', key);
            await this.send(data);

            const resp = await this.receive();
            const value = decode(resp);

            if (typeof value !== 'bigint') {
                throw new Error('Invalid RESP response');
            }

            return Number(value);
        } catch (error) {
            throw error;
        }
    }

    // EXPIRE - Set expiration time
    async expire(key, seconds) {
        try {
            const data = encode('EXPIRE', key, String(seconds));
            await this.send(data);

            const resp = await this.receive();
            const value = decode(resp);

            if (typeof value !== 'bigint') {
                throw new Error('Invalid RESP response');
            }

            return Number(value) === 1;
        } catch (error) {
            throw error;
        }
    }

    // EXISTS - Check if key exists
    async exists(key) {
        try {
            const data = encode('EXISTS', key);
            await this.send(data);

            const resp = await this.receive();
            const value = decode(resp);

            if (typeof value !== 'bigint') {
                throw new Error('Invalid RESP response');
            }

            return Number(value) === 1;
        } catch (error) {
            throw error;
        }
    }

    // Connect to Redis server
    async connect() {
        return this.conn.connect();
    }

    // Close connection
    close() {
        this.conn.close();
    }
}