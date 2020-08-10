package factory

import (
	log "github.com/sirupsen/logrus"
	"org.penitence/eureka-sidecar/discovery/eureka"
	"org.penitence/eureka-sidecar/env"
	"org.penitence/eureka-sidecar/types"
	"org.penitence/eureka-sidecar/utils/network"
	"strconv"
	"strings"
)

const (
	zoneUrlEnv = "zoneUrl"
	appEnv     = "app"
	ipEnv      = "ip"
	portEnv    = "port"
)

func CreateEurekaDiscovers() ([]*eureka.Register, error) {
	zoneUrl, err := env.FindEnv(zoneUrlEnv)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(env.FindEnvOrDefault(portEnv, "8080"))
	if err != nil {
		return nil, err
	}
	appName := env.FindEnvOrDefault(appEnv, "")
	ip := env.FindEnvOrDefaultSupplier(ipEnv, ipSupplier)
	instance, err := types.CreateEurekaRegister(ip, port, appName)
	if err != nil {
		return nil, err
	}

	urls := strings.Split(zoneUrl, ",")

	clients := make([]*eureka.Register, len(urls))

	for i, baseUrl := range urls {
		clients[i] = eureka.CreateRegister(baseUrl, instance)
	}

	return clients, nil
}

func ipSupplier() string {
	ip, err := network.FindCurrentIp()
	if err != nil {
		log.Warningf("获取当前ip失败:%v, 使用127.0.0.1这个ip", err)
		ip = "127.0.0.1"
	}
	return ip

}
