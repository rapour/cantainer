#!/bin/bash

set -eo pipefail

NAMESPACE=$1
ADDRESS=$2

if [[ -z ${NAMESPACE} ]]
then 
    echo namespace is not set 1>&2
    exit 1
fi

if [[ -z ${ADDRESS} ]]
then 
    echo ip address is not set 1>&2
    exit 1
fi

if sudo sudo ip netns exec ${NAMESPACE} ip addr show | grep -qw ${ADDRESS}; then 
    echo "the provided ip address already bound to interface ${NAMESPACE}-in"
else 
    sudo ip netns exec ${NAMESPACE} ip addr add ${ADDRESS} dev ${NAMESPACE}-in
    sudo ip netns exec ${NAMESPACE} ip link set ${NAMESPACE}-in up
    echo "ip address ${ADDRESS} binded to ${NAMESPACE}-in interface on namespace ${NAMESPACE}"
fi
