#!/bin/bash

go build -o worker main.go

docker build -t notify-worker:latest .