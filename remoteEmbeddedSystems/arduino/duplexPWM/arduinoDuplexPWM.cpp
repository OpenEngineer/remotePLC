# include <Arduino.h>

// parameters:
#define OUTPUT_PIN 8 
#define INPUT_PIN 2
#define INPUT_INTERRUPT_PIN 0 // 0 is the interrupt id of actual inputPin 2
#define DEBUG // comment out for no debugging
#define BUFFER_SIZE 263 // 255 max payload + 8 bytes header
#define WRITE_OPCODE 1
#define READ_OPCODE 2
#define HEADER_SIZE 9 // 9 bytes
#define MIN_SAMPLE_PERIOD 30 // number of Microseconds, determines rate at which the inputPin is sample
#define MAX_SAMPLES_PER_PULSE 40
#define SERIAL_BUFFER_SIZE 64
#define MARGIN 30 // default margin is 50micros

// write incoming serial messages to this pin:
int outputPin = OUTPUT_PIN; 

// read this pin and send to serial:
int inputPin = INPUT_PIN;  

// # buffer related functions and global declaration of buffer itself
// both the incoming serial messages (written to output pin),
//  and the outgoing serial messages (read from input pin),
//  are passed via this buffer to the WriteOutputBytes() function and
//  from the ReadInput function
uint8_t buffer[BUFFER_SIZE];
char serialResetBuffer[64];

// shift buffer to right
void ShiftBuffer(int numBytes, int numShift) {
  int i;
  for (i = numBytes; i >= 0; i--) {
    if (i+numShift < BUFFER_SIZE) {
      buffer[i+numShift] = buffer[i];
    }
  }
}

//
void delayMicrosecondsAccurate(int d) {
  int dMicros = d%1000;
  int dMillis = (d-dMicros)/1000;

  delayMicroseconds(dMicros);

  if (dMillis > 0) {
    delay(dMillis);
  }
}

// mostly used for headers
void PrependToBuffer(int numBytesTotal, int numBytesHeader, uint8_t *header) {
  // numBytes by numShift amount to the right
  ShiftBuffer(numBytesTotal, numBytesHeader);

  int i;
  for (i = 0; i < numBytesHeader; i++) {
    buffer[i] = header[i];
  }
}


// # WriteOutput related functions:

// write a bit, then wait a numer of microseconds (pulseWidth)
void WriteOutputBit(uint8_t bit, int pulseWidth) {
  if (bit != 0) {
    digitalWrite(outputPin, HIGH);
  } else {
    digitalWrite(outputPin, LOW);
  }

  delayMicrosecondsAccurate(pulseWidth);
}

// write a single byte, in char form, to the output pin
//  by cycling through its 8 bits.
// we use a shifting mask for this
void WriteOutputByte(uint8_t byte, int pulseWidth) {
  uint8_t mask = 128; // send most significant bit first (network byte order)

  int i;
  for (i = 0; i < 8; i++) {
    uint8_t bit = byte & mask;

    WriteOutputBit(bit, pulseWidth);



    // shift the mask bit to right (towards least siginificant bit)
    mask = mask >> 1; 
  }
}

// inputs to the WriteOutputBytes() function:
//  - numBytes: number of bytes to convert to bits and write to the output pin.
//     this number must be smaller than the size of the buffer
//  - pulseWidth: in microseconds, the width of a single bit
// if we want to repeat the message, then it is recommended to this upstream
//  by sending the extended byte stream via the serial line. So repeating is 
//  something that the client needs to handle
void WriteOutputBytes(int numBytes, int pulseWidth, int clearCount, int numRepeat) {
  int i, j;
  for (i = 0; i <= numRepeat; i++) { // numRepeat ==0 -> send message once, numRepeat==1->twice etc.
    for (j = 0; j < numBytes; j++) {
      WriteOutputByte(buffer[j], pulseWidth);
    }

    delayMicrosecondsAccurate(pulseWidth*clearCount);
  }
}


