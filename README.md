# outPorts

## Description

uses [portquiz](http://portquiz.net) to check outbound ports (you should check it out!)

## INSTALL

	go get github.com/nhooyr/outPorts

## USAGE

	outPorts min[-max]

### EXAMPLES
	outPorts

check ports 1 to 65535

	outPorts 20-30

check ports 20 to 30

	outPorts 20-10

check ports 20 to (20+10)

	outPorts 25

check port 25

	outPorts -h

to see documentation on the terminal

##NOTE

it works asynchronously so the output will not always be in order, but it is very fast.
