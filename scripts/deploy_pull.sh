#!/usr/bin/env bash
set -euo pipefail

cd /opt/projeto-m
git config --global --add safe.directory /opt/projeto-m
git pull origin main
