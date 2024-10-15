# Cantainer 

Cantainer is an educational project that tries to mimic some of the functionalities of container runtimes and orchestrators. It consists of two main components: a super simple container runtime that only runs alpine image on a linux host, and also a distributed orchestrator that currently can create an overlay network on top of a set of linux host machines which are connected through an arbitrary underlay network. 

The dynamic overlay network is implemented using vxlan and is managed through a cluster of daemons which are connected via [dqlite](https://dqlite.io/)'s raft implementation. 


## Host requirements (Ubuntu)

Install the following tools before launching cantainer on a host:
```bash
sudo add-apt-repository -y ppa:dqlite/dev
sudo apt update
sudo apt install libdqlite-dev gcc
```

## Local usage (Mac)

You can use Canonical's `multipass` to launch a set of Ubuntu virtual machines (we call them nodes) and then experiment with cantainer: 

1. Install multipass using homebrew:
    ```
    brew install multipass 
    ```
2. Bootstrap three instances of Ubuntu nodes on your host machines
    ```
    multipass launch --name host-1
    multipass launch --name host-2
    multipass launch --name host-3
    ```
    These virtual machines can access each other through the host's underlay network. You can access each instance using the following command:
    ```
    multipass shell $INSTALCE_NAME   
    ```
3. clone Cantainer into each of the nodes or alternatively clone it into the host and then mount Cantainer's directory on each node:
    ```
    multipass mount ./cantainer $INSTALCE_NAME:/tmp/dev
    ```
4. Go to the directory of the project inside each of the node and run the daemon:
    ```
    go run ./cmd/main.go daemon --seeds "192.168.105.17:9000"
    ```
    Since nodes are completely distributed and has no prior knowledge of other nodes in the cluster, the `seeds` option provide the IP address of at least one node so that each node can connect to other nodes in the cluster. It suffices to provide only on IP address, which belongs to one the nodes you have bootstrapped. You can find each node's IP address using the following command:
    ```
    multipass list
    ```
    Note: Currently, the daemon process is attached to the shell so that debugging become easier. 
5. Now you can create and run a container on each of the nodes:
    ```
    sudo go run ./cmd/main.go new -n "10.0.0.0/24"
    ```
    You can specify a subnet for the container to attach to it. Every container that is on the same subnet is accessible to other containers on that network regardless of the node on which they reside. 