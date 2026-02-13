#!/bin/sh
set -eu

cd /opt/projeto-m
sudo cp Caddyfile /etc/caddy/Caddyfile
sudo systemctl reload caddy
