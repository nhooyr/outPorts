# outPorts

## Description

Uses [portquiz](http://portquiz.net) to check outbound ports (you should check it out!)

__note: abandoned and has bugs__

## INSTALL

	go get github.com/nhooyr/outPorts

## USAGE

	outPorts [-t n][-w n][-c] min[-max][sf]...

The flag -t allows you to change the timeout for each connection in seconds, default is 0.
The flag -w allows you to change the number of workers used (concurrent TCP requests), so if you're getting false positives, use less. Default is 1024.
The flag -c allows for color/bold output which makes it easier to see.

### EXAMPLES
check from ports 1 to 65535

    outPorts all

check from ports 20-30 and then 40-50

    outPorts 20-30 40-50

check from ports 20-10 and then 40-10

    outPorts 20-10 40-10

check port 25

    outPorts 25

check from ports 1-65535 and only display failure

    outPorts allf

check from ports 20-25 and only display success 

    outPorts 20-25s

check from ports 20-25 and only display success, then check from ports 30-35 and only display failure, then check from ports 40-50 and display both.

    outPorts 20-25s 30-35f 40-50

to see documentation on the terminal

    outPorts -h

##NOTE

it works asynchronously so the output will not always be in order, but it is very fast.
