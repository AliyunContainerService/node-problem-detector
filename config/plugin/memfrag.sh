#!/bin/bash

# check max fd open files

OK=0
NONOK=1
UNKNOWN=2

filepath="/proc"
if [ -f "/host/proc" ]; then
filepath="/host/proc"
fi

nr=0

memTotal=`cat $filepath/meminfo | grep "MemTotal" | awk '{print $2}'`

if [[ $memTotal -lt 4194304 ]]; then
	echo "node has no fragment"
	exit $OK
fi

frag=`cat $filepath/buddyinfo | grep " Normal " |awk '{if ($9 + $10 + $11 + $12 + $13 + $14 + $15 < 100 ) print $0}' | wc -l`

if [[ $frag -gt 0 ]]; then
   echo "Memory buddy system fragment"
   exit $NONOK
fi
echo "node has no fragment"
exit $OK
