package main

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"org.penitence/eureka-sidecar/register"
	"org.penitence/eureka-sidecar/runner"
	"org.penitence/eureka-sidecar/types"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)


var (
	gitCommit = "1"
	buildStamp = "1900-01-01"
)

func main() {
	eurekaHome, err := findEnv("zoneUrl")
	if err != nil {
		panic(err)
	}
	period := 10 * time.Second
	done := make(chan os.Signal)
	var taskGroup sync.WaitGroup
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	instance, err := types.CreateEurekaRegister()
	if err != nil {
		panic(err)
	}
	sc, err := register.InitRegister(eurekaHome, instance)
	if err != nil {
		log.Fatalf("register app fail:%v", err)
	}
	for {
		select {
		case u := <-sc:
			go func() {
				taskGroup.Add(1)
				defer taskGroup.Done()
				rsignal := make(chan os.Signal)
				signal.Notify(rsignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
				runner.CreateTimeRunner(rsignal,
					period,
					register.CreateRunFunc(u, instance),
					register.CreateDoneFunc(u, instance),
				).Start()
			}()
			break
		case sg := <-done:
			fmt.Printf("accept signal %s", sg)
			taskGroup.Wait()
			os.Exit(0)
		}
	}
}

func findEnv(envname string) (v string, err error) {
	v = os.Getenv(envname)
	if v == "" {
		return "", errors.New("env not found")
	}
	return
}

func printVersionInfo()  {
	log.Infoln("Starting Web hook handler --->")
	log.Infof("Git Commit : %s\n", gitCommit)
	log.Infof("Build Stamp : %s\n", buildStamp)
}
