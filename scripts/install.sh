#!/bin/bash

SUDO=
if [ "$(id -u)" -ne 0 ]; then
    if ! available sudo; then
        error "This script requires superuser permissions. Please re-run as root."
    fi
    SUDO="sudo"
fi

tmp_dir=$(mktemp -d)

echo "Downloading kanmit..."
git clone https://github.com/1704mori/kanmit.git "$tmp_dir"

cd "$tmp_dir"

VERSION=$(git rev-parse --short HEAD)

echo "Building kanmit..."
go build -ldflags "-X main.version=$VERSION" main.go

$SUDO mv main /usr/local/bin/kanmit

rm -rf "$tmp_dir"

echo "Install complete"
