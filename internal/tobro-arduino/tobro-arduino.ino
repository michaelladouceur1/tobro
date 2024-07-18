#include <ArduinoJson.h>

String inputString = "";         // a string to hold incoming data
bool stringComplete = false;

float previousMillis = 0;
int delayTime = 1000;

StaticJsonDocument<200> doc;

const char* COMMAND_KEY = "c";
const char* PIN_KEY = "p";
const char* MODE_KEY = "m";
const char* VALUE_KEY = "v";

const int COMMAND_SETUP = 1;
const int COMMAND_DIGITAL_WRITE = 2;
const int COMMAND_ANALOG_WRITE = 3;

void setup() {
  Serial.begin(115200);
  while (!Serial) continue;
  
  pinMode(LED_BUILTIN, OUTPUT);
}

void loop() {
  attemptReadSerial();
  attemptSetupPin();
  attemptDigitalWritePin();
  attemptAnalogWritePin();
}

DeserializationError attemptReadSerial() {
  if (Serial.available() > 0) {
    while (Serial.available()) {
      char inChar = (char)Serial.read();
      if (inChar == '\n') {
        stringComplete = true;
        break;
      }
      inputString += inChar;
    }

    if (!stringComplete) {
      return DeserializationError::IncompleteInput;
    }

    char json[inputString.length() + 1];
    inputString.toCharArray(json, inputString.length() + 1);

    DeserializationError error = deserializeJson(doc, json);

    if (error) {
      Serial.println("deserializeJson() failed: ");
      Serial.println(error.f_str());
      return error;
    }

    // const char* command = doc[COMMAND_KEY];

    // if (strcmp(command, "delay") == 0) {
    //   int newDelay = doc["delay"];
    //   if (newDelay > 0) {
    //     delayTime = newDelay;
    //   }
    // }

    // stringComplete = false;
    // inputString = "";
  }
}

void resetSerial() {
  stringComplete = false;
  inputString = "";
  doc.clear();
}

void attemptSetupPin() {
  if (!checkCommandReady()) {
    return;
  }

  const int command = doc[COMMAND_KEY];
  if (command == COMMAND_SETUP) {
    int pin = doc[PIN_KEY];
    int mode = doc[MODE_KEY];

    if (pin < 0 || pin > 13) {
      serialWrite("Invalid pin number");
      return;
    }

    if (mode != 0 && mode != 1) {
      serialWrite("Invalid mode");
      return;
    }

    pinMode(pin, mode);

    resetSerial();
  }
}

void attemptDigitalWritePin() {  
  if (!checkCommandReady()) {
    return;
  }

  const int command = doc[COMMAND_KEY];
  if (command == COMMAND_DIGITAL_WRITE) {
    int pin = doc[PIN_KEY];
    int value = doc[VALUE_KEY];

    if (pin < 0 || pin > 13) {
      serialWrite("Invalid pin number");
      return;
    }

    if (value != 0 && value != 1) {
      serialWrite("Invalid value");
      return;
    }

    digitalWrite(pin, value);

    resetSerial();
  }
}

void attemptAnalogWritePin() {
  if (!checkCommandReady()) {
    return;
  }

  const int command = doc[COMMAND_KEY];
  if (command == COMMAND_ANALOG_WRITE) {
    int pin = doc[PIN_KEY];
    int value = doc[VALUE_KEY];

    if (pin < 0 || pin > 13) {
      serialWrite("Invalid pin number");
      return;
    }

    if (value < 0 || value > 255) {
      serialWrite("Invalid value");
      return;
    }

    analogWrite(pin, value);

    resetSerial();
  }
}

bool checkCommandReady() {
  return stringComplete && doc.containsKey(COMMAND_KEY);
}

void serialWrite(const char* message) {
  Serial.write(message);
  Serial.write("\n");
}