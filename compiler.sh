#!/bin/bash

go mod init src/main.go && go mod tidy && go build src/main.go

if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || "$OSTYPE" == "win32" ]]; then
  executable="./main.exe"
else
  executable="./main"
  chmod +x "$executable"
fi
"$executable"
