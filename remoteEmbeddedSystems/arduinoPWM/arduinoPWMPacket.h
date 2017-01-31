#ifndef arduinoPWMPacket_h
#define arduinoPWMPacket_h

#include <Arduino.h>

#define ARDUINO_PWM_PACKET_MAX_PAYLOAD_SIZE 255
#define ARDUINO_PWM_PACKET_MAX_PACKET_SIZE 267
#define ARDUINO_PWM_ERROR_OPCODE_NOT_RECOGNIZED 3
#define ARDUINO_PWM_OPCODE_WRITE 1
#define ARDUINO_PWM_OPCODE_READ 2
#define ARDUINOPWM_SYNCBYTE1 172
#define ARDUINOPWM_SYNCBYTE2 86
#define ARDUINOPWM_NUMSYNCBYTES 2

typedef struct {
  uint8_t opCode;
  uint8_t numBytes;
  int     pulseWidth;
  uint8_t errorCode;
} arduinoPWMHeader1;

typedef struct {
  uint8_t clearCount;
  int     timeOutCount;
  uint8_t numRepeat;
  uint8_t pulseMargin;
} arduinoPWMHeader2;

typedef struct {
  uint8_t syncBytes[ARDUINOPWM_NUMSYNCBYTES];
  arduinoPWMHeader1 header1;
  arduinoPWMHeader2 header2;

  uint8_t *payload;
} arduinoPWMPacket;

uint8_t *getArduinoPWMBuffer();

#endif
