package logging

import (
	"fmt"
	"github.com/lin07ux/go-gin-example/pkg/file"
	"github.com/lin07ux/go-gin-example/pkg/setting"
	"os"
	"time"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

func getLogFileName() string {
	return fmt.Sprintf(
		"%s-%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt,
	)
}

func openLogFile(filename, filepath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filepath
	if file.IsPermissionDenied(src) {
		return nil, fmt.Errorf("permission denied src: %s", src)
	}

	if err := file.MakeDirIfNotExist(src); err != nil {
		return nil, fmt.Errorf("failed to make directory, src: %s, err: %v", src, err)
	}

	f, err := file.Open(src + filename, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	return f, nil
}
