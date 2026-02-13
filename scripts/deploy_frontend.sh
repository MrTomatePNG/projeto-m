#!/bin/sh
set -eu

export PATH="$PATH:/usr/local/go/bin:$HOME/.bun/bin"

cd /opt/projeto-m/client
bun install
bun run build
