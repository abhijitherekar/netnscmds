#!/bin/bash

set -e 
set -x

container_name = $1

sudo docker ps | grep container_name

pid = sudo docker inspect -f '{{.State.Pid}}' container_id

sudo mkdir -p /var/run/netns

sudo ln -sf /proc/pid/ns/net "/var/run/netns/container_name"

sudo ip netns exec ovs ip a

#Create the veth "tap1" pair at the host and the peer test-veth
sudo ip link add tap1 type veth peer name container_name-veth

#now move the container_name-veth into container.
sudo ip link set container_name-veth netns contianer_name

