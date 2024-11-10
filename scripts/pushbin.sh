#!/bin/bash

go build -o bin/fzfoxide cmd/fzfoxide/main.go

sudo mv bin/fzfoxide /usr/local/bin