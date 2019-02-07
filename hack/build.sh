#!/bin/bash
set -e
OUTFILE=/usr/local/bin/go-sdk
go build -o $OUTFILE main.go

chmod +x $OUTFILE