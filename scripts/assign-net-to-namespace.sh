#!/bin/bash

set -eo pipefail

NAMESPACE=$1
VETH_NAME=$2
NETWORK=$3

if [[ -z ${NAMESPACE} ]]
then 
    echo namespace is not set 1>&2
    exit 1
fi

if [[ -z ${VETH_NAME} ]]
then 
    echo veth inside namespace is not set 1>&2
    exit 1
fi

if [[ -z ${NETWORK} ]]
then 
    echo prefix is not set 1>&2
    exit 1
fi

if sudo sudo ip netns exec ${NAMESPACE} ip addr show | grep -qw ${NETWORK}; then 
    echo "the provided prefix already bound to interface ${VETH_NAME}-in"
else 
    sudo ip netns exec ${NAMESPACE} ip addr add ${NETWORK} dev ${VETH_NAME}-in
    sudo ip netns exec ${NAMESPACE} ip link set ${VETH_NAME}-in up
    echo "prefix ${NETWORK} binded to ${VETH_NAME}-in interface on namespace ${NAMESPACE}"
fi
