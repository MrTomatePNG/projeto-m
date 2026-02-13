#!/usr/bin/env bash
set -euo pipefail

cd /opt/projeto-m
sudo cp Caddyfile /etc/caddy/Caddyfile
sudo systemctl reload caddy
