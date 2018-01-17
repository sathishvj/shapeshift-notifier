## shapeshift-notifier polls the shapeshift.io api to check exchange rates.

## Tested only on Mac
I've tested/used this only on a Mac.
The notification part, especially, might not work on other platforms.  Disable it with -popup=false if required.

## Requirements
* [dep](https://github.com/golang/dep)
* golang

## Get and Build
	go get github.com/sathishvj/shapeshift-notifier
	dep ensure
	go build

## Usage
	Usage example: shapeshift-notifier -popup=false -interval=32 "snt_bat>0.75" "eth_btc<0.01
	Defaults: popup=true, interval=30, args="eth_btc>0.1"
	Signs: Only > and < are allowed for operations.

