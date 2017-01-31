#include "arduinoPWMPacket.h"

uint8_t arduinoPWMBuffer[ARDUINO_PWM_PACKET_MAX_PAYLOAD_SIZE];

uint8_t *getArduinoPWMBuffer() {
  return arduinoPWMBuffer;
}
