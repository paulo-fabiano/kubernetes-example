#!/bin/bash
echo Deleting cluster ckad
k3d cluster deletee ckad > exec.log 2> /dev/null

echo Creating new cluster ckad
k3d cluster create ckad --servers 2 --agents 3 --port '8081:80@loadbalancer' > /dev/null