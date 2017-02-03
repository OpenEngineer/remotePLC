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
The configuration file in this example takes three numbers from an HTTP GET request on port 8080, then these are written to the 0th position of a file:
```
x HttpInput 8080 3
l Line x y
y FileOutput output.dat
```
Note how a block and line specification has the following format:
```
NAME BLOCK_TYPE BLOCK_ARGS
```
*Blocks* must always be named, but naming a *line* is optional, so the following configuration is equivalent:
```
x HttpInput 8080 3
Line x y
y FileOutput output.dat
```
If the *line* simply connects the input of one block to the output of another, you can substitute it by a *pipe* character. So the following configuration is also equivalent:
```
x HttpInput 8080 3 | y FileOutput output.dat
```
Because long lines can be difficult to read, you can split it using a *newline* character:
```
x HttpInput 8080 3 | \
y FileOutput output.dat
```
Finally you can add comments and empty lines for readability:
```
# take three input numbers and send them to output.dat

x HttpInput 8080 3 | \ # this listens for http requests

y FileOutput output.dat # this puts the numbers, separated by spaces, into output.dat, 
# starting at the 0th position
```

## Example 2
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
