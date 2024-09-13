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