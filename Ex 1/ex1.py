#!/usr/bin/python
from threading import Thread

i = 0

def adder():
    global i
    for x in range(0, 1000000):     
        i += 1

def subtract():
    global i
    for x in range(0, 1000000):     
        i -= 1


def main():
    adder_thr = Thread(target = adder)
    subtract_thr = Thread(target = subtract)

    adder_thr.start()
    subtract_thr.start()
    adder_thr.join()
    subtract_thr.join()
    print("Done: " + str(i))


main()
