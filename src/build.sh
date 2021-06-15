#!/usr/bin/env bash
go build -ldflags "-w -s" main.go
mv -f main cmdb_agent
upx --brute cmdb_agent