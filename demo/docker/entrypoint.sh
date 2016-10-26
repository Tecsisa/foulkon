#!/bin/sh

usage() { echo "Usage: proxy|api|web" 1>&2; exit 1; }
if [ "$1" = 'proxy' ]; then
    proxy -proxy-file=/proxy.toml
elif [ "$1" = 'api' ]; then
    api
elif [ "$1" = 'web' ]; then
    web
else
	usage
fi