// # ReadInputBytes related functions:
int READ_DESYNC = 0;

int calcDesync(int highCount, int highPosCount, int lowCount, int lowPosCount) {
  int highPos = highPosCount/highCount;
  int lowPos = lowPosCount/lowCount;

  int desync;
  if (lowCount < highCount) {
    if (lowPos < highPos) {
      desync = -lowCount;
    } else {
      desync = lowCount;
    }
  } else { // highCount < lowCount
    if (highPos < lowPos) {
      desync = -highPos;
    } else {
      desync = highPos;
    }
  }

  return desync;
}

volatile unsigned long START_HIGH = 0;
volatile unsigned long START_LOW = 0;
volatile unsigned long END_HIGH = 0;
volatile unsigned long END_LOW = 0;
volatile int PULSE_WIDTH = 0;
volatile int PULSE_COUNT_LOW = 0;
volatile int PULSE_COUNT_HIGH = 0;
volatile int PULSE_ID = 0;

// interrupt functions
void sampleRising() {
  volatile unsigned long endLow = micros();

  if (endLow > START_LOW + (unsigned long)(PULSE_WIDTH - MARGIN)) { // otherwise ignore (assume that we were always in a high state)
#ifdef DEBUG
    //digitalWrite(outputPin, HIGH);
#endif
    END_LOW = endLow;
    START_HIGH = endLow;
    int diff = int(END_LOW - START_LOW);
    PULSE_COUNT_LOW = (diff + MARGIN)/PULSE_WIDTH; // with some margin
  }
}

// only change the pulse id if a high pulse has been detected for at least half the pulseWidth
void sampleFalling() {
  volatile unsigned long endHigh = micros();

  if (endHigh > START_HIGH + (unsigned long)(PULSE_WIDTH - MARGIN)) { // otherwise ignore (assume that we were always in a low state)
#ifdef DEBUG
    //digitalWrite(outputPin, LOW);
#endif
    END_HIGH = endHigh;
    START_LOW = endHigh;
    int diff = int(END_HIGH - START_HIGH);
    PULSE_COUNT_HIGH = (diff + MARGIN)/PULSE_WIDTH; // the actual pulse can be MARGIN shorter than the expected pulse (but no less)

    PULSE_ID = (PULSE_ID%32767) + 1;
  }
}

// combination of rising and falling, because apparantly only one function can be attached to an interruptpin at a time
void sampleChange() {
  if (digitalRead(inputPin) == HIGH) { // rising, end a low pulse, start a high pulse
    sampleRising();
  } else { // falling, end a high pulse, start a low pulse (ends the pair low/high)
    sampleFalling();
  }
}

// pulseWidth in microSeconds
// return true for HIGH and false for LOW
bool ReadInputBit(int pulseWidth) {
  // keep a desync count inside
  
  // the the sampling rate
  int samplePeriod = MIN_SAMPLE_PERIOD;
  int samplesPerPulse = pulseWidth/MIN_SAMPLE_PERIOD;
  if (samplesPerPulse > MAX_SAMPLES_PER_PULSE) {
    samplesPerPulse = MAX_SAMPLES_PER_PULSE;
    samplePeriod = pulseWidth/samplesPerPulse;
  }

  // set the counters
  int count;
  int highCount = 0;
  int lowCount = 0;
  int highPosCount = 0;
  int lowPosCount = 0;

  // count the number of high and low reads, as well as the cumulative high and low positions
  for (count=READ_DESYNC; count<samplesPerPulse; count++){
    if (digitalRead(inputPin) == HIGH) {
      highCount += 1;
      highPosCount += count;
    } else {
      lowCount += 1;
      lowPosCount += count;
    }

    // now delay a little
    delayMicroseconds(samplePeriod); // will definitely be smaller than 1000
  }

  // handle the average position counts in order to determine the desycn
  READ_DESYNC = calcDesync(highCount, highPosCount, lowCount, lowPosCount);
  
  // if the desync is negative, perhaps we should just wait this long and set it to zero
  if (READ_DESYNC < 0) {
    delayMicroseconds(-samplePeriod*READ_DESYNC);
    READ_DESYNC = 0;
  }

  bool isHigh;
  if (highCount >= lowCount) {
    isHigh = true;
  } else {
    isHigh = false;
  }

  return isHigh;
}

