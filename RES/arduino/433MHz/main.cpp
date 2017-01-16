# include <Arduino.h>
// 433MHz example

// serial tx and rx setup
int transPin = 13;
int receivePin = 12;

void setup() {
  Serial.begin(9600, SERIAL_8N2);
  pinMode(transPin, OUTPUT);
  pinMode(receivePin, INPUT);
}

void serial2radio(int data[25], int spaceDelay, int repeat) {
  int longDelay = 3*spaceDelay;
  int shortOnDelay = spaceDelay;
  int shortOffDelay = longDelay - shortOnDelay;
  int repeatDelay = 10*longDelay;

  int i, j;
  for (j = 0; j < repeat; j++) {
    for (i = 0; i < 25; i++) {
      if (data[i] == 1) {
        digitalWrite(transPin, HIGH);
        delayMicroseconds(longDelay);
      }
      else {
        digitalWrite(transPin, HIGH);
        delayMicroseconds(shortOnDelay);
        digitalWrite(transPin, LOW);
        delayMicroseconds(shortOffDelay);
      }
      digitalWrite(transPin, LOW);
      delayMicroseconds(spaceDelay);
    }
    delayMicroseconds(repeatDelay);
  }
}

// data from serial comm
char serialBuffer[100];
int serialData[25];
int serialSpaceDelay = 452;
int serialRepeat = 10;

bool readSerial() {
  bool isNew = false;

  if (Serial.available() > 0) {
    Serial.readBytes(serialBuffer, 27);

    serialSpaceDelay = int(serialBuffer[0])*10;
    serialRepeat = int(serialBuffer[1]);
    //serialSpaceDelay = 256*int(serialBuffer[0]) + 196;

    int i0 = 2;
    int i;
    for (i = 0; i < 25; i++) {
      if (int(serialBuffer[i+i0]) != serialData[i]) {
        isNew = true;
        serialData[i] = int(serialBuffer[i+i0]);
      }
    }
  }

  return isNew;
}

// data received via receivePin 
int numReceived = 0;
int receiveOn = 0;
unsigned long previousRamp;
int receive[100];
int receiveProtocol = 0;
char receiveBytes[3];

int tDiff() {
  unsigned long dt;
  unsigned long t = micros();
  if (t < previousRamp) {
    dt = t + 4294967295 - previousRamp;
  } else {
    dt = t - previousRamp;
  }
  previousRamp = t;

  return int(dt);
}

void fillReceiveBuffer() { // fill up to 25 bits only
  if (numReceived < 25) {
    int dt;
    if (receiveOn == 0 && digitalRead(receivePin)== HIGH) {
      receiveOn = 1;
      tDiff();
    } else if (receiveOn == 1 && digitalRead(receivePin) == LOW) {
      receiveOn = 0;
      dt = tDiff();

      if (dt > 1500) {
        numReceived = 0;
      } else {
        receive[numReceived] = dt;
        numReceived = numReceived + 1;
      }
    }
  }
}

int getReceiveProtocol(int max, int min, int bit0) {
  int prot;
  if (max < 600 && max > 300 && min < 300 && min > 0 && bit0 == 1) {
    prot = 1; // new one
  } else if (max < 1500 && max > 1000 && min > 0 && min < 700 && bit0 == 0) {
    prot = 2; // Chacon
  } else {
    prot = 0;
  }

  return prot;
}

char bits2byte(int x[8]) {
  int i;
  char c = char(0);
  for (i = 0; i < 8; i++) {
    c = c | (char(x[i]) << (8-i-1));
  }

  return c;
}

bool writeSerial() {
  bool isNew = false;

  if (numReceived == 25) {
    int i;
    int max = 0;
    int min = 32767;
    for (i = 0; i < 25; i++) {
      if (receive[i] > max) {
        max = receive[i];
      }
      if (receive[i] < min) {
        min = receive[i];
      }
    }
    //
    int avg = (max + min)/2;
    for (i = 0; i < 25; i++) {
      if (receive[i] < avg) {
        receive[i] = 0;
      } else {
        receive[i] = 1;
      }
    }

    receiveProtocol = getReceiveProtocol(max, min, receive[0]);
    if (receiveProtocol != 0) {
      receiveBytes[0] = bits2byte(&receive[1]);
      receiveBytes[1] = bits2byte(&receive[1+8]);
      receiveBytes[2] = bits2byte(&receive[1+16]);

      isNew = true;
      numReceived = 0;
      receiveOn = 0;
    }

    // reset receiveBuffer
    numReceived = 0;
    receiveOn = 0;
  }

  return isNew;
}

void radio2serial() {
  uint8_t message[4];
  message[0] = uint8_t(receiveProtocol);
  if (receiveProtocol == 0) {
    message[0] = 66;
  } else if (receiveProtocol < 0) {
    receiveProtocol = -receiveProtocol + 3;
  }
  message[1] = uint8_t(receiveBytes[0]);
  message[2] = uint8_t(receiveBytes[1]);
  message[3] = uint8_t(receiveBytes[2]);

  /*message[0] = receiveProtocol;
  message[1] = k;
  message[2] = 254;
  message[3] = 255;*/

  if (receiveProtocol != 0) {
    Serial.write(message, 4);
    //Serial.flush();
  }
}

void loop() {
  // read serial
  /*if (readSerial()) {
    serial2radio(serialData, serialSpaceDelay, serialRepeat);
  }*/

  fillReceiveBuffer(); 

  if (writeSerial()) {
    radio2serial();
  }
}

int main(void) {
  init();
  setup();
  for (;;) {
    loop();
  }
}
