# remotePLC
soft plc, configurable through text files

# documentation
see doc/remotePLC.pdf

# compile and install
In package root directory:
> make.sh
> static build: make.sh -s
copies to ~/bin/ if this directory exists

# internet of things:
* Philips Hue Bridge supported, user needs to specify an IP address and a user string (see Philips Hue API reference). I included a script in ./tutorials/philipsHue/ that can return these
* Arduino serial (tutorial with 433MHz example eg. for cheap remote switches)

# license
MIT, see LICENSE.txt
