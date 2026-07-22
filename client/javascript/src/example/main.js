/**
 * RESP Protocol Usage Examples
 */

import { encode, decode, ErrInvalidRESP } from '../protocol.js';

// ============================================
// ENCODING EXAMPLES
// ============================================

console.log('=== ENCODING EXAMPLES ===\n');

// Example 1: Simple PING command
const pingEncoded = encode('PING');
console.log('PING command:');
console.log('Encoded:', pingEncoded);
console.log('Hex:', pingEncoded.toString('hex'));
console.log('');

// Example 2: SET command with key and value
const setEncoded = encode('SET', 'mykey', 'myvalue');
console.log('SET command:');
console.log('Encoded:', setEncoded);
console.log('Hex:', setEncoded.toString('hex'));
console.log('');

// Example 3: GET command
const getEncoded = encode('GET', 'mykey');
console.log('GET command:');
console.log('Encoded:', getEncoded);
console.log('Hex:', getEncoded.toString('hex'));
console.log('');

// Example 4: DEL command
const delEncoded = encode('DEL', 'mykey');
console.log('DEL command:');
console.log('Encoded:', delEncoded);
console.log('Hex:', delEncoded.toString('hex'));
console.log('');

// ============================================
// DECODING EXAMPLES
// ============================================

console.log('=== DECODING EXAMPLES ===\n');

// Example 1: Decode simple string (PONG response)
const pongResponse = Buffer.from('+PONG\r\n', 'utf8');
try {
    const decodedPong = decode(pongResponse);
    console.log('Decoded PONG:', decodedPong);
} catch (error) {
    console.error('Error:', error.message);
}

// Example 2: Decode bulk string (GET response with value)
const bulkStringResponse = Buffer.from('$7\r\nmyvalue\r\n', 'utf8');
try {
    const decodedBulk = decode(bulkStringResponse);
    console.log('Decoded bulk string:', decodedBulk);
} catch (error) {
    console.error('Error:', error.message);
}

// Example 3: Decode null bulk string
const nullBulkResponse = Buffer.from('$-1\r\n', 'utf8');
try {
    const decodedNull = decode(nullBulkResponse);
    console.log('Decoded null bulk string:', decodedNull);
} catch (error) {
    console.error('Error:', error.message);
}

// Example 4: Decode integer (DEL response)
const integerResponse = Buffer.from(':1\r\n', 'utf8');
try {
    const decodedInt = decode(integerResponse);
    console.log('Decoded integer:', decodedInt, 'Type:', typeof decodedInt);
} catch (error) {
    console.error('Error:', error.message);
}

// Example 5: Decode error
const errorResponse = Buffer.from('-ERR unknown command\r\n', 'utf8');
try {
    const decodedError = decode(errorResponse);
    console.log('Decoded error:', decodedError);
} catch (error) {
    console.error('Caught error:', error.message);
}

// Example 6: Decode array
const arrayResponse = Buffer.from('*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n', 'utf8');
try {
    const decodedArray = decode(arrayResponse);
    console.log('Decoded array:', decodedArray);
} catch (error) {
    console.error('Error:', error.message);
}

// ============================================
// ENCODING-DECODING ROUNDTRIP
// ============================================

console.log('\n=== ROUNDTRIP EXAMPLES ===\n');

const commands = [
    ['PING'],
    ['SET', 'name', 'John'],
    ['GET', 'name'],
    ['DEL', 'name'],
    ['EXPIRE', 'key', '3600']
];

for (const cmd of commands) {
    const encoded = encode(...cmd);
    console.log(`Command: ${cmd.join(' ')}`);
    console.log(`Encoded: ${encoded.toString('hex')}`);
    console.log(`Readable: ${encoded.toString('utf8').replace(/\r\n/g, '\\r\\n')}`);
    console.log('');
}

// ============================================
// PRACTICAL CLIENT USAGE
// ============================================

console.log('=== PRACTICAL USAGE IN CLIENT ===\n');

class SimpleRedisClient {
    constructor(socket) {
        this.socket = socket;
    }

    async sendCommand(...args) {
        const encoded = encode(...args);
        this.socket.write(encoded);
    }

    async parseResponse(data) {
        try {
            const value = decode(data);
            return { success: true, value };
        } catch (error) {
            return { success: false, error: error.message };
        }
    }
}

// Usage example
console.log('Example client usage:');
console.log('const client = new SimpleRedisClient(socket);');
console.log('await client.sendCommand("SET", "key", "value");');
console.log('const result = client.parseResponse(responseBuffer);');