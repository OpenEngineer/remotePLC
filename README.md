# remotePLC
Soft plc, configurable through text files. Inspired by Matlab Simulink and its deficiencies.

# documentation
see doc/remotePLC.pdf. I will move the introductory stuff to this readme.

# compile and install
In package root directory:
> make.sh
> static build: make.sh -s
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
* parser for Block and Line constructor arguments
* automatic documentation
* compilation for MS Windows
