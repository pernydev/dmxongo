import time
from dmx import DMX
import logging 
import threading
import websocket
import json
import numpy as np

logging.getLogger().addHandler(logging.StreamHandler())

dmx = DMX(num_of_channels=512, serial_number="AL05ZDF6")

def update_universe(ws, message):
    start = time.time()
    dmx.data = np.array(json.loads(message), dtype=np.uint8)
    dmx.send()
    print("Time to send: ", time.time() - start)
    
if __name__ == "__main__":
    ws = websocket.WebSocketApp("ws://127.0.0.1:8080/ws/universe",
                              on_message = lambda ws, message: update_universe(ws, message),
                              on_error = lambda ws, error: print(error),
                              on_close = lambda ws: print("### closed ###"))
    ws.on_open = lambda ws: print("### opened ###")
    ws.run_forever()
    del dmx