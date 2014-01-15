#!/usr/bin/python

from threading import Thread, RLock

i = 0
lock = RLock()
def adder():
    global i
    for x in range(0, 1000000):     
        lock.acquire()
        i += 1
        lock.release()

def subtract():
    global i
    for x in range(0, 1000000):     
        lock.acquire()
        i -= 1
        lock.release()

def main():
    adder_thr = Thread(target = adder)
    subtract_thr = Thread(target = subtract)
    adder_thr.start()
    subtract_thr.start()
    adder_thr.join()
    subtract_thr.join()
    print("Done: " + str(i))


main()
