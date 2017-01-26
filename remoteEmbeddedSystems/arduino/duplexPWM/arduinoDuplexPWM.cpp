# include <Arduino.h>

// parameters:
#define OUTPUT_PIN 8 
#define INPUT_PIN 12
#define DEBUG // comment out for no debugging
#define BUFFER_SIZE 263 // 255 max payload + 8 bytes header
#define WRITE_OPCODE 1
#define READ_OPCODE 2
#define HEADER_SIZE 8 // 8 bytes
#define MIN_SAMPLE_PERIOD 10 // number of Microseconds, determines rate at which the inputPin is sample
#define MAX_SAMPLES_PER_PULSE 20
#define SERIAL_BUFFER_SIZE 64

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
void WriteOutputBytes(int numBytes, int pulseWidth) {
  int i;
  for (i = 0; i < numBytes; i++) {
    WriteOutputByte(buffer[i], pulseWidth);
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
  int errorCode = 1; // assume failed
  int pulseCount = 0;


  // using Nyquist theory we know we need to sample each pulse at 
  //  least twice in order to get all the information
  
  // the inputPin must be in a low state for the length of a message-2.
  //  -2 because the message must be bounded by at least 2 high states
  //  for practical reasons we just use the length of the message (so we don't need handling of numBytes=0, 1 or 2)
  // waiting this long assures that the next high state we read is from the start of a message,
  //  not somewhere halfway
  pulseCount = WaitForClearInput(numBytes, pulseWidth, clearCount, timeOutCount);

  int byteI = 0;

  // ReadFirstInputByte() waits for the high state, and only then starts sampling
  pulseCount += ReadFirstInputByte(pulseWidth, timeOutCount, pulseCount, byteI);
  byteI += 1;

  while (pulseCount <= timeOutCount) {

    pulseCount += ReadInputByte(pulseWidth, byteI);
    byteI += 1;

    if (byteI >= numBytes) { // we managed to read all the bytes we needed
      errorCode = 0; // success
      break;
    }
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
void SetHeader(int opCode, int numBytes, int pulseWidth, int clearCount, int timeOutCount, int errorCode, uint8_t header[HEADER_SIZE]) {
  header[0] = uint8_t(opCode);

  header[1] = uint8_t(numBytes);

  header[2] = intToFirstByte(pulseWidth);
  header[3] = intToSecondByte(pulseWidth);

  header[4] = uint8_t(clearCount);

  header[5] = intToFirstByte(timeOutCount);
  header[6] = intToSecondByte(timeOutCount);

  header[7] = uint8_t(errorCode);
}

// deserialization of header bytes
void GetHeader(uint8_t header[HEADER_SIZE], int *opCode, int *numBytes, int *pulseWidth, int *clearCount, int *timeOutCount, int *errorCode) {
  *opCode = int(header[0]);
  *numBytes = int(header[1]);
  *pulseWidth = twoBytesToInt(header[2], header[3]);
  *clearCount = int(header[4]);
  *timeOutCount = twoBytesToInt(header[5], header[6]);
  *errorCode = int(header[7]);
}

// HandleMessage() inputs:
//  - the timeOutCount is ignored in the case of the write instruction opCode (WRITE_OPCODE)
// do nothing in case the opCode isn't recognized
//  - errorCode: mostly dummy input (overwritten internally), but can be used as default value in replies
void HandleMessage(int opCode, int numBytes, int pulseWidth, int clearCount, int timeOutCount, int errorCode) {
  // after every message a reply needs to be sent upstream
  //  this is to assure that everything is synchronous
  
  switch(opCode) {
    case WRITE_OPCODE: {
      WriteOutputBytes(numBytes, pulseWidth);

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

  SetHeader(opCode, numBytes, pulseWidth, clearCount, timeOutCount, errorCode, header);

  PrependToBuffer(numBytes, HEADER_SIZE, header);

  Serial.write(buffer, HEADER_SIZE + numBytes);
}

// parse incoming messages sequentially
// responses are handled with HandleMessage
void ParseSerialMessages() {
  while (Serial.available() > 0) {
    // read the header 
    Serial.readBytes((char*)buffer, HEADER_SIZE);

    int opCode, numBytes, pulseWidth, clearCount, timeOutCount, errorCode; // the errorCode is unused for incoming messages, so this is a dummy variable
    GetHeader(buffer, &opCode, &numBytes, &pulseWidth, &clearCount, &timeOutCount, &errorCode);

    // read the body of the message
    Serial.readBytes((char*)buffer, numBytes);

    // this function handles downstream, and also upstream reply:
    // use the incoming errorCode as default value (mostly 0)
    HandleMessage(opCode, numBytes, pulseWidth, clearCount, timeOutCount, errorCode);

    // after handling a single message, we throw remaining bytes away
    //Serial.flush(); 

    // reset the remaining serial buffer if it is overflowing
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
  pinMode(outputPin, OUTPUT);
  pinMode(inputPin , INPUT);
}

// program entry point
int main(void) {
  init();
  setup();
  for (;;) {
    loop();
  }
}
