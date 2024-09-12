#!/bin/bash

set -eo pipefail

NAME=$1
BRIDGE_NAME=$2
VID=$3

if [[ -z ${NAME} ]]
then 
    echo vxlan name is not set 1>&2
    exit 1
fi

if [[ -z ${BRIDGE_NAME} ]]
then 
    echo bridge name to connect to vxlan is not set 1>&2
    exit 1
fi

if [[ -z ${VID} ]]
then 
    echo vxlan id is not set 1>&2
    exit 1
fi

re='^[0-9]+$'
if ! [[ $VID =~ $re ]] ; then
   echo "vxlan id ${VID} is not a number" >&2; 
   exit 1
fi

if sudo ip link list type vxlan | grep -qw ${NAME}; then 
    echo "vxlan ${NAME} on host is already set up"
else 
    DEFAULT_INTERFACE=$(sudo ip route | grep default | awk '{print $5}')
    DEFAULT_INTERFACE_ADDRESS=$(sudo ip route | grep default | awk '{print $9}')
    sudo ip link add ${NAME} type vxlan id ${VID} local ${DEFAULT_INTERFACE_ADDRESS} dstport 0 dev ${DEFAULT_INTERFACE}
    sudo ip link set ${NAME} master ${BRIDGE_NAME}
    sudo ip link set ${NAME} up
    echo "vxlan ${NAME} set up on host"
fi