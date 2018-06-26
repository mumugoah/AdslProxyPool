package main

import (
	"math/rand"
	"net/http"
	"time"

	"fmt"

	"github.com/gin-gonic/gin"
)

type Proxy struct {
	Ip   string
	Port string
}

var proxies = map[string]Proxy{}

func main() {
	r := gin.Default()
	r.GET("/add", addProxy)
	r.GET("/del", delProxy)
	r.GET("/get", getProxy)
	r.Run(":9010") // listen and serve on 0.0.0.0:8080
}

type addProxyQuery struct {
	Id   string `form:"id"`
	Port string `form:"port"`
}

//添加proxy
func addProxy(c *gin.Context) {
	var query addProxyQuery
	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.Error(err)
	}
	if query.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "id不能为空"})
	}
	if query.Port == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "port不能为空"})
	}

	//获取IP
	ip := c.ClientIP()
	proxy := Proxy{
		Ip: ip, Port: query.Port,
	}
	//添加
	proxies[query.Id] = proxy

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

//获取proxy
func getProxy(c *gin.Context) {
	var px = []string{}
	for _, p := range proxies {
		px = append(px, fmt.Sprintf("%s:%s", p.Ip, p.Port))
	}
	if len(px) > 0 {
		rand.Seed(time.Now().UnixNano())
		i := rand.Intn(len(px))
		c.String(http.StatusOK, px[i])
	}
	c.String(http.StatusOK, "")
}

type delProxyQuery struct {
	Id   string `form:"id"`
}


//删除proxy
func delProxy(c *gin.Context) {
	var query delProxyQuery
	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.Error(err)
	}
	if query.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "id不能为空"})
	}
	delete(proxies, query.Id)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}