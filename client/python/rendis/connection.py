import websocket

class ClientConnection:
    def __init__(self, url: str, key: str):
        connect_url = url + "/connect"
        # Using websocket-client to connect synchronously
        self.conn = websocket.create_connection(
            connect_url,
            header=[f"X-RENDIS-Key: {key}"]
        )

    def close(self):
        self.conn.close()

    def send(self, data: bytes):
        self.conn.send_binary(data)

    def receive(self) -> bytes:
        return self.conn.recv()
