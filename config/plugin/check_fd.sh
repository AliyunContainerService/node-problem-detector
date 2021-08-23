#!/bin/bash

# check max fd open files

OK=0
NONOK=1
UNKNOWN=2

filepath="/proc"
if [ -f "/host/proc" ]; then
filepath="/host/proc"
fi


max=$(cat $filepath/sys/fs/file-max)
file_nr=$(cat $filepath/sys/fs/file-nr | awk '{print $1}')

if [[ $file_nr -gt $((max*80/100)) ]]; then
   echo "current fd usage is $count and max is $max"
   exit $NONOK
fi
echo "node has no fd pressure"
exit $OK
