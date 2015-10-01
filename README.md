# outPorts

## Description

uses [portquiz](http://portquiz.net) to check outbound ports (you should check it out!)

## INSTALL

	go get github.com/nhooyr/outPorts

## USAGE

	outPorts min[-max]

### EXAMPLES
check from port 1 to 65535
	outPorts

check from port 20 to 30 and then 40-50
	outPorts 20-30 40-50

check from port 20 to 10 and then 40 to 10
	outPorts 20-10 40-10

check port 25
	outPorts 25

check from port 20-25 and only display failure (carries onto next port(s))
	outPorts 20-25f

check from port 20-25 and only display success (carries onto next port(s))
	outPorts 20-25s

check from port 20-25 and display failure/success (needed to reset only failure/success)
	outPorts 20-25a

to see documentation on the terminal
	outPorts -h

##NOTE

it works asynchronously so the output will not always be in order, but it is very fast.