// in case of timeOutCount==0, the function returns immediately.
//  this can be used to read constant states (ie. not time varying)
int WaitForClearInput(int numBytes, int pulseWidth, int clearCount, int timeOutCount) {
  int count = 0; 
  int pulseCount = 0;

  while(count < clearCount && pulseWidth < timeOutCount) {
    count = ReadInputBit(pulseWidth) == true ?  0 : count + 1;

    pulseCount += 1;
  }

  return pulseCount;
}


// the first bit of the first byte will always be 0x1
int ReadFirstInputByte(int pulseWidth, int timeOutCount, int pulseCount, int byteI) {
  while (ReadInputBit(pulseWidth)==false && pulseCount < timeOutCount) {
    pulseCount += 1;
  }

  if (pulseCount < timeOutCount) {
    uint8_t byte = 128;
    uint8_t mask = 128; // first pulse is always high

    int i;
    for (i = 1; i < 8; i++) {
      mask = mask >> 1;

      // add a bit to the byte if the pulse is high
      if (ReadInputBit(pulseWidth)==true) {
        byte = byte | mask;
      }
      
      pulseCount += 1;
    }

    // finally put the byte into the buffer
    buffer[byteI] = byte;
  } else {
    buffer[byteI] = 0;
  }

  return pulseCount;
}

int ReadInputByte(int pulseWidth, int byteI) {
  int pulseCount = 0;

  uint8_t byte = 0;
  uint8_t mask = 128;

  int i;
  for (i = 0; i < 8; i++) {

    // add a bit to the byte if the pulse is high
    if (ReadInputBit(pulseWidth)==true) {
      byte = byte | mask;
    }

    pulseCount += 1;

    mask = mask >> 1;
  }

  // finally put the byte into the buffer
  buffer[byteI] = byte;

  return pulseCount;
}

unsigned long TIME_OUT_END;

// pulseWidth in microSeconds
void startTimeOut(int pulseWidth, int timeOutCount) {
  long diff = ((long)pulseWidth*(long)timeOutCount-1L)/1000L+1L;
  TIME_OUT_END = millis() + (unsigned long)(diff);
}

bool isTimeOut() {
  unsigned long currentTime = millis();

  if (currentTime > TIME_OUT_END) {
    return true;
  } else {
    return false;
  }
}

