#!/bin/sh

/usr/bin/system-docker run --rm --net=host etcdctl mk /coreos.com/network/config '{"Network":"10.244.0.0/16", "Backend": {"Type": "vxlan"}}'
/flannel --iface=eth1
