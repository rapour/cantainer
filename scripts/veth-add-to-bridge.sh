#!/bin/bash

set -eo pipefail

VETH_NAME=$1
BRIDGE_NAME=$2

if [[ -z ${VETH_NAME} ]]
then 
    echo veth pair name is not set 1>&2
    exit 1
fi

if [[ -z ${BRIDGE_NAME} ]]
then 
    echo bridge name is not set 1>&2
    exit 1
fi

if sudo ip link list type veth | grep -qw "${VETH_NAME}-out" | grep -qw "${BRIDGE_NAME}" ; then 
    echo "bridge "${BRIDGE_NAME} is already the master of ${VETH_NAME}-out"
else 
    sudo ip link set ${VETH_NAME}-out down
    sudo ip link set ${VETH_NAME}-out master ${BRIDGE_NAME}
    sudo ip link set ${VETH_NAME}-out up
fi
