import threading

def cpu_load():
    while True:
        pass

thread = threading.Thread(target=cpu_load)
thread.start()
