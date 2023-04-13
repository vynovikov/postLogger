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

func NewSaver(rootPath string, limit int64) (*SaverStruct, error) {
	//logger.L.Infof("saver.NewSaver invoked with rootPath: %s, limit %d\n", rootPath, limit)

	_, err := os.Stat(rootPath)
	if err != nil {
		if os.IsNotExist(err) {

			// Log folder not found

			os.Mkdir(rootPath, 0777)
			//logger.L.Infoln("in store.NewStore folder created")

			return &SaverStruct{path: rootPath, limit: limit}, nil
		}

		//logger.L.Errorf("in store.NewStore error %v\n", err)
		return nil, fmt.Errorf("in saver.NewSaver error finding \"%s\": %v\n", rootPath, err)
	}
	// Log folder found

	filePath := ""

	err = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		//logger.L.Infof("in saver.NewSaver checking path: %s, d name: %s\n", path, d.Name())
		if err != nil {
			return fmt.Errorf("in saver.NewSaver unable to walk through %q: %v\n", path, err)
		}

		if !d.IsDir() {
			//name := path[strings.LastIndex(path, "\\")+1:]
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

	//return nil, nil
	if len(filePath) > 0 {
		rootPath = filePath
	}
	ss := &SaverStruct{path: rootPath, limit: limit}
	//logger.L.Infof("in saver.NewSaver returning %v\n", ss)
	return ss, nil
}

func (s *SaverStruct) Save(l string) error {

	//logger.L.Infof("store.Save invoked with len(l) = %d, s.path: %q\n", len(l), s.path)

	f, pathUpd, err := GetFile(s.path, s.limit)
	if err != nil {
		return fmt.Errorf("in store.Save error while getting file %v\n", err)
	}

	//logger.L.Infof("in store.Save path %q pathUpd %q\n", s.path, pathUpd)

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
	//logger.L.Infof("in store.IsTS %s matched? %t\n", s, res)
	return res
}

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

	//logger.L.Infof("in saver.GetFile path: %s, lastPartPath %q\n", path, lastPartPath)

	if strings.Contains(lastPartPath, ".txt") &&
		IsTS(lastPartPath) {
		//logger.L.Infof("in store.GetFile opening %q\n", path)
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
		//pathUpd = path[:strings.Index(path, "/logs/")+len("/logs/")]
		pathUpd = path[:len(path)-len(lastPartPath)]
		//logger.L.Infof("in store.GetFile pathUpd 1 %q\n", pathUpd)
	} else {
		pathUpd = path + sep
	}

	pathUpd += repo.NewTS() + ".txt"
	//logger.L.Infof("in store.GetFile pathUpd 2 %q\n", pathUpd)

	f, err := os.Create(pathUpd)

	if err != nil {
		//logger.L.Errorf("in store.GetFile error %v\n", err)
		return nil, "", err
	}

	return f, pathUpd, nil

}
func (s *SaverStruct) SetPath(p string) {
	s.path = p
}
