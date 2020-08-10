package types

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

func CreateEurekaRegister(ip string, port int, appName string) (*EurekaInstanceCreate, error) {
	var err error
	if ip == "" {
		return nil, errors.New("ip不能为空")
	}
	if port <= 0 || port > 65535 {
		return nil, errors.New("端口访问错误")
	}
	instanceId, err := createInstanceId(ip, appName)
	if err != nil {
		return nil, err
	}
	return &EurekaInstanceCreate{
		Instance: EurekaInstance{
			InstanceId:       instanceId,
			App:              strings.ToUpper(os.Getenv("app")),
			HostName:         ip,
			IpAddr:           ip,
			Status:           UP,
			OverriddenStatus: UNKNOWN,
			Port: EurekaPort{
				Port:    port,
				Enabled: true,
			},
			SecurePort: EurekaPort{
				Port:    443,
				Enabled: false,
			},
			CountryId: 1,
			DataCenterInfo: EurekaDataCenterInfo{
				Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
				Name:  "MyOwn",
			},
			LeaseInfo: EurekaLeaseInfo{
				RenewalIntervalInSecs: 5,
				DurationInSecs:        20,
				RegistrationTimestamp: FindCurrentTimestampToMillisecond(),
				ServiceUpTimestamp:    FindCurrentTimestampToMillisecond(),
			},
			Metadata: map[string]string{
				"management.port": string(port),
			},
			HomePageUrl:                   createHomePageUrl(ip, port),
			StatusPageUrl:                 createStatusPageUrl(ip, port),
			HealthCheckUrl:                createHealthCheckUrl(ip, port),
			VipAddress:                    appName,
			SecureVipAddress:              appName,
			IsCoordinatingDiscoveryServer: "false",
			LastUpdatedTimestamp:          FindCurrentTimestampToMillisecond(),
			LastDirtyTimestamp:            FindCurrentTimestampToMillisecond(),
		},
	}, nil
}

func createInstanceId(ip, appName string) (string, error) {
	var err error
	if appName == "" {
		appName, err = os.Hostname()
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s:%s", appName, ip), nil
}

func FindCurrentTimestampToMillisecond() int64 {
	return time.Now().UnixNano() / int64(math.Pow10(6))
}

func createHomePageUrl(ip string, port int) string {
	return fmt.Sprintf("http://%s:%d/", ip, port)
}

func createStatusPageUrl(ip string, port int) string {
	return fmt.Sprintf("http://%s:%d/actuator/info", ip, port)
}

func createHealthCheckUrl(ip string, port int) string {
	return fmt.Sprintf("http://%s:%d/actuator/health", ip, port)
}
