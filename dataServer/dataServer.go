/**
 * heart beat send message inclued ip to mq
 * accept put request and save data
 * accept get request and load data
 */
package main

import (
	"fmt"
	"net/http"
	"github.com/skyhackvip/dragon/controller/datas"
	"github.com/skyhackvip/dragon/controller/temp"
	"github.com/skyhackvip/dragon/service/heartbeat"
	"github.com/skyhackvip/dragon/service/locate"
)

func main() {
	fmt.Println("dataserver")
	locate.CollectObjects()
	//heart beat 发送心跳包 bing dataServers
	go heartbeat.StartHeartbeat()
	//locate 回复位置信息 bind locate
	go locate.StartLocate()

	//support get,put
	http.HandleFunc("/datas/", datas.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	http.ListenAndServe(":8878", nil)
}