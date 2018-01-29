## shapeshift-notifier polls the shapeshift.io api to check exchange rates.

## Tested only on Mac
I've tested/used this only on a Mac.
The notification part, especially, might not work on other platforms.  Disable it with -popup=false if required.

## Download and Usage
You can download the executable version for your platform (windows/linux/mac) from the releases tab and run it as shown in the Usage example below.

## Requirements
* [dep](https://github.com/golang/dep)
* golang

## Get and Build
	go get github.com/sathishvj/shapeshift-notifier
	dep ensure
	go build

## Usage
	Usage example: shapeshift-notifier -popup=false -interval=32 "snt_bat,>0.75,=100000" "eth_btc<0.01" "rlc_gnt,=150"
	Defaults: popup=true, interval=30, args="eth_btc,>0.1,=0"
	Signs: Only > and < are allowed for operations. = indicates the amount to convert.  Only the first part with token codes is mandatory.

