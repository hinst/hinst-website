#!/bin/bash
set -e
go run build/go_mod_no_replace/main.go > go.mod.new
mv go.mod.new go.mod