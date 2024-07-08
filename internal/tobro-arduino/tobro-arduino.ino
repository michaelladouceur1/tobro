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

  checkForToggleLight();
  attemptReadSerial();
  // attemptUpdateDelay();
  // delay(10);
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
    serialWrite(json);

    DeserializationError error = deserializeJson(doc, json);

    if (error) {
      Serial.println("deserializeJson() failed: ");
      Serial.println(error.f_str());
      return error;
    }

    const char* command = doc["command"];

    if (strcmp(command, "delay") == 0) {
      int newDelay = doc["delay"];
      if (newDelay > 0) {
        delayTime = newDelay;
      }
    }

    stringComplete = false;
    inputString = "";
  }
}

void attemptUpdateDelay() {
  if (inputString.length() > 0) {
    int newDelay = inputString.toInt();
    if (newDelay > 0) {
      delayTime = newDelay;
    }
  }
}

void checkForToggleLight() {
  if (previousMillis + delayTime < millis()) {
    previousMillis = millis();
    digitalWrite(LED_BUILTIN, !digitalRead(LED_BUILTIN));
  }
}

void serialWrite(const char* message) {
  Serial.write(message);
  Serial.write("\n");
}