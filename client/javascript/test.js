import { Client } from './src/index.js';

async function test() {
    console.log('Initializing Rendis Client...');
    const client = new Client('ws://localhost:8080', 'test');

    try {
        await client.connect();
        console.log('✅ Connected successfully!');

        console.log('Sending PING...');
        const pingResult = await client.ping();
        console.log(`✅ PING result: ${pingResult}`);

        console.log('Setting key "test_key" to "hello"...');
        await client.set('test_key', 'hello');
        console.log('✅ Key set successfully!');

        console.log('Getting key "test_key"...');
        const value = await client.get('test_key');
        console.log(`✅ Retrieved value: ${value}`);

    } catch (error) {
        console.error('\n❌ Connection or execution failed.');
        console.error('Note: Make sure a Redis/Rendis server is running on localhost:6379.');
        console.error('Error Details:', error.message);
    } finally {
        client.close();
        console.log('Client connection closed.');
    }
}

test();
