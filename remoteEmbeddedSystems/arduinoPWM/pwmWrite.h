#ifndef pwmWrite_h
#define pwmWrite_h

#include "arduinoPWMPacket.h"

void pwmWriteSetup(int outputPin);

// return p upon success
arduinoPWMPacket pwmWrite(arduinoPWMPacket p);

#endif
