#!/bin/bash

set -euo pipefail

print_usage() {
    cat <<EOF
Usage:
  ./deploy <server_name> <server_ip> <git_ref>
  server_name: staging, ropsten, ...
EOF
}

server_name=$1
server_ip=$2
git_ref=$3

ssh ubuntu@$server_ip << EOF
    cd wallet-cache
    git pull
    git checkout $git_ref
    sudo docker-compose -f docker-compose-$server_name.yml up --build -d
EOF
