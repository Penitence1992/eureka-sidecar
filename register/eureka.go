package register

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	_ "org.penitence/eureka-sidecar/env"
	"org.penitence/eureka-sidecar/runner"
	"org.penitence/eureka-sidecar/types"
	"strconv"
	"strings"
	"time"
)

const ApplicationJson = "application/json"

// register app to eureka
//zone url must be http://abc:abc@127.0.0.1:8761/eureka format
func InitRegister(zoneUrl string, create *types.EurekaInstanceCreate) (<-chan string, error) {
	body, err := json.Marshal(create)
	if err != nil {
		return nil, err
	}
	sUrlChan := make(chan string)
	for _, u := range strings.Split(zoneUrl, ",") {
		reqUrl := createRequestUrl(u, create.Instance.App)
		successUrl := u
		go func() {
			for !RequestCreate(reqUrl, body) {
				log.Infof("wait 5 seconds and retry")
				time.Sleep(5 * time.Second)
			}
			sUrlChan <- successUrl
		}()
	}
	return sUrlChan, nil
}

func RequestCreate(reqUrl string, body []byte) bool {
	reps, err := http.Post(reqUrl, ApplicationJson, bytes.NewReader(body))
	if err != nil {
		log.Errorf("client request error %v", err)
		return false
	}
	b, _ := ioutil.ReadAll(reps.Body)
	log.Infof("client request response , code: %d, context: %s", reps.StatusCode, string(b))

	if isSuccess(reps.StatusCode) {
		log.Infof("register success")
		return true
	}
	return false
}

func CreateDoneFunc(url string, create *types.EurekaInstanceCreate) runner.RunFunc {
	return func() {
		reqUrl := createDoneUrl(url, create.Instance.App, create.Instance.InstanceId)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, reqUrl, nil)
		if err != nil {
			log.Errorf("client request error %v", err)
			return
		}
		reps, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Errorf("client request error %v", err)
			return
		} else {
			b, _ := ioutil.ReadAll(reps.Body)
			log.Infof("delete instance request response , code: %d, context: %s", reps.StatusCode, string(b))
		}
	}
}

func CreateRunFunc(url string, create *types.EurekaInstanceCreate) runner.RunFunc {
	return func() {
		reqUrl := createUpdateUrl(url, create.Instance.App, create.Instance.InstanceId, types.UP, strconv.FormatInt(create.Instance.LastDirtyTimestamp, 10))
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPut, reqUrl, nil)

		if err != nil {
			log.Errorf("client request error %v", err)
			return
		}
		reps, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Errorf("client request error %v", err)
			return
		} else {
			b, _ := ioutil.ReadAll(reps.Body)
			log.Infof("update instance request response , code: %d, context: %s", reps.StatusCode, string(b))
		}
	}
}

func isSuccess(code int) bool {
	return code >= 200 && code <= 299
}

func createRequestUrl(u, app string) string {
	return fmt.Sprintf("%s/apps/%s", u, strings.ToUpper(app))
}

func createDoneUrl(base, app, instantId string) string {
	return fmt.Sprintf("%s/apps/%s/%s", base, strings.ToUpper(app), instantId)
}

func createUpdateUrl(base, app, instantId string, status types.EurekaStatus, lastDirtyTimestamp string) string {
	reqUrl := fmt.Sprintf("%s/apps/%s/%s", base, strings.ToUpper(app), instantId)
	u, _ := url.Parse(reqUrl)
	params := url.Values{}
	params.Set("status", status.String())
	params.Set("lastDirtyTimestamp", lastDirtyTimestamp)
	u.RawQuery = params.Encode()
	return u.String()
}
