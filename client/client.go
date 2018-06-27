package main

import (
	"bytes"
	"flag"
	"fmt"
	"os/exec"
	"time"

	"strings"

	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/parnurzeal/gorequest"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func main() {
	var host string
	var id string
	var port string
	var changeInterval int

	flag.StringVar(&host, "host", "", "server host")
	flag.StringVar(&id, "id", "", "client id")
	flag.StringVar(&port, "port", "", "server host")
	flag.IntVar(&changeInterval, "changeInterval", 180, "change IP minutes")

	flag.Parse()

	if host == "" {
		log.Fatal("host 必填")
	}
	if id == "" {
		log.Fatal("id 必填")
	}
	if port == "" {
		log.Fatal("port 必填")
	}

	for {
		//发起请求删除
		err := sendDelete(host, id)
		if err != nil {
			log.WithError(err).Error("删除代理错误")
		} else {
			log.Info("删除代理成功")
		}

		err = updateIP()
		if err != nil {
			log.WithError(err).Error("更新IP错误")
			time.Sleep(5 * time.Second)
			continue
		}
		log.Info("更新IP成功")

		//发送请求
		err = sendUpdate(host, id, port)
		if err != nil {
			log.WithError(err).Error("发送代理错误")
			time.Sleep(3 * time.Second)
			continue
		}
		log.Info("发送代理成功")
		time.Sleep(time.Duration(changeInterval) * time.Second)
	}
}

func sendDelete(host, id string) error {
	u := fmt.Sprintf("%s/del?id=%s", host, id)
	err := req(u)
	if err != nil {
		return err
	}
	return nil
}

func startProxy(port string) {
	proxy := goproxy.NewProxyHttpServer()
	log.Fatal(http.ListenAndServe(":"+port, proxy))
}

func sendUpdate(host, id, port string) error {
	u := fmt.Sprintf("%s/add?id=%s&port=%s", host, id, port)
	err := req(u)
	if err != nil {
		return err
	}
	return nil
}

func req(u string) error {
	client := gorequest.New()
	_, body, errs := client.Get(u).End()
	if len(errs) > 0 {
		return fmt.Errorf("请求失败: %s", errs)
	}
	if gjson.Get(body, "status").Str == "error" {
		return fmt.Errorf("请求失败: %s", gjson.Get(body, "error").Str)
	}
	return nil
}

func updateIP() error {
	res, err := execShell("pppoe-stop")
	log.Debug(err, res)

	res, err = execShell("pppoe-start")
	if err != nil {
		return fmt.Errorf("pppoe-start: %s", err)
	}
	log.Debug(res)
	time.Sleep(1 * time.Second)

	res, err = execShell("pppoe-status")
	if err != nil {
		return fmt.Errorf("pppoe-status: %s", err)
	}
	log.Debug(res)
	if strings.Contains(res, "Link is down") {
		return fmt.Errorf("pppoe-status get error")
	}
	return nil
}

//阻塞式的执行外部shell命令的函数,等待执行完毕并返回标准输出
func execShell(s string) (string, error) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", s)

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
