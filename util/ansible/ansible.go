package utilansible

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"path"
)

//const logpath = "/opt/ansible/"
const logpath = "/Users/rubinus/app/go-check-k8s/"

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func CreateAnsibleLogWriterWithId(clusterName string, logId string) (io.Writer, error) {
	dirName := path.Join(logpath, clusterName)
	if !Exists(dirName) {
		err := os.MkdirAll(dirName, 0755)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("open ansible log file with id failed: %v", err))
		}
	}
	fileName := path.Join(dirName, fmt.Sprintf("%s.log", logId))
	writer, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("open ansible log file with id failed: %v", err))
	}
	return writer, nil
}
