class ErrInvalidRESP(Exception):
    pass

def encode(*args: str) -> bytes:
    encoded = f"*{len(args)}\r\n"
    for arg in args:
        encoded += bulk_string(arg)
    return encoded.encode("utf-8")

def bulk_string(value: str) -> str:
    return f"${len(value)}\r\n{value}\r\n"

def decode(message: bytes):
    if not message:
        raise ErrInvalidRESP("Empty message")
    
    prefix = chr(message[0])
    
    if prefix == '+':
        return message[1:-2].decode('utf-8')
    elif prefix == '$':
        try:
            cr_index = message.index(b'\r')
        except ValueError:
            raise ErrInvalidRESP("Invalid bulk string header")
            
        length = int(message[1:cr_index].decode('utf-8'))
        if length == -1:
            return None
            
        header_end = message.find(b'\r\n')
        if header_end == -1:
            raise ErrInvalidRESP("Invalid bulk string header end")
            
        start = header_end + 2
        end = start + length
        
        if end + 2 != len(message):
            raise ErrInvalidRESP("Invalid bulk string length")
            
        return message[start:end].decode('utf-8')
        
    elif prefix == '-':
        error_msg = message[1:-2].decode('utf-8')
        raise Exception(error_msg)
        
    elif prefix == ':':
        return int(message[1:-2].decode('utf-8'))
        
    else:
        raise ErrInvalidRESP(f"Unknown prefix: {prefix}")
