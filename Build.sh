#!/bin/bash
#go tool dist list
#0.0.0
#clean up the old builds
rm ./builds/*

#compile experimental: need CGO for go-sqlite, can't use cgo due to portability issues
#GOOS=linux GOARCH=amd64 go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o builds/apachelogparser.tmp .
GOOS=linux GOARCH=amd64 go build -o builds/apachelogparser.tmp .
GOOS=windows GOARCH=amd64 go build -o builds/apachelogparser.exe .

#compress the linux binary, this will be the one i distribute the most (due to most webservers running on linux)
upx -f --brute -o builds/apachelogparser builds/apachelogparser.tmp

#copying config.ini in the builds directory
cp config.ini builds/config.ini

#compressing binarys for supported os's and architectures
tar cvzf builds/apachelogparser-v$1-linux-amd64.tar.gz builds/apachelogparser builds/config.ini
7zr a builds/apachelogparser-v$1-windows-amd64.zip builds/apachelogparser.exe builds/config.ini

#create the nfpm.yaml including my current version 
sed "s/vmyversion/v$1/" nfpm.yaml.template > nfpm.yaml

#packaging .deb and .rpm
nfpm pkg --packager deb --target /root/apachelogparser/builds/
nfpm pkg --packager rpm --target /root/apachelogparser/builds/

#signing
gpg --detach-sign builds/*.deb
gpg --detach-sign builds/*.rpm
gpg --detach-sign builds/apachelogparser-v$1-linux-amd64.tar.gz
gpg --detach-sign builds/apachelogparser-v$1-windows-amd64.zip

#cleaning stuff i don't need anymore
rm ./builds/apachelogparser
rm ./builds/apachelogparser.exe
rm ./builds/apachelogparser.tmp
rm ./builds/config.ini
