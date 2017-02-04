# remotePLC

This is a *soft* PLC program. It is configured through a text file.

You specify *blocks* and connect them with *lines*. The *blocks* process arrays of 64 bit floating point numbers, and the *lines* pass these arrays between the *blocks*.

There are five types of *blocks*:

* input blocks (remote sensors, switches...)
* stateless node blocks (math operators, bool operators...)
* logic blocks, these have a state (PID controls, time shift delays...)
* output blocks (actuators, lights...)
* stop blocks, stop the program on certain conditions (timeout, stability...)

By configuring and connecting *blocks* in the right way you can automate your outputs based on your user and environment inputs.

## Usage
From your shell:
```
remotePLC FILE_NAME [-c CMD_STRING] [-t DELTA_T] [-s LOG_INTERVAL]
```

* `FILE_NAME`     name of configuration file, described below
* `CMD_NAME`      string of commands, in same format as file
* `DELTA_T`       PLC cycle time
* `LOG_INTERVAL`  save a record to the log every this many cycles

## Installation
Download this repository:
```
git clone https://github.com/christianschmitz/remotePLC remotePLC
```
Make sure you have the golang dev tools install, for example with aptitude:
```
sudo apt-get install golang
```
Go into the repository root directory, and run the compilation script:
```
cd ./remotePLC
./make.sh
```
For statically linked compilation:
```
./make.sh -s
```
This script also copies the *remotePLC* binary to `$HOME/bin` if this directory exists. Proper installation to some `/usr/bin` directory is not yet included.

# Examples 

## Example 1
The configuration file in this example takes three numbers from an HTTP GET request on port 8080, then these are written to the 0th position of `output.dat`:
```
x HttpInput 8080 3
l Line x y
y FileOutput output.dat
```
Note how *block* and *line* specifications have the following format:
```
NAME BLOCK_TYPE BLOCK_ARGS
```
*Blocks* must always be named, but naming a *line* is optional, so the following configuration is equivalent:
```
x HttpInput 8080 3
Line x y
y FileOutput output.dat
```
If the *line* simply connects the input of one *block* to the output of another, you can substitute it by a *pipe* character. So the following configuration is also equivalent:
```
x HttpInput 8080 3 | y FileOutput output.dat
```
Because long text lines can be difficult to read, you can split them using a *newline* character:
```
x HttpInput 8080 3 | \
y FileOutput output.dat
```
You can add comments and empty lines for readability:
```
# take three input numbers and send them to output.dat

x HttpInput 8080 3 | \ # this block listens for http requests of the form: http://127.0.0.1:8080/1.0,2.0,3.0

# (this is bad usage of empty lines, but works)

# this block puts the numbers, separated by spaces, into output.dat, 
# starting at the 0th position:
y FileOutput output.dat 
```
The output arrays of *blocks* or *lines*, with names ending with an underscore (eg. `x_`), are not saved in the data log:
```
x_ HttpInput 8080 3 | y_ FileOutput output.dat # no data is logged
```

## Example 2
In this example three numbers from an http request are sent as a brightness, hue, and saturation value to three Philips Hue lights. In order not to interfere with the smart phone app, the three values are only sent after a new valid http request has been received:
```
x HttpInput 8080 3 | n Node

ForkLine n light1 light2 light3

light1 PhilipsHueOutput 192.168.1.6 T08t2C8GF9KEqXYRI8PBzb3M6vDjteT3hxdERW8z 1
light2 PhilipsHueOutput 192.168.1.6 T08t2C8GF9KEqXYRI8PBzb3M6vDjteT3hxdERW8z 2
light3 PhilipsHueOutput 192.168.1.6 T08t2C8GF9KEqXYRI8PBzb3M6vDjteT3hxdERW8z 3
```

Every time a client sends a valid http request, `HttpInput` propagates those numbers. In an idle state `HttpInput` propagates `UNDEFINED` numbers. The `PhilipsHueOutput` detects `UNDEFINED` numbers and does nothing. This way the lights are only switched when a new http request is received, thus upon a client action.
## Example 3
This example combines an `HttpInput` with a 433MHz receiver. Both inputs are used to switch Philips Hue lights, and 433MHz lights.

