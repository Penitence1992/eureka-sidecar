package types

import "bytes"

type EurekaStatus int

const (
	UP EurekaStatus = iota
	DOWN
	UNKNOWN
)

func (e EurekaStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(e.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e EurekaStatus) String() string {
	switch e {
	case 0:
		return "UP"
	case 1:
		return "DOWN"
	case 2:
		return "UNKNOWN"
	default:
		return ""
	}
}



type EurekaInstanceCreate struct {
	Instance EurekaInstance `json:"instance"`
}

type EurekaInstance struct {
	InstanceId                    string               `json:"instanceId"`
	HostName                      string               `json:"hostName"`
	App                           string               `json:"app"`
	IpAddr                        string               `json:"ipAddr"`
	Status                        EurekaStatus         `json:"status"`
	OverriddenStatus              EurekaStatus         `json:"overriddenStatus"`
	Port                          EurekaPort           `json:"port"`
	SecurePort                    EurekaPort           `json:"securePort"`
	CountryId                     int                  `json:"countryId"`
	DataCenterInfo                EurekaDataCenterInfo `json:"dataCenterInfo"`
	LeaseInfo                     EurekaLeaseInfo      `json:"leaseInfo"`
	Metadata                      map[string]string    `json:"metadata"`
	HomePageUrl                   string               `json:"homePageUrl"`
	StatusPageUrl                 string               `json:"statusPageUrl"`
	HealthCheckUrl                string               `json:"healthCheckUrl"`
	VipAddress                    string               `json:"vipAddress"`
	SecureVipAddress              string               `json:"secureVipAddress"`
	IsCoordinatingDiscoveryServer string               `json:"isCoordinatingDiscoveryServer"`
	LastUpdatedTimestamp          int64                `json:"lastUpdatedTimestamp"`
	LastDirtyTimestamp            int64                `json:"lastDirtyTimestamp"`
}

type EurekaPort struct {
	Port    int  `json:"$"`
	Enabled bool `json:"@enabled"`
}

type EurekaSecurePort struct {
	EurekaPort
}

type EurekaDataCenterInfo struct {
	Class string `json:"@class"`
	Name  string `json:"name"`
}

type EurekaLeaseInfo struct {
	RenewalIntervalInSecs int   `json:"renewalIntervalInSecs"`
	DurationInSecs        int   `json:"durationInSecs"`
	RegistrationTimestamp int64 `json:"registrationTimestamp"`
	LastRenewalTimestamp  int64 `json:"lastRenewalTimestamp"`
	EvictionTimestamp     int64 `json:"evictionTimestamp"`
	ServiceUpTimestamp    int64 `json:"serviceUpTimestamp"`
}

// 更新LeaseInfo中的时间内容
func (i *EurekaLeaseInfo) Renew() {
	i.RegistrationTimestamp = FindCurrentTimestampToMillisecond()
	i.ServiceUpTimestamp = FindCurrentTimestampToMillisecond()
}
