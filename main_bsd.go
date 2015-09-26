// +build darwin dragonfly freebsd netbsd openbsd

package main

import "syscall"

const ioctlReadTermios = syscall.TIOCGETA
const ioctlWriteTermios = syscall.TIOCSETA
