#!/bin/bash

set -eo pipefail

NAME=$1

if [[ -z ${NAME} ]]
then 
    echo veth pair name is not set 1>&2
    exit 1
fi

if ! $(sudo ip link list type veth | grep -qw "${NAME}-out"); then 
    exit 0
else 
    sudo ip link del "${NAME}-in"
    echo "veth ${NAME}-in/${NAME}-out deleted"
fi