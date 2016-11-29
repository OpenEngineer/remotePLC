# remotePLC
soft plc, configurable through text files

# details
I wrote this out of frustation with unconfigurable controller software that requires rebuilding whenever there is a change. For linux based industrial controllers this is really unecessary. Also the fact that PLC logic is now a utility, means that it can easily be driven by higher level software (like dakota) for doing DoEs.

# building
In package root directory:
> go build
Static build:
> make.sh

# internet of things:
* Philips Hue Bridge supported, user needs to specify an IP address and a user string (see Philips Hue API reference). I included a script in ./tutorials/philipsHue/ that can return these
