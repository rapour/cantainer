#!/bin/bash

set -eo pipefail

VXLAN_NAME=$1
REMOTE_ADDRESS=$2

if [[ -z ${VXLAN_NAME} ]]
then 
    echo vxlan name is not set 1>&2
    exit 1
fi

if [[ -z ${REMOTE_ADDRESS} ]]
then 
    echo remote address of peer is not set 1>&2
    exit 1
fi

if ! $(bridge fdb show dev ${VXLAN_NAME} | grep -qw ${REMOTE_ADDRESS}); then 
    exit 0
else 
    sudo bridge fdb del 00:00:00:00:00:00 dst ${REMOTE_ADDRESS} dev ${VXLAN_NAME}
fi