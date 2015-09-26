# outPorts

## INSTALL

	go get github.com/nhooyr/outPorts

## USAGE

	outPorts min[-max]

### EXAMPLES
	outPorts 20-30

check ports 20 to 30

	outPorts 20-10

check ports 20 to (20+10)

	outPorts 25

check port 25

	outports -h

to see documentation on terminal

##NOTE

it works asynchronously so the output will not always be in order, but it is MUCH MUCH faster.

##SPECIAL THANKS TO CREATOR OF PORTQUIZ.NET!
