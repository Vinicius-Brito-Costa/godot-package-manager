#! /bin/bash

path_to_go="$(which go)"

(export GOOS="windows"; export GOARCH="amd64"; exec "$path_to_go" build -o gpm-windows-amd64.exe)
(export GOOS="windows"; export GOARCH="386"; exec "$path_to_go" build -o gpm-windows-386.exe)
(export GOOS="linux"; export GOARCH="arm"; exec "$path_to_go" build -o gpm-linux-arm)
(export GOOS="darwin"; export GOARCH="arm64"; exec "$path_to_go" build -o gpm-darwin-arm64.dmg)
