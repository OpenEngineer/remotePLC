#ifndef pwmRead_h
#define pwmRead_h

#include "arduinoPWMPacket.h"

namespace pwmRead {

void pwmReadSetupUnoPin2();

arduinoPWMPacket pwmRead(arduinoPWMPacket question);

}
#endif
