package log

import (
	"github.com/laconiz/metis/utils/json"
)

type JsonLog struct {
	Level   string          `json:"level"`
	Time    string          `json:"time"`
	Data    json.RawMessage `json:"data,omitempty"`
	Context json.RawMessage `json:"context,omitempty"`
	Message string          `json:"message"`
}

func Json() *JsonFormatter {
	return &JsonFormatter{timeLayout: DefaultTimeLayout}
}

type JsonFormatter struct {
	timeLayout string
}

func (formatter *JsonFormatter) TimeLayout(layout string) *JsonFormatter {
	return &JsonFormatter{timeLayout: layout}
}

func (formatter *JsonFormatter) Format(log *Log) ([]byte, error) {

	var layout string
	if formatter.timeLayout != "" {
		layout = log.Time.Format(formatter.timeLayout)
	} else {
		layout = log.Time.String()
	}

	return json.Marshal(&JsonLog{
		Level:   string(log.Level.Grade()),
		Time:    layout,
		Message: log.Message,
		Data:    log.Data.Raw(),
		Context: log.Context.Raw(),
	})
}
