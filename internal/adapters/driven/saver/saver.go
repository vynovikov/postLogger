package saver

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"postLogger/internal/logger"
	"postLogger/internal/repo"
	"regexp"
	"strings"
)

type Saver interface {
	Save(string) error
	SetPath(string)
}

type SaverStruct struct {
	path  string
	limit int64
}

// NewSaver creates new example of *StoreStruct.
// If logs folder is not exists, NewSaver creates it.
// Tested in saver_test.go
func NewSaver(rootPath string, limit int64) (*SaverStruct, error) {

	_, err := os.Stat(rootPath)
	if err != nil {
		if os.IsNotExist(err) {

			// Log folder not found

			os.Mkdir(rootPath, 0777)

			return &SaverStruct{path: rootPath, limit: limit}, nil
		}
		return nil, fmt.Errorf("in saver.NewSaver error finding \"%s\": %v\n", rootPath, err)
	}
	// Log folder found

	filePath := ""

	err = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("in saver.NewSaver unable to walk through %q: %v\n", path, err)
		}

		if !d.IsDir() {
			if IsTS(d.Name()) {

				info, err := d.Info()
				if err != nil {
					return fmt.Errorf("in saver.NewSaver unable to check info of \"%s\": %v\n", path, err)
				}
				if info.Size() < limit {
					filePath = path
					return nil
				}
				return nil
			}
		}
		return nil
	})
	if err != nil {
		logger.L.Errorf("in saver.NewSaver error: %v\n", err)
	}
	if len(filePath) > 0 {
		rootPath = filePath
	}
	ss := &SaverStruct{path: rootPath, limit: limit}
	return ss, nil
}

// Save adds l to log file
func (s *SaverStruct) Save(l string) error {

	f, pathUpd, err := GetFile(s.path, s.limit)
	if err != nil {
		return fmt.Errorf("in store.Save error while getting file %v\n", err)
	}
	if pathUpd == s.path {
		_, err = f.WriteString("\r\n")
		if err != nil {
			return fmt.Errorf("in store.Save error while writing CRLF %v\n", err)
		}
	}
	if pathUpd != s.path {
		s.SetPath(pathUpd)
	}
	_, err = f.WriteString(l)
	if err != nil {
		return fmt.Errorf("in store.Save error while writing logstring %v\n", err)
	}
	err = f.Close()
	if err != nil {
		return fmt.Errorf("in store.Save error while closingl %v\n", err)
	}
	return nil

}

func IsTS(s string) bool {
	r := regexp.MustCompile(`^[0-2]\d.[0-2]\d.[0-2]\d\d\d [0-2]\d_\d\d_\d\d.\d{3,4}`)
	res := r.MatchString(s)
	return res
}

// GetFile looks for valid log file.
// Creates one if not found.
func GetFile(path string, limit int64) (*os.File, string, error) {

	pathUpd, lastPartPath, sep := "", "", ""

	switch {
	case !strings.Contains(path, "\\"):
		sep = "/"
		lastPartPath = path[strings.LastIndex(path, "/")+1:]
	default:
		sep = "\\"
		lastPartPath = path[strings.LastIndex(path, "\\")+1:]
	}

	if strings.Contains(lastPartPath, ".txt") &&
		IsTS(lastPartPath) {
		fileStat, err := os.Stat(path)
		if err != nil {
			return nil, "", err
		}
		if fileStat.Size() < limit {
			f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				return nil, "", err
			}
			return f, path, nil
		}
		pathUpd = path[:len(path)-len(lastPartPath)]
	} else {
		pathUpd = path + sep
	}

	pathUpd += repo.NewTS() + ".txt"

	f, err := os.Create(pathUpd)

	if err != nil {
		return nil, "", err
	}

	return f, pathUpd, nil

}
func (s *SaverStruct) SetPath(p string) {
	s.path = p
}
