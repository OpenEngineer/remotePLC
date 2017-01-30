#include <Arduino.h>

#include "arduinoPWMPacket.h"

#define PWM_READ_DEFAULT_PULSE_MARGIN 50 // microseconds

// 
// internal functions and variables
//
int pwmReadInputPin;

// for communication with the interrupt functions
typedef struct {
  unsigned long lowPulseStartTime;
  unsigned long lowPulseEndTime;
  unsigned long highPulseStartTime;
  unsigned long highPulseEndTime;
} pwmReadInterruptTimes_t;

typedef struct {
  int pulseWidth;
  int pulseMargin;
  unsigned long minPulseWidth; // pulseWidth - pulseMargin
} pwmReadInterruptSettings_t;

typedef struct {
  int lowPulseCount;
  int highPulseCount;
  int pulsePairId;
} pwmReadInterruptCounts_t;

volatile pwmReadInterruptTimes_t    pwmReadInterruptTimes;
volatile pwmReadInterruptSettings_t pwmReadInterruptSettings;
volatile pwmReadInterruptCounts_t   pwmReadInterruptCounts;

void pwmReadDetectLowPulseEnd() {
  // get data
  unsigned long t  = micros();
  unsigned long t0 = pwmReadInterruptTimes.lowPulseStartTime;

  unsigned long dt    = pwmReadInterruptSettings.pulseWidth;
  unsigned long et    = pwmReadInterruptSettings.pulseMargin;
  unsigned long dtmin = pwmReadInterruptSettings.minPulseWidth;

  unsigned long t1min = t0 + dtmin;

  if (t > t1min) {
    unsigned long t1 = t;
    int Dt = int(t1 - t0);
    int N  = (Dt + et)/dt;

    // set data
    pwmReadInterruptTimes.highPulseStartTime = t;
    pwmReadInterruptTimes.lowPulseEndTime    = t1;
    pwmReadInterruptCounts.lowPulseCount     = N;
  }
}

void pwmReadDetectHighPulseEnd() {
  // get data
  unsigned long t  = micros();
  unsigned long t0 = pwmReadInterruptTimes.highPulseStartTime;

  unsigned long dt    = pwmReadInterruptSettings.pulseWidth;
  unsigned long et    = pwmReadInterruptSettings.pulseMargin;
  unsigned long dtmin = pwmReadInterruptSettings.minPulseWidth;

  int i = pwmReadInterruptCounts.pulsePairId;

  unsigned long t1min = t0 + dtmin;


  if (t > t1min) {
    unsigned long t1 = t;
    int Dt = int(t1 - t0);
    int N = (Dt + et)/dt;

    i = (i%32767) + 1;

    // set data
    pwmReadInterruptTimes.lowPulseStartTime = t;
    pwmReadInterruptTimes.highPulseEndTime = t1;
    pwmReadInterruptCounts.highPulseCount = N;
    pwmReadInterruptCounts.pulsePairId = i;
  }
}

void pwmReadDetectHighLowPulses() {
  if (digitalRead(pwmReadInputPin) != LOW) {
    pwmReadDetectLowPulseEnd();
  } else {
    pwmReadDetectHighPulseEnd();
  }
}

void pwmReadSetInterruptParameters(arduinoPWMPacket p) {
  pwmReadInterruptSettings.pulseWidth = int(p.header.pulseWidth);
  pwmReadInterruptSettings.pulseMargin = PWM_READ_DEFAULT_PULSE_MARGIN;
  pwmReadInterruptSettings.minPulseWidth = pwmReadInterruptSettings.pulseWidth - pwmReadInterruptSettings.pulseMargin;
}

unsigned long pwmReadGetTimeOutDeadline(uint16_t pulseWidth, uint16_t timeOutCount) {
  long deltaTimeOut = ((long)pulseWidth*(long)timeOutCount -1L)/1000L + 1L;

  unsigned long deadline = millis() + (unsigned long)(deltaTimeOut);
  return deadline;
}

typedef struct {
  bool isStarted;
  bool isEnded;
  int errorCode;

  int byteI;
  int bitI;

  uint8_t byte;
  uint8_t mask;
} pwmReadState_t;

void pwmReadByte(pwmReadState_t *state, arduinoPWMPacket *p) { 
  p->payload[state->byteI] = state->byte;
  state->byteI += 1;
  state->bitI = 0;
  state->byte = 0;
  state->mask = 128;

  if (state->byteI >= p->header.numBytes) {
    state->isEnded = true;
    p->header.errorCode = 0;
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

bool pwmReadReadyToStart(pwmReadInterruptCounts_t interruptState, arduinoPWMPacket p) {
  if (interruptState.lowPulseCount >= p.header.clearCount &&
      interruptState.highPulseCount < p.header.clearCount) {
    return true;
  } else {
    return false;
  }
}
//
// exported functions
//
void pwmReadSetupUnoPin2() {
  pwmReadInputPin = 2;
  int correspondingInterruptPin = 0;

  attachInterrupt(correspondingInterruptPin, pwmReadDetectHighLowPulses, CHANGE);
}

arduinoPWMPacket pwmRead(arduinoPWMPacket p) {
  pwmReadSetInterruptParameters(p);

  unsigned long deadline = pwmReadGetTimeOutDeadline(p.header.pulseWidth, p.header.timeOutCount);

  pwmReadInterruptCounts_t interruptState;
  interruptState.lowPulseCount = pwmReadInterruptCounts.lowPulseCount;
  interruptState.highPulseCount= pwmReadInterruptCounts.highPulseCount;
  interruptState.pulsePairId = pwmReadInterruptCounts.pulsePairId;
  int prevPulsePairId = interruptState.pulsePairId;

  pwmReadState_t readState;
  arduinoPWMPacket answer = p;

  while (millis() < deadline) {
    if (prevPulsePairId != interruptState.pulsePairId) {
      // store volatile interrupt variables locally
      interruptState.lowPulseCount = pwmReadInterruptCounts.lowPulseCount;
      interruptState.highPulseCount= pwmReadInterruptCounts.highPulseCount;
      interruptState.pulsePairId = pwmReadInterruptCounts.pulsePairId;
      prevPulsePairId = interruptState.pulsePairId;
      
      if (!readState.isStarted && pwmReadReadyToStart(interruptState, p)) {
        readState.isStarted = true;
      } else {
        continue;
      }

      int i;
      if (!(readState.bitI==0 && readState.byteI==0)) { // this is not the first pulse (the first bit can never be 0)
        for (i = 0; i < interruptState.lowPulseCount; i++) {
          pwmReadLowBit(&readState, &answer);
        }
      }

      for (i = 0; i < interruptState.highPulseCount; i++) {
        pwmReadHighBit(&readState, &answer);
      }

      if (readState.isEnded) {
        break;
      }
    }
  }

  return answer;
}
