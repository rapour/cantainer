#!/bin/bash

set -eo pipefail

VETH_NAME=$1
NAMESPACE=$2

if [[ -z ${VETH_NAME} ]]
then 
    echo veth pair name is not set 1>&2
    exit 1
fi

if [[ -z ${NAMESPACE} ]]
then 
    echo namespace is not set 1>&2
    exit 1
fi

if sudo ip netns exec ${NAMESPACE} ip link show | grep -qw "${VETH_NAME}-in"; then 
    echo "veth ${VETH_NAME} is already plugged into namespace ${NAMESPACE}"
else 
    sudo ip link set "${VETH_NAME}-in" netns ${NAMESPACE}
    sudo ip netns exec ${NAMESPACE} ip link set ${VETH_NAME}-in up
    echo "veth ${VETH_NAME} plugged into namespace ${NAMESPACE}"
fi
