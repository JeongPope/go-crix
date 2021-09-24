package logger

import (
	"errors"
	"os"
	"testing"
)

func Test_Logger(t *testing.T) {
	f, _ := os.Create("logger.log")
	Log.SetOut(f)
	Log.SetLevel(DEBUG)
	Log.Debug("This is debug log")
	Log.Debugf("%.8f", 0.1808195842)
	Log.Info("This is info log")
	Log.Warn(errors.New("Test error"))
	//Log.Fatal("This is fatal log")
	//Log.Panicf("%s","This is panic log")
}
