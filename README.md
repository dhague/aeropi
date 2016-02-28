# aeropi
Air density monitor &amp; display for Raspberry Pi using Adafruit LCD display and BME280 temperature/humidity/pressure sensor.

From Windows prompt:
set GOARCH=arm
set GOARM=6
set GOOS=linux
go build -v

Side note: make sure to run aeropi with *sudo*

Current status: Reports sensor data & air density to stdout
