

package main

import (
"fmt"
"github.com/afex/hystrix-go/hystrix"
"log"
"net"
"net/http"
"os"
"os/signal"
"syscall"
"time"
)

func main() {
	hystrix.ConfigureCommand("default", hystrix.CommandConfig{
		Timeout:                1000, // 单次请求 超时时间
		MaxConcurrentRequests:  1,    // 最大并发量
		SleepWindow:            5000, // 熔断后多久去尝试服务是否可用
		RequestVolumeThreshold: 1,    // 验证熔断的 请求数量, 10秒内采样
		ErrorPercentThreshold:  1,    // 验证熔断的 错误百分比
	})

	//开启一个http监控服务
	//可以使用hystrix-dashboard面板查看具体情况
	//https://github.com/mlabouardy/hystrix-dashboard-docker
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go func() {
		err := http.ListenAndServe(net.JoinHostPort("", "8888"), hystrixStreamHandler)
		log.Fatal(err)
	}()

	for i := 0; i < 100000; i++ {
		Do(i)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func Do(params int) {
	err := hystrix.Do("defalut", func() error {
		fmt.Println(params)
		return nil
	}, nil)

	if err != nil {
		//加入自动降级处理，如获取缓存数据等
		switch err {
		case hystrix.ErrCircuitOpen:
			fmt.Println("circuit error:" + err.Error())
		case hystrix.ErrMaxConcurrency:
			fmt.Println("circuit error:" + err.Error())
		default:
			fmt.Println("circuit error:" + err.Error())
		}

		time.Sleep(1 * time.Second)
		log.Println("sleep 1 second")
	}
}