bool recordBitsIf8(int *bitI, uint8_t *byte, uint8_t *mask, int *byteI, int numBytes) {
  bool isEnd = false;
  if (*bitI == 8) {
    buffer[*byteI] = *byte;
    *byteI += 1;
    *bitI = 0;
    *mask = 128;
    *byte = 0;

    if (*byteI >= numBytes) { 
      isEnd = true;
    }
  }

  return isEnd;
}
// inputs to the ReadInputBytes() function:
//  - numBytes: fill the global buffer with this number of bytes,
//     only then is the reading considered a success.
//     the incoming message is numBytes*pulseWidth long
//  - pulseWidth: in microseconds, the width of a single bit
//  - timeOutCount: multiply pulseWidth by this number to get the timeOut time. 
//     stop trying to read input after this time.
//     this requires counting the pulses and comparing to this number
// outputs of the ReadInputBytes() function:
//  buffer[]: parsed bytes are saved into the global buffer
//  bool return value: true for succes, false for failure
//   in case of failure the buffer can contain an incomplete message, 
//   the downstream function should then ignore this
int ReadInputBytes(int numBytes, int pulseWidth, int clearCount, int timeOutCount) {
  // set the interrupt variables right
  PULSE_WIDTH = pulseWidth;

  // start the timeOut counter
  startTimeOut(pulseWidth, timeOutCount);

  // state variables and counters
  int prevPulseId = PULSE_ID;
  int errorCode = 1; // assume failed
  int byteI = 0;
  int bitI = 0;
  uint8_t byte = 0;
  uint8_t mask = 128;

  bool isStarted = false;
  bool isFirstPulse = true;

  int i; // used in the for loops when generating bits and bytes

  // TODO: refactor 
  while (!isTimeOut()) {
    if (PULSE_ID != prevPulseId) { // only handle the low/high pulses if they are different from the last pair
      // local copies of interrupts state variables (the interrupts can 
      //   otherwise change these variables during the processing below)
      prevPulseId = PULSE_ID;
      int pulseCountLow = PULSE_COUNT_LOW;
      int pulseCountHigh = PULSE_COUNT_HIGH;


      // detect if this is the first pulse of a message
      //  TODO: more robust criteria
      if (!isStarted && (pulseCountLow >= clearCount && pulseCountHigh < clearCount)) {
        isStarted = true;
      } else if (!isStarted) {
        // wait until the pulseCount of the first low pulse is large enough:
        continue;
      }

      // take the low pulses into account if this is not the first pulse (for the 0's in the bytes)
      if (!isFirstPulse) {
        for (i = 0; i < pulseCountLow; i++) {
          bitI += 1;
          mask = mask >> 1;

          if (recordBitsIf8(&bitI, &byte, &mask, &byteI, numBytes) == true) {
            errorCode = 0;
            break;
          }
        }
      } else {
        // this is the first pulse, from now on allow the processing the pulseCountLow into bits
        isFirstPulse = false;
      }

      for (i = 0; i < pulseCountHigh; i++) {
        byte = byte | mask;
        bitI += 1;
        mask = mask >> 1;

        // record and reset bit state
        if (recordBitsIf8(&bitI, &byte, &mask, &byteI, numBytes) == true) {
          errorCode = 0;
          break;
        }
      }

      if (byteI >= numBytes) { // an extra check
        errorCode = 0;
        break;
      }
    } // end of change detection
  } // end of while

  // write the last byte, even though it might be incomplete
  if (byteI < numBytes) {
    buffer[numBytes-1] = byte;
  }

  return errorCode;
}


// # Serial message related functions:

// # a message has the following structure
// header (7 bytes):
//  - message type (1 byte) : WRITE_OPCODE for a write instruction
//                              READ_OPCODE for a read instruction
//  - numBytes, between 0 and 255 (1 byte)
//  - pulseWidth int (2bytes/16 bits)
//  - clearCount uint8_t (1 byte): for READ_OP, start recording after the 
//      inputPin is in a low state for this many pulses
//  - timeOutCount int16 (only for the read instruction, otherwise ignored)
//  - errorCode (1 byte): used in replies
// body:
//  - remaining bytes are written/read from/to the pins

// TODO: are there arduino library functions that do this better?
int twoBytesToInt(uint8_t i0, uint8_t i1) {
  int i = int(i0) + int(i1)*256;
  return i;
}

uint8_t intToFirstByte(int i) {
  // i = i0 + i1*256
  uint8_t i0 = i%256;
  return i0;
}

uint8_t intToSecondByte(int i) {
  // i = i0 + i1*256
  uint8_t i1 = (i-intToFirstByte(i))/256;
  return i1;
}

