package types

import (
	"fmt"
	"math"
	"org.penitence/eureka-sidecar/utils/network"
	"os"
	"strconv"
	"strings"
	"time"
)

func CreateEurekaRegister() (*EurekaInstanceCreate, error) {
	var currentIp string
	var err error
	currentIp = os.Getenv("ip")
	if currentIp == "" {
		currentIp, err = network.FindCurrentIp()
		if err != nil {
			return nil, err
		}
	}

	instanceId, err := createInstanceId(currentIp)
	if err != nil {
		return nil, err
	}
	p, err := strconv.Atoi(findEnvOrDefault("port", "8080"))
	if err != nil {
		return nil, err
	}
	return &EurekaInstanceCreate{
		Instance: EurekaInstance{
			InstanceId:       instanceId,
			App:              strings.ToUpper(os.Getenv("app")),
			HostName:         currentIp,
			IpAddr:           currentIp,
			Status:           UP,
			OverriddenStatus: UNKNOWN,
			Port: EurekaPort{
				Port:    p,
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
				"management.port": string(p),
			},
			HomePageUrl:                   createHomePageUrl(currentIp, p),
			StatusPageUrl:                 createStatusPageUrl(currentIp, p),
			HealthCheckUrl:                createHealthCheckUrl(currentIp, p),
			VipAddress:                    os.Getenv("app"),
			SecureVipAddress:              os.Getenv("app"),
			IsCoordinatingDiscoveryServer: "false",
			LastUpdatedTimestamp:          FindCurrentTimestampToMillisecond(),
			LastDirtyTimestamp:            FindCurrentTimestampToMillisecond(),
		},
	}, nil
}

func createInstanceId(ip string) (string, error) {
	appName := os.Getenv("app")
	var err error
	if appName == "" {
		appName, err = os.Hostname()
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s:%s", appName, ip), nil
}

func findEnvOrDefault(envname string, defaultV string) (v string) {
	v = os.Getenv(envname)
	if v == "" {
		return defaultV
	}
	return
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
