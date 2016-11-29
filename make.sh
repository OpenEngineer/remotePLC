#!/bin/bash
if [ $# -eq 1 -a "$1" == "-s" ]
then
  export CGO_ENABLED=0
  go build -ldflags "-linkmode external -extldflags -static"
else
  go build
fi

if [ -d ~/bin/ ]
then
  cp remotePLC ~/bin/
fi
