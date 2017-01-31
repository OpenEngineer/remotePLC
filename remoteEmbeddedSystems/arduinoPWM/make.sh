#!/bin/bash
cd ${0%/*}

make
sudo make upload
