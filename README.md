# remotePLC: user's guide

## Synopsis
Soft plc, configurable through a text file. Inspired by Matlab Simulink.

Specify *blocks* and connect them with *lines*. The *blocks* process arrays of floats, and the *lines* pass these floats between the *blocks*.

## Usage
```
remotePLC FILE_NAME [-c CMD_STRING] [-t DELTA_T] [-s LOG_INTERVAL]
```

* FILE_NAME: name of configuration file, described below
* CMD_NAME: string of commands, in same format as file
* DELTA_T: cycle time
* LOG_INTERVAL: save a record to the log every LOG_INTERVAL cycles

## Example 1

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
