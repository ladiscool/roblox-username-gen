#!/bin/bash

go build src/main.go

if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || "$OSTYPE" == "win32" ]]; then
  executable="./main.exe"
else
  executable="./main"
  chmod +x "$executable"
fi
"$executable"
