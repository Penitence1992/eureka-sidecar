package main

import (
	log "github.com/sirupsen/logrus"
	"org.penitence/eureka-sidecar/discovery"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	gitCommit  = "1"
	buildStamp = "1900-01-01"
)

func main() {
	printVersionInfo()
	clients, err := discovery.CreateDiscover()
	if err != nil {
		panic(err)
	}
	period := 10 * time.Second
	done := make(chan os.Signal)
	taskGroup := &sync.WaitGroup{}
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	runAsync(clients, period, taskGroup)
	for {
		select {
		case sg := <-done:
			log.Infof("主线程接收到信号:%s\n", sg)
			taskGroup.Wait()
			os.Exit(0)
		}
	}
}

func runAsync(clients []discovery.Discovery, period time.Duration, group *sync.WaitGroup) {
	for _, client := range clients {
		rsignal := make(chan os.Signal)
		signal.Notify(rsignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		executor := discovery.CreateTimerExecutor(client, rsignal, period)
		go func() {
			group.Add(1)
			defer group.Done()
			executor.Start()
		}()
	}
}

func printVersionInfo() {
	log.Infoln("Starting Web hook handler --->")
	log.Infof("Git Commit : %s\n", gitCommit)
	log.Infof("Build Stamp : %s\n", buildStamp)
}
