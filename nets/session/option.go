package session

import (
	"github.com/laconiz/metis/nets/cipher"
	"github.com/laconiz/metis/nets/encoder"
	"github.com/laconiz/metis/nets/event"
	"time"
)

// 会话配置信息
type Option struct {
	Queue   int
	Cipher  cipher.Maker
	Encoder encoder.Maker
	Invoker event.Invoker
	Timeout time.Duration
}

func (opt *Option) Parse() {

	if opt.Cipher == nil {
		opt.Cipher = cipher.EmptyMaker{}
	}

	if opt.Encoder == nil {
		opt.Encoder = encoder.NameMaker{}
	}

	if opt.Invoker == nil {
		opt.Invoker = event.NewStdInvoker()
	}

	if opt.Timeout <= 0 {
		opt.Timeout = time.Second * 15
	}
}
