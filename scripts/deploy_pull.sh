#!/bin/sh
set -eu

#atualiza o projeto
cd /opt/projeto-m
git config --global --add safe.directory /opt/projeto-m
git pull origin main
