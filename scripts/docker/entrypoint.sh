#!/bin/sh

usage() { echo "Usage: worker|proxy" 1>&2; exit 1; }
if [ "$1" = 'worker' ]; then
    worker -config-file=/worker.toml
elif [ "$1" = 'proxy' ]; then
    proxy -proxy-file=/proxy.toml
else
	usage
fi