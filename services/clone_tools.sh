#!/bin/bash
SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"

if [ ! -d $SCRIPTPATH/gobuster ]; then
	git clone https://github.com/iepathos/gobuster $SCRIPTPATH/gobuster
fi
if [ ! -d $SCRIPTPATH/sublist3r ]; then
	git clone https://github.com/iepathos/sublist3r $SCRIPTPATH/sublist3r
fi
if [ ! -d $SCRIPTPATH/altdns ]; then
	git clone https://github.com/iepathos/altdns $SCRIPTPATH/altdns
fi