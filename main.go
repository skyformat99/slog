package slog

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/themester/GoSlaves"
)

const (
	Snappy = iota
	Gzip
	info
	panic
	warning
	perror
	fatal
)

var (
	defaultTime      = time.Second * 30
	defaultFormatter = &logrus.TextFormatter{
		DisableColors:   false,
		DisableSorting:  false,
		FullTimestamp:   true,
		ForceColors:     true,
		TimestampFormat: "15:04:05 02/01/2006",
	}
)

type t struct {
	data int
	obj  interface{}
}

type Log struct {
	q    *slaves.Queue
	ch   chan struct{}
	stop chan struct{}
}

// NewLog creates new log structure which you can
// write log without waiting disk writting.
//
// file parameter is the output file for logging.
// checkTime is left time to check the file size. If 0 defaultTime is set.
// size is the max log file size. In KB. If 0 default size will be 5120.
// algo is the algorightm to compress the file. (See #Constants)/
func NewLog(file string, checkTime, size, algo int) (*Log, error) {
	log := &Log{}

	switch algo {
	case Snappy:
	case Gzip:
	default:
		return nil, fmt.Errorf("Unknow algorithm")
	}

	if size == 0 {
		size = 1024 * 5
	}
	if checkTime > 0 {
		defaultTime = time.Second * time.Duration(checkTime)
	}

	osf, err := os.OpenFile(file,
		os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600,
	)
	if err != nil {
		return nil, err
	}
	logrus.SetOutput(osf)
	logrus.SetFormatter(defaultFormatter)

	log.ch = make(chan struct{}, 1)
	log.stop = make(chan struct{}, 1)
	log.q = slaves.DoQueue(2, func(obj interface{}) {
		wr := obj.(t)
		switch wr.data {
		case info:
			logrus.Info(wr.obj)
		case panic:
			logrus.Panic(wr.obj)
		case warning:
			logrus.Warning(wr.obj)
		case perror:
			logrus.Error(wr.obj)
		case fatal:
			logrus.Fatal(wr.obj)
		}
	})
	go func() {
		for {
			select {
			case <-time.After(defaultTime):
				if l, _ := osf.Stat(); l.Size() < int64(size) {
					continue
				}
				log.ch <- struct{}{}
			case <-log.ch:
				log.q.Stop()
				switch algo {
				case Snappy:
					cSnappy(osf)
				case Gzip:
					cGzip(osf)
				}
				osf, err = os.OpenFile(file,
					os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600,
				)
				if err != nil {
					return
				}
				logrus.SetOutput(osf)
				log.q.Resume()
			case <-log.stop:
				osf.Close()
				return
			}
		}
		osf.Close()
	}()

	return log, nil
}

func (log *Log) Compress() {
	log.ch <- struct{}{}
}

func (log *Log) Info(obj interface{}) {
	log.q.Serve(t{info, obj})
}

func (log *Log) Panic(obj interface{}) {
	log.q.Serve(t{panic, obj})
}

func (log *Log) Warning(obj interface{}) {
	log.q.Serve(t{warning, obj})
}

func (log *Log) Error(obj interface{}) {
	log.q.Serve(t{perror, obj})
}

func (log *Log) Fatal(obj interface{}) {
	log.q.Serve(t{fatal, obj})
}

func (log *Log) Close() {
	log.q.Close()
	log.stop <- struct{}{}
}
