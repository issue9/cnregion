#!/bin/sh

go build -v ./
unlink ../../data/regions.db
./fetch build -output=../../data/regions.db -data=./data