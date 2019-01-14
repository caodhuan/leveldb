package leveldb

import "fmt"

const (
	kLogFile = iota
	kDBLockFile
	kTableFile
	kDescriptorFile
	kCurrentFile
	kTempFile
	kInfoLogFile  // Either the current one, or an old one
)

func makeFileName(name string, number uint64, suffix string) string {
	buf := fmt.Sprintf("/%06llu.%s", number, suffix)
	return name + buf;
}

// Return the name of the log file with the specified number
// in the db named by "dbname".  The result will be prefixed with
// "dbname".
func LogFileName(name string, number uint64) string {
	return makeFileName(name, number, "log");
}

// Return the name of the info log file for "dbname".
func InfoLogFileName(name string) string {
	return name + "/LOG";
}

// Return the name of the old info log file for "dbname".
func OldInfoLogFileName(name string) string {
	return name + "/LOG.old";
}