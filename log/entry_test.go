package log

import (
	"errors"
	"testing"
)

func TestEntry(t *testing.T) {

	logger := std

	type Data struct {
		Int int    `json:"int"`
		Str string `json:"string"`
	}

	value0 := &Data{Int: 0, Str: "hello"}
	value1 := &Data{Int: 1, Str: "world"}
	value2 := &Data{Int: 2, Str: "Go!"}
	err := errors.New("text error")

	logger.Debug("hello world")
	logger.Debugf("running on %s", "Go")

	logger.Data(value0).Info("value 0")
	logger.Data(value1).Infof("value %d", 1)

	logger = logger.Field("module", "test")

	logger.Data(value0, value1, value2).Warn("values 0, 1, 2")
	logger.Data(value0, value1).Warnf("values %d, %d", 0, 1)

	logger.Data(value0, err).Error("value with error")
	logger.Data(value0, value1, value2, err).Errorf("%s with %s", "values", "error")

	logger = logger.Data(value0, value1, err)
	logger.Fatal("fatal")
	logger.Fatalf("%s", "fatalf")
}
