package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"syscall"
)

func main() {
	path := "./"
	if len(os.Args) >= 2 {
		path = os.Args[1]
	}

	fd, err := syscall.Open(path, 0|syscall.O_CLOEXEC, 0)
	logrus.Infof("syscall.Open fd:%d, err:%v", fd, err)

	var stat syscall.Stat_t
	err = syscall.Stat(path, &stat)
	logrus.Infof("syscall.Stat err:%v, stat:%+v", err, stat)

	file, err := os.Open("")
	file.Stat()
	err = syscall.Close(fd)
	logrus.Infof("syscall.Close err:%v", err)
}
