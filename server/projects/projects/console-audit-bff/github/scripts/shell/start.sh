#!/bin/bash

set -x

if [ "$EUID" -ne 0 ]; then
    echo "changing binary ownership to pismo user"
    chown pismo:pismo /application
    su pismo
fi

/application
