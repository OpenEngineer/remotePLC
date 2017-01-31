#include <Arduino.h>

#include "arduinoPWMPacket.h"

//
// internal functions and variables
//
uint8_t serialBuffer[ARDUINO_PWM_PACKET_MAX_PACKET_SIZE]; // includes payload AND header

uint16_t serialTwoBytesToUint16(uint8_t byte1, uint8_t byte2) {
  uint16_t result = uint16_t(byte1) + 256*(uint16_t(byte2));
  return result;
}

uint8_t serialUint16ToByte1(uint16_t i) {
  uint8_t byte1 = i%256;
  return byte1;
}

uint8_t serialUint16ToByte2(uint16_t i) {
  uint8_t byte1 = serialUint16ToByte1(i);
  uint16_t byte2 = (i - uint16_t(byte1))/256; // make sure all this arithmatic takes place in uint16_t
  return uint8_t(byte2);
}

arduinoPWMHeader serialReadHeader() {
  Serial.readBytes((char*)serialBuffer, sizeof(arduinoPWMHeader));

  // fill a packet with the incoming bytes
  arduinoPWMHeader h;

  h.opCode       = serialBuffer[0];
  h.numBytes     = serialBuffer[1];
  h.pulseWidth   = serialTwoBytesToUint16(serialBuffer[2], serialBuffer[3]);
  h.clearCount   = serialBuffer[4];
  h.timeOutCount = serialTwoBytesToUint16(serialBuffer[5], serialBuffer[6]);
  h.numRepeat    = serialBuffer[7];
  h.errorCode    = serialBuffer[8];

  return h;
}

uint8_t *serialReadPayload(int payloadSize) {
  Serial.readBytes((char*)serialBuffer, payloadSize);

  // use the arduinoPWMBuffer as a destination
  uint8_t *dst = getArduinoPWMBuffer();

  int i;
  for (i = 0; i < payloadSize; i++) {
    dst[i] = serialBuffer[i];
  }

  return dst;
}

void serialWriteHeader(arduinoPWMPacket p) { 
  serialBuffer[0] = p.header.opCode;
  serialBuffer[1] = p.header.numBytes;
  serialBuffer[2] = serialUint16ToByte1(p.header.pulseWidth);
  serialBuffer[3] = serialUint16ToByte2(p.header.pulseWidth);
  serialBuffer[4] = p.header.clearCount;
  serialBuffer[5] = serialUint16ToByte1(p.header.timeOutCount);
  serialBuffer[6] = serialUint16ToByte2(p.header.timeOutCount);
  serialBuffer[7] = p.header.numRepeat;
  serialBuffer[8] = p.header.errorCode;
}

void serialWritePayload(arduinoPWMPacket p) {
  int i;

  for (i = 0; i < p.header.numBytes; i++) {
    serialBuffer[i+sizeof(arduinoPWMHeader)] = p.payload[i];
  }
}

void serialSend(int numTotalBytes) {
  Serial.write(serialBuffer, numTotalBytes);
}

//
// exported functions
//
void serialSetup(int bitRate) {
  Serial.begin(bitRate, SERIAL_8N2);
}

bool serialReadWriteIsReady() {
  if (Serial.available() > 0) {
    return true;
  } else {
    return false;
  }
}

arduinoPWMPacket serialReadMessage() {
  arduinoPWMPacket p;

  p.header = serialReadHeader();

  p.payload = serialReadPayload(p.header.numBytes);

  return p;
}

void serialWriteMessage(arduinoPWMPacket p) {
  // write the header to the serialBuffer
  serialWriteHeader(p);

  // write the payload to the serialBuffer
  serialWritePayload(p);

  // send the content of the serialBuffer
  serialSend(p.header.numBytes + sizeof(arduinoPWMHeader));
}
