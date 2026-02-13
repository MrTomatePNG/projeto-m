#!/usr/bin/env bash
set -euo pipefail

export PATH="$PATH:/usr/local/go/bin:$HOME/.bun/bin"

cd /opt/projeto-m
go build -o backend ./cmd/api
