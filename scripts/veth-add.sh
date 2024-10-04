#!/bin/bash

set -eo pipefail

NAME=$1

if [[ -z ${NAME} ]]
then 
    echo veth pair name is not set 1>&2
    exit 1
fi

if sudo ip link list type veth | grep -qw "${NAME}-out"; then 
    echo "veth pair already exists"
else 
    sudo ip link add "${NAME}-in" type veth peer name "${NAME}-out"
    echo "veth ${NAME}-in/${NAME}-out created"
fi