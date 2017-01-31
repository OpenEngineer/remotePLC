#include <Arduino.h>

#include "arduinoPWMPacket.h"

namespace pwmRead {

// 
// internal functions and variables
//
int inputPin;

// for communication with the interrupt functions
volatile unsigned long lowPulseStartTime = 0;
volatile unsigned long lowPulseEndTime = 0;
volatile unsigned long highPulseStartTime = 0;
volatile unsigned long highPulseEndTime = 0;

volatile int pulseWidth;
volatile int pulseMargin;
volatile unsigned long minPulseWidth; // pulseWidth - pulseMargin

volatile int lowPulseCount = 0;
volatile int highPulseCount = 0;
volatile int pulsePairId = 0;

void pwmReadDetectLowPulseEnd() {
  // get data
  unsigned long t  = micros();
  unsigned long t0 = lowPulseStartTime;

  unsigned long dt    = pulseWidth;
  unsigned long et    = pulseMargin;
  unsigned long dtmin = minPulseWidth;

  unsigned long t1min = t0 + dtmin;

  if (t > t1min) {
    unsigned long t1 = t;
    int Dt = int(t1 - t0);
    int N  = (Dt + et)/dt;

    // set data
    highPulseStartTime = t;
    lowPulseEndTime    = t1;
    lowPulseCount      = N;
  }
}

void pwmReadDetectHighPulseEnd() {
  // get data
  unsigned long t  = micros();
  unsigned long t0 = highPulseStartTime;

  unsigned long dt    = pulseWidth;
  unsigned long et    = pulseMargin;
  unsigned long dtmin = minPulseWidth;

  int i = pulsePairId;

  unsigned long t1min = t0 + dtmin;


  if (t > t1min) {
    unsigned long t1 = t;
    int Dt = int(t1 - t0);
    int N = (Dt + et)/dt;

    i = (i+1)%32767;

    // set data
    lowPulseStartTime = t;
    highPulseEndTime  = t1;
    highPulseCount    = N;
    pulsePairId       = i;
  }
}

void pwmReadDetectHighLowPulses() {
  if (digitalRead(inputPin) != LOW) {
    pwmReadDetectLowPulseEnd();
  } else {
    pwmReadDetectHighPulseEnd();
  }
}

void pwmReadSetInterruptParameters(arduinoPWMPacket p) {
  pulseWidth = int(p.header1.pulseWidth);
  pulseMargin = int(p.header2.pulseMargin);
  minPulseWidth = pulseWidth - pulseMargin;
}

unsigned long pwmReadGetTimeOutDeadline(int timeOutCount) {
  long deltaTimeOut = long((long)pulseWidth*(long)timeOutCount -1L)/1000L + 1L;

  unsigned long deadline = millis() + (unsigned long)(deltaTimeOut);
  return deadline;
}

typedef struct {
  bool isStarted;
  bool isEnded;

  int byteI;
  int bitI;

  uint8_t byte;
  uint8_t mask;
} pwmReadState_t;

pwmReadState_t initReadState() {
  pwmReadState_t state;
  state.isStarted = false;
  state.isEnded = false;
  state.byteI = 0;
  state.bitI = 0;

  state.byte = 0;
  state.mask = 128;

  return state;
}

void pwmReadByte(pwmReadState_t *state, arduinoPWMPacket *p) { 
  if (state->byteI >= p->header1.numBytes) {
    state->isEnded = true;
    p->header1.errorCode = 0;
  } else {
    p->payload[state->byteI] = state->byte;
    state->byteI += 1;
    state->bitI = 0;
    state->byte = 0;
    state->mask = 128;
  }
}

void pwmReadBit(pwmReadState_t *state, arduinoPWMPacket *p) {
  state->bitI += 1;
  state->mask = state->mask >> 1;

  if (state->bitI == 8) {
    pwmReadByte(state, p);
  }
}

void pwmReadLowBit(pwmReadState_t *state, arduinoPWMPacket *p) {
  pwmReadBit(state, p);
}

void pwmReadHighBit(pwmReadState_t *state, arduinoPWMPacket *p) {
  state->byte = state->byte | state->mask;
  pwmReadBit(state, p);
}

bool pwmReadReadyToStart(int tmpLowPulseCount, int tmpHighPulseCount, arduinoPWMPacket p) {
  if (tmpLowPulseCount >= p.header2.clearCount &&
      tmpHighPulseCount < p.header2.clearCount) {
    return true;
  } else {
    return false;
  }
}
//
// exported functions
//
void pwmReadSetupUnoPin2() {
  inputPin = 2;
  int correspondingInterruptPin = 0;

  attachInterrupt(correspondingInterruptPin, pwmReadDetectHighLowPulses, CHANGE);
}

arduinoPWMPacket pwmRead(arduinoPWMPacket p) {
  pwmReadSetInterruptParameters(p);

  int prevPulsePairId = pulsePairId;

  pwmReadState_t readState = initReadState();
  arduinoPWMPacket answer = p;
  answer.header1.errorCode = 1;

  unsigned long deadline = pwmReadGetTimeOutDeadline(int(p.header2.timeOutCount));

  while (millis() <= deadline) {
    if (prevPulsePairId != pulsePairId) {
      // store volatile interrupt variables locally
      prevPulsePairId       = pulsePairId;
      int tmpLowPulseCount  = lowPulseCount;
      int tmpHighPulseCount = highPulseCount;
      
      if (!readState.isStarted) { 
        if (pwmReadReadyToStart(tmpLowPulseCount, tmpHighPulseCount, p)) {
          readState.isStarted = true; 
        } else {
          continue;
        }
      } 

      int i;
      if (!(readState.bitI==0 && readState.byteI==0)) { // this is not the first pulse (the first bit can never be 0)
        for (i = 0; i < tmpLowPulseCount; i++) {
          pwmReadLowBit(&readState, &answer);
        }
      }

      for (i = 0; i < tmpHighPulseCount; i++) {
        pwmReadHighBit(&readState, &answer);
      }

      if (readState.isEnded) {
        break;
      }
    } 
  }

  return answer;
}

} // end of pwmRead namespace
