package leveldb

import "fmt"

type Logger interface {
	Logv(format string, a ...interface{})
}

type defaultLogger struct{
	WritableFile
}

func makeLogger(fname string) Logger {

	return &defaultLogger {}
}

func (this *defaultLogger) Logv(format string, a ...interface{}) {

	str := fmt.Sprintf(format, a)
	
	this.Append(str)
	
	this.Flush()
}