The 433MHz lights protocol uses OOK. An Arduino can be used along with a cheap transmitter and receiver to interface with these external 433MHz receivers and transmitters. *remotePLC* includes a protocol that, via a serial port, writes or reads a PWM signal to an Arduino pin, in turn connected to the cheap transmitter or receiver.
```
# handle the http input
in1 HttpInput 8080 3
SplitLine 1 in1 switch1 hue sat
switch1_ Node

# save the hue and sat state
hue_ DefineLogic 0.1 #default values if upstream is UNDEFINED
sat_ DefineLogic 0.1

# handle the 433MHz input
#  the arguments are: 
#   PORT 
#   BITRATE 
#   NUMBYTES 
#   PULSEWIDTH 
#   CLEARCOUNT 
#   TIMEOUTCOUNT 
#   PULSEMARGIN
in2_ ArduinoPWMInput /dev/ttyACM0 9600 40 210 30 20000 30 | \
switch2_ MapNode map_in.dat 40 1 Mode exact

# combine the inputs
JoinLine switch1 switch2 zero_or_one
# numbers smaller or equl to 0.5 are set 0, greater than 0.5 are set to 1
# UNDEFINED numbers are left unchanged
zero_or_one_ IfElseElseNode 0 0.5 1 | \
switch_ ReductionNode Or # 0, UNDEFINED, or 1

# write to the Philips Hue lights
JoinLine switch_ hue sat ph_state
ph_state Node
ForkLine ph_state light1 light2 light3
light1 PhilipsHueOutput 192.168.1.6 T08t2C8GF9KEqXYRI8PBzb3M6vDjteT3hxdERW8z 1
light2 PhilipsHueOutput 192.168.1.6 T08t2C8GF9KEqXYRI8PBzb3M6vDjteT3hxdERW8z 2
light3 PhilipsHueOutput 192.168.1.6 T08t2C8GF9KEqXYRI8PBzb3M6vDjteT3hxdERW8z 3

# write to the 433Mhz lights
Line switch rf_state
rf_state_ MapNode map_out.dat 1 40 Mode exact | relay1_ ArduinoPWMOutput /dev/ttyACM0 9600 210 30 10
```

`map_in.dat` and `map_out.dat` can contain comments and can look like:
```
# off code
128   5   2 129  64 129  80  40  20   8 20  10 129   2 160  64 160  84  10   5   2 129  64 160  80  40  20  10   5   2 129  64 160  80  32  84  10   4   0   0               0
# on code
128   5   2 129  64 129  80  40  20   8 20  10 129   2 160  64 160  84  10   5   2 129  64 160  80  40  20  10   5   2 129  64 129  80  32  84  10   4   0   0               1
```

Long configurations like this are a form of meta-programming, and bugs can quickly be introduced. That is why a data log is written every few cycles. This log contains a column for every number (if not hidden). For arrays containing more than one number, an index subscript is added to the block name in the log header.

Once a configuration has been debugged an underscore can be appended to the names of the *blocks*, this limits the amount of data being written. The underscores don't need to be added to every occurance of that blockname in the configuration though. The *lines* see regular names and underscores as identifiers of the same *block*, the underscore only acts as a hiding marker.


# Internet of things:
*remotePLC* is intended as a utility for easy home automation. Some of the devices that are supported:

* Philips Hue lights. The user needs to specify an IP address and a user string (see Philips Hue API reference). A script is included in `./tutorials/philipsHue/` that can return these.
* Arduino PWM in and out (eg. cheap 433MHz remote switches)

## Remote Embedded Systems
The `./remoteEmbeddedSystems/` folder contains source code intended for:

* Arduino: duplexPWM code (requires the avr and arduino dev tools)

# License
GPLv3

# Acknowledgements

* Martin Ling and others: *libserialport*
* mikepb: golang interface for *libserialport*

# Developer's guide

All blocks can expect upstream arrays of constant length, and must generate downstream arrays of constant length. This is because some blocks depend on positional numbers, and any change in length anywhere in the configuration can generate errors in those blocks.

## Input blocks
Input blocks should return actions, and not states: eg. "I flipped a switch", and not "the switch is on". This is because it is easier to turn actions into states than vice versa.

There are two types of input blocks:

1. Polling type
2. Serving type

A reply to a poll can be considered an action, eg. a reading of a sensor.
A request to a server can also be considered an action, eg. change light state via http request.

Switches will be idle most of the time. For this *actionless* state, the downstream arrays should be set to `UNDEFINED`.

For polling inputs with high latency, or any serving inputs, the polling/serving function should run in the background. This function should be launched upon construction of the input block and set to loop infinitely by itself. The `Update()` function should then poll the output of the polling/serving function via an intermediate array.

## Todo
* automatic documentation
* support for MS Windows
