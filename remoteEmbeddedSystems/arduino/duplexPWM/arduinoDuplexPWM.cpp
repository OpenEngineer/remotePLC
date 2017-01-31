#include <Arduino.h>

#include "arduinoPWMPacket.h"
#include "pwmRead.h"
#include "pwmWrite.h"
#include "serialReadWrite.h"

// parameters:
#define OUTPUT_PIN 8 
#define SERIAL_BITRATE 9600

arduinoPWMPacket handleMessage(arduinoPWMPacket question) {
  // after every message a reply needs to be sent upstream.
  // this is to assure that everything is synchronous
  arduinoPWMPacket answer;

  switch(question.header.opCode) {
    case ARDUINO_PWM_OPCODE_WRITE: {
      answer = pwmWrite(question);
    } break;
    case ARDUINO_PWM_OPCODE_READ: {
      answer = pwmRead::pwmRead(question);
    } break;
    default:
      answer.header.errorCode = ARDUINO_PWM_ERROR_OPCODE_NOT_RECOGNIZED;
      break;
  }

  return answer;
}


void setup() {
  // use the slowest baudrate (9600 bps) for robustness
  serialSetup(SERIAL_BITRATE);

  pwmWriteSetup(OUTPUT_PIN);

  pwmRead::pwmReadSetupUnoPin2();
}

void loop() {
  if (serialReadWriteIsReady()) {
    arduinoPWMPacket question = serialReadMessage();

    arduinoPWMPacket answer = handleMessage(question);
    
    serialWriteMessage(answer);
  }
}

// program entry point
int main(void) {
  init();
  setup();
  for (;;) {
    loop();
  }
}
