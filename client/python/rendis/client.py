from .connection import ClientConnection
from .protocol import encode, decode, ErrInvalidRESP

class Client:
    def __init__(self, url: str, key: str):
        self._conn = ClientConnection(url, key)

    def close(self):
        self._conn.close()

    def ping(self) -> None:
        self._conn.send(encode("PING"))
        resp = self._conn.receive()
        value = decode(resp)
        if value != "PONG":
            raise Exception(f"unexpected response: {value}")

    def set(self, key: str, value: str) -> None:
        self._conn.send(encode("SET", key, value))
        resp = self._conn.receive()
        decode(resp)

    def get(self, key: str) -> str:
        self._conn.send(encode("GET", key))
        resp = self._conn.receive()
        value = decode(resp)
        if value is None:
            return ""
        if not isinstance(value, str):
            raise ErrInvalidRESP("Expected string response")
        return value

    def delete(self, key: str) -> int:
        self._conn.send(encode("DEL", key))
        resp = self._conn.receive()
        value = decode(resp)
        if not isinstance(value, int):
            raise ErrInvalidRESP("Expected int response")
        return value

    def ttl(self, key: str) -> int:
        self._conn.send(encode("TTL", key))
        resp = self._conn.receive()
        value = decode(resp)
        if not isinstance(value, int):
            raise ErrInvalidRESP("Expected int response")
        return value

    def expire(self, key: str, seconds: int) -> bool:
        self._conn.send(encode("EXPIRE", key, str(seconds)))
        resp = self._conn.receive()
        value = decode(resp)
        if not isinstance(value, int):
            raise ErrInvalidRESP("Expected int response")
        return value == 1

    def exists(self, key: str) -> bool:
        self._conn.send(encode("EXISTS", key))
        resp = self._conn.receive()
        value = decode(resp)
        if not isinstance(value, int):
            raise ErrInvalidRESP("Expected int response")
        return value == 1
