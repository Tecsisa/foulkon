#!/bin/sh

usage() { echo "Usage: $0 [-w] [-p]" 1>&2; exit 1; }
while getopts "wp" arg; do
case "${arg}" in
	w)
/go/bin/worker -config-file=/config_env_vars.toml
    ;;
	p)
/go/bin/proxy -proxy-file=/proxy_env_vars.toml
    ;;
	*)
	usage
	;;
esac
done