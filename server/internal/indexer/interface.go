package indexer

import "time"

type FileInfo struct {
	SzPath      string
	SzName      string
	SzExtension string
	InSize      int64
	TmModTime   time.Time
	SzHash		string
}

type ScannerInterface interface {
	ScanDirectories() ([]FileInfo, error)
}

type ManagerInterface interface {
	IndexAll() error
	StartWatcher(tmInterval time.Duration)
}