// serialization into header bytes
void SetHeader(int opCode, int numBytes, int pulseWidth, int clearCount, int timeOutCount, int numRepeat, int errorCode, uint8_t header[HEADER_SIZE]) {
  header[0] = uint8_t(opCode);

  header[1] = uint8_t(numBytes);

  header[2] = intToFirstByte(pulseWidth);
  header[3] = intToSecondByte(pulseWidth);

  header[4] = uint8_t(clearCount);

  header[5] = intToFirstByte(timeOutCount);
  header[6] = intToSecondByte(timeOutCount);

  header[7] = uint8_t(numRepeat);

  header[8] = uint8_t(errorCode);
}

// deserialization of header bytes
void GetHeader(uint8_t header[HEADER_SIZE], int *opCode, int *numBytes, int *pulseWidth, int *clearCount, int *timeOutCount, int *numRepeat, int *errorCode) {
  *opCode = int(header[0]);
  *numBytes = int(header[1]);
  *pulseWidth = twoBytesToInt(header[2], header[3]);
  *clearCount = int(header[4]);
  *timeOutCount = twoBytesToInt(header[5], header[6]);
  *numRepeat = int(header[7]);
  *errorCode = int(header[8]);
}

// HandleMessage() inputs:
//  - the timeOutCount is ignored in the case of the write instruction opCode (WRITE_OPCODE)
// do nothing in case the opCode isn't recognized
//  - errorCode: mostly dummy input (overwritten internally), but can be used as default value in replies
void HandleMessage(int opCode, int numBytes, int pulseWidth, int clearCount, int timeOutCount, int numRepeat, int errorCode) {
  // after every message a reply needs to be sent upstream
  //  this is to assure that everything is synchronous
  
  switch(opCode) {
    case WRITE_OPCODE: {
      WriteOutputBytes(numBytes, pulseWidth, clearCount, numRepeat);

      // assume always success
      // set the reply variables
      errorCode = 0;
    } break;
    case READ_OPCODE: {
      errorCode = ReadInputBytes(numBytes, pulseWidth, clearCount, timeOutCount);
    } break;
    default:
      // do nothing in case the opCode isn't recognized
      break;
  }

  // send the reply message via the serial line
  // the reply contains the same header as the received message
  uint8_t header[HEADER_SIZE];

  SetHeader(opCode, numBytes, pulseWidth, clearCount, timeOutCount, numRepeat, errorCode, header);

  PrependToBuffer(numBytes, HEADER_SIZE, header);

  Serial.write(buffer, HEADER_SIZE + numBytes);
}

// parse incoming messages sequentially
// responses are handled with HandleMessage
void ParseSerialMessages() {
  while (Serial.available() > 0) {
    // read the header 
    Serial.readBytes((char*)buffer, HEADER_SIZE);

    int opCode, numBytes, pulseWidth, clearCount, timeOutCount, numRepeat, errorCode; // the errorCode is unused for incoming messages, so this is a dummy variable
    GetHeader(buffer, &opCode, &numBytes, &pulseWidth, &clearCount, &timeOutCount, &numRepeat, &errorCode);

    // read the body of the message
    Serial.readBytes((char*)buffer, numBytes);

    // this function handles downstream, and also upstream reply:
    // use the incoming errorCode as default value (mostly 0)
    HandleMessage(opCode, numBytes, pulseWidth, clearCount, timeOutCount, numRepeat, errorCode);

    // after handling a single message, we throw remaining bytes away
    //Serial.flush(); 

    // reset the remaining serial buffer if it is overflowing
    // TODO: do the syncronization with start bytes
    while (Serial.available() >= 64) {
      Serial.readBytes(serialResetBuffer, 64);
    }
  }
}

void loop() {
  ParseSerialMessages();
}

// the slowest baudrate (9600 bps) for robustness
void setup() {
  Serial.begin(9600, SERIAL_8N2);

  // for WRITE_OP
  pinMode(outputPin, OUTPUT);

  // for READ_OP
  attachInterrupt(INPUT_INTERRUPT_PIN, sampleChange, CHANGE);
}

// program entry point
int main(void) {
  init();
  setup();
  for (;;) {
    loop();
  }
}
