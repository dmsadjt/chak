package indexer

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)


type DirectoryScanner struct {
	watchedDirs []string
	allowedExts []string
	maxFileSize int64
}

func NewDirectoryScanner(dirs []string, exts []string, fileSize int64) *DirectoryScanner {
	return &DirectoryScanner{
		watchedDirs: dirs,
		allowedExts: exts,
		maxFileSize: fileSize,
	}
}

func (dirMgr *DirectoryScanner) ScanDirectories() ([]FileInfo, error) {
	var files []FileInfo

	for _, dir := range dirMgr.watchedDirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if info.IsDir() {
				return nil
			}

			if info.Size() > dirMgr.maxFileSize {
				return nil
			}

			ext := strings.ToLower(filepath.Ext(path))
			if !dirMgr.isAllowedExtension(ext) {
				return nil 
			}
			
			hash, err := hashFile(path)
			if err != nil {
				return nil 
			}

			files = append(files, FileInfo{
				SzPath: path,
				SzName: info.Name(),
				SzExtension: ext,
				InSize: info.Size(),
				TmModTime: info.ModTime(),
				SzHash: hash,
			})

			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("Error scanning %s: %w", dir, err)
		}
	}
	return files, nil
}

func (dirMgr *DirectoryScanner) isAllowedExtension(ext string) bool {
	if len(dirMgr.allowedExts) == 0 {
		return true
	}

	for _, allowedExt := range dirMgr.allowedExts {
		if ext == allowedExt {
			return true
		}
	}

	return false
}

func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
