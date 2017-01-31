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

void serialSync() {
  uint8_t syncByte1;
  uint8_t syncByte2;

  while (Serial.available() > 0) {
    Serial.readBytes((char*)&syncByte1, sizeof(char));
    if (syncByte1 == ARDUINOPWM_SYNCBYTE1) {
      Serial.readBytes((char*)&syncByte2, sizeof(char));
      if (syncByte2 == ARDUINOPWM_SYNCBYTE2) {
        break;
      }
    }
  }
}

arduinoPWMHeader1 serialReadHeader1() {
  Serial.readBytes((char*)serialBuffer, sizeof(arduinoPWMHeader1));

  // fill a packet with the incoming bytes
  arduinoPWMHeader1 h;

  h.opCode       = serialBuffer[0];
  h.numBytes     = serialBuffer[1];
  h.pulseWidth   = serialTwoBytesToUint16(serialBuffer[2], serialBuffer[3]);
  h.errorCode    = serialBuffer[4];

  return h;
}

arduinoPWMHeader2 serialReadHeader2() {
  Serial.readBytes((char*)serialBuffer, sizeof(arduinoPWMHeader2));

  arduinoPWMHeader2 h;

  h.clearCount   = serialBuffer[0];
  h.timeOutCount = serialTwoBytesToUint16(serialBuffer[1], serialBuffer[2]);
  h.numRepeat    = serialBuffer[3];
  h.pulseMargin  = serialBuffer[4];

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

void serialWriteSyncBytes() {
  serialBuffer[0] = ARDUINOPWM_SYNCBYTE1;
  serialBuffer[1] = ARDUINOPWM_SYNCBYTE2;
}

void serialWriteHeader1(arduinoPWMPacket p) { 
  int offset = ARDUINOPWM_NUMSYNCBYTES;
  serialBuffer[offset+0] = p.header1.opCode;
  serialBuffer[offset+1] = p.header1.numBytes;
  serialBuffer[offset+2] = serialUint16ToByte1(p.header1.pulseWidth);
  serialBuffer[offset+3] = serialUint16ToByte2(p.header1.pulseWidth);
  serialBuffer[offset+4] = p.header1.errorCode;
}

void serialWriteHeader2(arduinoPWMPacket p) { 
  int offset = ARDUINOPWM_NUMSYNCBYTES + sizeof(arduinoPWMHeader1);
  serialBuffer[offset+0] = p.header2.clearCount;
  serialBuffer[offset+1] = serialUint16ToByte1(p.header2.timeOutCount);
  serialBuffer[offset+2] = serialUint16ToByte2(p.header2.timeOutCount);
  serialBuffer[offset+3] = p.header2.numRepeat;
  serialBuffer[offset+4] = p.header2.pulseMargin;
}

void serialWritePayload(arduinoPWMPacket p) {
  int i;
  int offset = ARDUINOPWM_NUMSYNCBYTES + sizeof(arduinoPWMHeader1) + sizeof(arduinoPWMHeader2);

  for (i = 0; i < p.header1.numBytes; i++) {
    serialBuffer[offset+i] = p.payload[i];
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

  serialSync();

  p.header1 = serialReadHeader1();

  p.header2 = serialReadHeader2();

  p.payload = serialReadPayload(p.header1.numBytes);

  return p;
}

void serialWriteMessage(arduinoPWMPacket p) {
  serialWriteSyncBytes();

  // write the header to the serialBuffer
  serialWriteHeader1(p);
  serialWriteHeader2(p);

  // write the payload to the serialBuffer
  serialWritePayload(p);

  // send the content of the serialBuffer
  serialSend(p.header1.numBytes + ARDUINOPWM_NUMSYNCBYTES + sizeof(arduinoPWMHeader1) + sizeof(arduinoPWMHeader2));
}
