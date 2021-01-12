#!/bin/bash
set -e

#Run app
if [ -z "$@" ]; then
  exec go run main.go
else
  exec PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin $@
fi
