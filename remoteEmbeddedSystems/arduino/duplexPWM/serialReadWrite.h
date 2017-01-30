#ifndef serialReadWrite_h
#define serialReadWrite_h

#include "arduinoPWMPacket.h"

void serialSetup(int bitRate);

// blocks until full message is read
arduinoPWMPacket serialReadMessage();

void serialWriteMessage(arduinoPWMPacket p);

#endif
