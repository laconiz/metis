package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"strconv"
	"strings"
)

type Register struct {
	client *Client
}

func (rg *Register) Register(svcName, address string, meta map[string]string) (reg *api.AgentServiceRegistration, err error) {
	addrPair := strings.Split(address, ":")
	port, _ := strconv.Atoi(addrPair[1])

	reg = &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s_%s_%s", svcName, addrPair[0], addrPair[1]),
		Name:    svcName,
		Address: addrPair[0],
		Port:    port,
		Meta:    meta,
	}

	if reg.Check == nil {

		var addr string
		if reg.Address == "" {
			addr = "127.0.0.1"
		} else {
			addr = reg.Address
		}

		reg.Check = &api.AgentServiceCheck{
			TCP:                            fmt.Sprintf("%s:%d", addr, reg.Port),
			Interval:                       "10s",
			Timeout:                        "30s",
			DeregisterCriticalServiceAfter: "30s",
		}
	}

	return reg, rg.client.Agent().ServiceRegister(reg)
}

func (rg *Register) Deregister(svciD string) {
	rg.client.Agent().ServiceDeregister(svciD)
}
