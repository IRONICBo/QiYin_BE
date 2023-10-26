package log

import (
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/IRONICBo/QiYin_BE/pkg/utils"
)

type filelineHook struct{}

func newFilelineHook() *filelineHook {
	return &filelineHook{}
}

func (hook *filelineHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *filelineHook) Fire(entry *logrus.Entry) error {
	var s string
	_, file, line, _ := runtime.Caller(8)
	i := strings.SplitAfter(file, "/")
	if len(i) > 3 {
		s = i[len(i)-3] + i[len(i)-2] + i[len(i)-1] + ":" + utils.IntToString(line)
	}
	entry.Data["FilePath"] = s

	return nil
}
