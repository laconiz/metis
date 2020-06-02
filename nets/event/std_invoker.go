package event

import (
	"github.com/laconiz/metis/nets/packet"
)

func NewStdInvoker() *StdInvoker {
	return &StdInvoker{handlers: map[packet.MetaID]func(*Event){}}
}

type StdInvoker struct {
	handlers map[packet.MetaID]func(*Event)
}

func (invoker *StdInvoker) Invoke(event *Event) {

	handler, ok := invoker.handlers[event.Meta.ID()]
	if !ok {
		return
	}

	handler(event)
}

func (invoker *StdInvoker) Register(msg interface{}, handler func(*Event)) {

	meta := packet.MetaByMsg(msg)
	if meta == nil {
		return
	}

	invoker.handlers[meta.ID()] = handler
}
