#include <ArduinoJson.h>

String inputString = "";         // a string to hold incoming data
bool stringComplete = false;

float previousMillis = 0;
int delayTime = 1000;

StaticJsonDocument<200> doc;


void setup() {
  Serial.begin(9600);
  while (!Serial) continue;
  
  pinMode(LED_BUILTIN, OUTPUT);
}

void loop() {

  attemptReadSerial();
  checkSetupPin();
  checkDigitalWritePin();
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

    // const char* command = doc["command"];

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

void checkSetupPin() {
  if (!stringComplete) {
    return;
  }

  if (!doc.containsKey("command")) {
    return;
  }

  const char* command = doc["command"];
  if (strcmp(command, "setup_pin") == 0) {
    int pin = doc["pin"];
    int mode = doc["mode"];

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

void checkDigitalWritePin() {  
  if (!stringComplete) {
    return;
  }

  if (!doc.containsKey("command")) {
    return;
  }

  const char* command = doc["command"];
  if (strcmp(command, "digital_write_pin") == 0) {
    int pin = doc["pin"];
    int value = doc["value"];

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

// void attemptUpdateDelay() {
//   if (inputString.length() > 0) {
//     int newDelay = inputString.toInt();
//     if (newDelay > 0) {
//       delayTime = newDelay;
//     }
//   }
// }

// void checkForToggleLight() {
//   if (previousMillis + delayTime < millis()) {
//     previousMillis = millis();
//     digitalWrite(LED_BUILTIN, !digitalRead(LED_BUILTIN));
//   }
// }

void serialWrite(const char* message) {
  Serial.write(message);
  Serial.write("\n");
}