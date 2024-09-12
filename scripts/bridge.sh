#!/bin/bash

set -eo pipefail

NAME=$1

if [[ -z ${NAME} ]]
then 
    echo bridge name is not set 1>&2
    exit 1
fi

if sudo ip link list type bridge | grep -qw ${NAME}; then 
    echo "the host bridge ${NAME} is already set"
else 
    sudo ip link add ${NAME} type bridge
    sudo ip link set ${NAME} up
    echo "bridge ${NAME} set up on host"
fi