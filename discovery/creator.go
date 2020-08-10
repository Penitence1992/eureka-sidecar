package discovery

import (
	"org.penitence/eureka-sidecar/discovery/factory"
)

func CreateDiscover() ([]Discovery, error) {
	//TODO 后续可能添加更多的discovery模块
	clients, err := factory.CreateEurekaDiscovers()
	if err != nil {
		return nil, err
	}
	discoverys := make([]Discovery, len(clients))
	for i, c := range clients {
		discoverys[i] = c
	}
	return discoverys, nil
}
