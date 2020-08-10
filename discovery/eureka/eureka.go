package eureka

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"org.penitence/eureka-sidecar/discovery"
	_ "org.penitence/eureka-sidecar/env"
	"org.penitence/eureka-sidecar/types"
	"strconv"
	"strings"
	"time"
)

const ApplicationJson = "application/json"

type Register struct {
	BaseUrl  string
	Instance *types.EurekaInstanceCreate
	client   *http.Client
}

func CreateRegister(baseUrl string, create *types.EurekaInstanceCreate) *Register {
	c := &http.Client{
		Timeout: 5 * time.Second,
	}
	return &Register{
		BaseUrl:  baseUrl,
		Instance: create,
		client:   c,
	}
}

func (r *Register) IsAppExists() (bool, error) {
	resp, err := r.client.Get(r.createAppUrl())
	if err != nil {
		return false, err
	}
	if resp.StatusCode == 404 {
		return false, nil
	}
	if resp.StatusCode == 401 {
		return false, errors.New("eureka的验证信息错误")
	}
	return isSuccess(resp.StatusCode), nil
}

func (r *Register) CreateInstance() (bool, error) {
	appUrl := r.createAppUrl()
	body, err := json.Marshal(r.Instance)
	if err != nil {
		return false, err
	}
	reps, err := r.client.Post(appUrl, ApplicationJson, bytes.NewReader(body))
	if err != nil {
		return false, err
	}
	b, _ := ioutil.ReadAll(reps.Body)
	log.Infof("创建app应用实例, 响应code为: %d", reps.StatusCode)
	if isSuccess(reps.StatusCode) {
		return true, nil
	}
	return false, errors.New(string(b))
}

func (r *Register) Heartbeat() (bool, error) {
	reqUrl := r.createUpdateUrl(types.UP, strconv.FormatInt(r.Instance.Instance.LastDirtyTimestamp, 10))
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPut, reqUrl, nil)
	if err != nil {
		return false, err
	}
	reps, err := r.client.Do(req)
	if err != nil {
		return false, err
	} else {
		b, _ := ioutil.ReadAll(reps.Body)
		log.Infof("更新app实例, 响应code为: %d", reps.StatusCode)
		if isSuccess(reps.StatusCode) {
			return true, nil
		} else {
			return false, discovery.NewError(string(b), reps.StatusCode)
		}
	}
}

func (r *Register) RemoveInstance() (bool, error) {
	reqUrl := r.createInstanceUrl()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, reqUrl, nil)
	if err != nil {
		return false, err
	}
	reps, err := r.client.Do(req)
	if err != nil {
		return false, nil
	} else {
		b, _ := ioutil.ReadAll(reps.Body)
		log.Infof("删除app实例, 响应code为: %d", reps.StatusCode)
		if isSuccess(reps.StatusCode) {
			return true, nil
		} else {
			return false, errors.New(string(b))
		}
	}
}

func isSuccess(code int) bool {
	return code >= 200 && code <= 299
}

func (r *Register) createAppUrl() string {
	return fmt.Sprintf("%s/apps/%s", r.BaseUrl, strings.ToUpper(r.Instance.Instance.App))
}

func (r *Register) createInstanceUrl() string {
	return fmt.Sprintf("%s/apps/%s/%s", r.BaseUrl, strings.ToUpper(r.Instance.Instance.App), r.Instance.Instance.InstanceId)
}

func (r *Register) createUpdateUrl(status types.EurekaStatus, lastDirtyTimestamp string) string {
	app := r.Instance.Instance.App
	instantId := r.Instance.Instance.InstanceId
	reqUrl := fmt.Sprintf("%s/apps/%s/%s", r.BaseUrl, strings.ToUpper(app), instantId)
	u, _ := url.Parse(reqUrl)
	params := url.Values{}
	params.Set("status", status.String())
	params.Set("lastDirtyTimestamp", lastDirtyTimestamp)
	u.RawQuery = params.Encode()
	return u.String()
}
