#include <Arduino.h>

#include "arduinoPWMPacket.h"

//
// internal functions and variables
//
int pwmWriteOutputPin;

void pwmWriteDelayMicroseconds(uint16_t d) {
  uint16_t dMicros = d%1000;
  uint16_t dMillis = (d - dMicros)/1000;

  delayMicroseconds(int(dMicros)); // doesn't work for dMicros > 1000, so we need to use delay for the number of milliseconds

  if (dMillis > 0) {
    delay(int(dMillis));
  }
}

void pwmWriteBit(uint8_t bit, uint16_t pulseWidth) {
  if (bit != 0) {
    digitalWrite(pwmWriteOutputPin, HIGH);
  } else {
    digitalWrite(pwmWriteOutputPin, LOW);
  }

  pwmWriteDelayMicroseconds(pulseWidth);
}

void pwmWriteByte(uint8_t byte, uint16_t pulseWidth)  {
  uint8_t mask = 128; // send most significant bit first (network byte order)

  int i;
  for (i = 0; i < 8; i++) {
    uint8_t bit = byte & mask;

    pwmWriteBit(bit, pulseWidth);

    // shift the mask bit to right (towards least significant bit)
    mask = mask >> 1; 
  }
}

//
// exported functions
//
void pwmWriteSetup(int outputPin) {
  pwmWriteOutputPin = outputPin;

  pinMode(outputPin, OUTPUT);
}

arduinoPWMPacket pwmWrite(arduinoPWMPacket p) {
  int i, j;
  for (i = 0; i <= p.header.numRepeat; i++) { // numRepeat==0 -> send message once, numRepeat==1 -> twice, etc.
    for (j = 0; j < p.header.numBytes; j++) {
      pwmWriteByte(p.payload[j], p.header.pulseWidth);
    }

    pwmWriteDelayMicroseconds(p.header.pulseWidth*p.header.clearCount);
  }

  p.header.errorCode = 0;
  return p;
}
