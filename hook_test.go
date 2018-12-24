package asynchook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

type fileHook struct {
	file *os.File
}

func (f *fileHook) Exec(entry *logrus.Entry) error {
	entry.Logger.Formatter = new(logrus.JSONFormatter)
	s, err := entry.String()
	if err != nil {
		return err
	}
	f.file.WriteString(s)
	f.file.Write([]byte("\n"))
	return nil
}

func (f *fileHook) Close() error {
	return f.file.Close()
}

func TestHook(t *testing.T) {
	fileName := fmt.Sprintf("%s/logrus_async_hook", os.TempDir())

	defer os.Remove(fileName)

	file, err := os.Create(fileName)
	if err != nil {
		t.Error(err.Error())
		return
	}

	hook := New(
		SetExec(&fileHook{file}),
		SetFilter(func(entry *logrus.Entry) *logrus.Entry {
			if _, ok := entry.Data["foo2"]; ok {
				delete(entry.Data, "foo2")
			}
			return entry
		}),
		SetExtra(map[string]interface{}{"foo": "bar"}),
	)
	log := logrus.New()
	log.AddHook(hook)

	log.WithField("foo2", "bar").Infof("test foo")
	hook.Flush()

	rfile, err := os.Open(fileName)
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer rfile.Close()

	buf, err := ioutil.ReadAll(rfile)
	if err != nil {
		t.Error(err.Error())
		return
	}

	var item map[string]interface{}
	err = json.Unmarshal(buf, &item)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if item == nil {
		t.Error("Not expected value:nil")
		return
	}

	if reflect.DeepEqual(item["level"], logrus.InfoLevel) {
		t.Errorf("Not expected value:%v", item["level"])
		return
	}

	if !reflect.DeepEqual(item["msg"], "test foo") {
		t.Errorf("Not expected value:%v", item["message"])
		return
	}

	if !reflect.DeepEqual(item["foo"], "bar") {
		t.Errorf("Not expected value:%v", item["foo"])
		return
	}

}
