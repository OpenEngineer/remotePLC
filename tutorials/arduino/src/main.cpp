# include <Arduino.h>
// 433MHz example

// serial tx and rx setup
int transPin = 13;

void setup() {
  Serial.begin(9600);
  pinMode(transPin, OUTPUT);
}

void transmit(int data[25], int spaceDelay, int repeat) {
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

void loop() {
  // read serial
  if (readSerial()) {
    transmit(serialData, serialSpaceDelay, serialRepeat);
  }
}

int main(void) {
  init();
  setup();
  for (;;) {
    loop();
    delay(10);
  }
}
