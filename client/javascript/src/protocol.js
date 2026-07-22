/**
 * RESP (Redis Serialization Protocol) Encoder/Decoder
 * Converts between JavaScript values and RESP byte format
 */

class RESPError extends Error {
    constructor(message) {
        super(message);
        this.name = 'RESPError';
    }
}

const ErrInvalidRESP = new RESPError('invalid RESP message');

/**
 * Encode command arguments into RESP protocol format
 * @param {...string} args - Command arguments to encode
 * @returns {Buffer} Encoded RESP message
 */
function encode(...args) {
    let result = `*${args.length}\r\n`;

    for (const arg of args) {
        result += bulkString(arg);
    }

    return Buffer.from(result, 'utf8');
}

/**
 * Decode RESP protocol message into JavaScript value
 * @param {Buffer|Uint8Array|string} message - RESP encoded message
 * @returns {any} Decoded value
 * @throws {RESPError} If message format is invalid
 */
function decode(message) {
    // Convert to Buffer if needed
    if (typeof message === 'string') {
        message = Buffer.from(message, 'utf8');
    } else if (message instanceof Uint8Array && !(message instanceof Buffer)) {
        message = Buffer.from(message);
    }

    if (message.length === 0) {
        throw ErrInvalidRESP;
    }

    const prefix = String.fromCharCode(message[0]);

    switch (prefix) {
        case '+': {
            // Simple String: +OK\r\n
            const value = message.toString('utf8', 1, message.length - 2);
            return value;
        }

        case '$': {
            // Bulk String: $6\r\nfoobar\r\n
            const crlfIndex = message.indexOf('\r\n');
            if (crlfIndex === -1) {
                throw ErrInvalidRESP;
            }

            const lengthStr = message.toString('utf8', 1, crlfIndex);
            const length = parseInt(lengthStr, 10);

            if (isNaN(length)) {
                throw ErrInvalidRESP;
            }

            if (length === -1) {
                // Null bulk string
                return null;
            }

            const start = crlfIndex + 2;
            const end = start + length;

            if (end + 2 !== message.length) {
                throw ErrInvalidRESP;
            }

            const value = message.toString('utf8', start, end);
            return value;
        }

        case '-': {
            // Error: -Error message\r\n
            const errorMsg = message.toString('utf8', 1, message.length - 2);
            throw new RESPError(errorMsg);
        }

        case ':': {
            // Integer: :1000\r\n
            const valueStr = message.toString('utf8', 1, message.length - 2);
            const value = BigInt(valueStr);
            return value;
        }

        case '*': {
            // Array: *2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n
            const crlfIndex = message.indexOf('\r\n');
            if (crlfIndex === -1) {
                throw ErrInvalidRESP;
            }

            const countStr = message.toString('utf8', 1, crlfIndex);
            const count = parseInt(countStr, 10);

            if (isNaN(count)) {
                throw ErrInvalidRESP;
            }

            if (count === -1) {
                // Null array
                return null;
            }

            const result = [];
            let offset = crlfIndex + 2;

            for (let i = 0; i < count; i++) {
                const remaining = message.subarray(offset);

                if (remaining.length === 0) {
                    throw ErrInvalidRESP;
                }

                try {
                    // Find where this element ends and recursively decode
                    const nextCrlf = remaining.indexOf('\r\n');
                    if (nextCrlf === -1) {
                        throw ErrInvalidRESP;
                    }

                    const elementPrefix = String.fromCharCode(remaining[0]);

                    let elementLength = 0;
                    switch (elementPrefix) {
                        case '+':
                        case '-':
                        case ':':
                            elementLength = nextCrlf + 2;
                            break;
                        case '$': {
                            const lengthStr = remaining.toString('utf8', 1, nextCrlf);
                            const length = parseInt(lengthStr, 10);
                            if (length === -1) {
                                elementLength = nextCrlf + 2;
                            } else {
                                elementLength = nextCrlf + 2 + length + 2;
                            }
                            break;
                        }
                        default:
                            throw ErrInvalidRESP;
                    }

                    const element = remaining.subarray(0, elementLength);
                    result.push(decode(element));
                    offset += elementLength;
                } catch (error) {
                    throw ErrInvalidRESP;
                }
            }

            return result;
        }

        default:
            throw ErrInvalidRESP;
    }
}

/**
 * Format a string as a RESP bulk string
 * @param {string} value - String value to format
 * @returns {string} Formatted bulk string
 */
function bulkString(value) {
    return `$${value.length}\r\n${value}\r\n`;
}

export {
    encode,
    decode,
    bulkString,
    ErrInvalidRESP,
    RESPError
};