# remotePLC: user's guide

## Synopsis
This is a *soft* PLC program. It is configured through a text file.

You specify *blocks* and connect them with *lines*. The *blocks* process arrays of 64 bit floating point numbers, and the *lines* pass these arrays between the *blocks*.

There are four types of *blocks*:

* input blocks (remote sensors, switches...)
* stateless node blocks (math operators, bool operators...)
* logic blocks (PID controls, time shift delays...)
* output blocks (actuators, lights...)

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
x HttpInput 8080 3

UndefineLine x n
n Node
ForkLine n light1 light2 light3

light1 PhilipsHueOutput 192.168.1.6 T08t2C8GF9KEqXYRI8PBzb3M6vDjteT3hxdERW8z 1
light2 PhilipsHueOutput 192.168.1.6 T08t2C8GF9KEqXYRI8PBzb3M6vDjteT3hxdERW8z 2
light3 PhilipsHueOutput 192.168.1.6 T08t2C8GF9KEqXYRI8PBzb3M6vDjteT3hxdERW8z 3
```

The `UndefineLine` takes the numbers from `HttpInput` and sends them to `Node n`. Then the output of `HttpInput` is set to `UNDEFINED`. The `PhilipsHueOutput` detects `UNDEFINED` numbers and does nothing.
This scheme assures that the lights are only switched if a new http request is received.

## Example 3

## documentation
see doc/remotePLC.pdf. I will move the introductory stuff to this readme.

## compile and install
In package root directory:
```
> make.sh
> static build: make.sh -s
```
copies to ~/bin/ if this directory exists

# internet of things:
* Philips Hue Bridge supported, user needs to specify an IP address and a user string (see Philips Hue API reference). I included a script in ./tutorials/philipsHue/ that can return these
* Arduino serial (tutorial with 433MHz example eg. for cheap remote switches)

## Remote Embedded Systems
the ./remoteEmbeddedSystems/ folder contains source code intended for eg:
* the arduino duplexPWM code

# license
GPL3

# TODO:
* automatic documentation
* compilation for MS Windows
