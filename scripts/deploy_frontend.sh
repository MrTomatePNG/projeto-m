#!/usr/bin/env bash
set -euo pipefail

export PATH="$PATH:/usr/local/go/bin:$HOME/.bun/bin"

cd /opt/projeto-m/client
bun install
bun run build
