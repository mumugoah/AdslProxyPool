package main

import (
	"os/exec"
	"bytes"
	log "github.com/sirupsen/logrus"
	"fmt"
	"time"
	"flag"
	"github.com/parnurzeal/gorequest"
)

func main()  {
	host := flag.String("host","", "server host")
	id := flag.String("id","", "client id")
	port := flag.String("port","", "server host")
	changeIPMinutes := flag.Int("changeIPMinutes",3, "change IP minutes")
	if *host == "" {
		log.Fatal("host 必填")
	}
	if *id == "" {
		log.Fatal("id 必填")
	}
	if *port == "" {
		log.Fatal("port 必填")
	}

	for {

		ok := updateIP()
		if !ok {
			continue
			time.Sleep(3 * time.Second)
		}
		//发送请求
		sendUpdate(*host,*id, *port)
		time.Sleep(time.Duration(*changeIPMinutes) * time.Minute)
	}
}

func sendUpdate(host, id, port string) {
	clinet := gorequest.New()
	clinet.Get(fmt.Sprintf("http://%s/add?id=%s&port=%s", host, id, port))
}

func updateIP() bool{
	res, err := execShell("pppoe-stop")
	if err != nil {
		log.Error(err)
	}
	fmt.Println(res)

	res, err = execShell("pppoe-start")
	if err != nil {
		log.Error(err)
	}
	fmt.Println(res)

	res, err = execShell("pppoe-status")
	if err != nil {
		log.Error(err)
	}
	fmt.Println(res)
	//TODO
	return true
}


//阻塞式的执行外部shell命令的函数,等待执行完毕并返回标准输出
func execShell(s string) (string, error){
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

