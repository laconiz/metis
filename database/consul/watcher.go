// consul监视服务

package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api/watch"
)

type Watcher struct {
	client *Client
}

func (watcher *Watcher) Prefix(prefix string, handler watch.HandlerFunc) (Plan, error) {

	if handler == nil {
		return nil, fmt.Errorf("nil handler")
	}

	plan, err := watch.Parse(param{"type": "keyprefix", "prefix": prefix})
	if err != nil {
		return nil, err
	}
	plan.Handler = handler

	go plan.Run(watcher.client.addr)
	return plan, nil
}

type param = map[string]interface{}

type Plan interface {
	Stop()
}
