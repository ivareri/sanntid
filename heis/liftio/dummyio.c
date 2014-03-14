// Wrapper for libComedi I/O.
// These functions provide and interface to libComedi limited to use in
// the real time lab.
//
// 2006, Martin Korsgaard

//
// YOU DO NOT NEED TO EDIT THIS FILE
//

#include "dummyio.h"
#include "channels.h"


int io_init(){
    return 1;
}

void io_set_bit(int channel){
}

void io_clear_bit(int channel){
}

void io_write_analog(int channel, int value){
}

int io_read_bit(int channel){
    return 0;
}

int io_read_analog(int channel){
    return 200;
}
