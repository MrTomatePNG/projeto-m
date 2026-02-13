#!/bin/sh
set -eu

export PATH="$PATH:/usr/local/go/bin:$HOME/.bun/bin"

cd /opt/projeto-m
go build -o backend ./cmd
