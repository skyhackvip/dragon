package service

import (
	"fmt"
	"github.com/skyhackvip/dragon/lib/rabbitmq"
	"github.com/skyhackvip/dragon/lib/utils"
	"strconv"
	"sync"
	"time"
)

//data send queue
func StartHeartbeat() {
	fmt.Println("start heartbeat")
	q := rabbitmq.New("amqp://test:test888@10.12.35.8:5672/test-host")
	defer q.Close()
	ip, _ := utils.ExternalIP()
	for {
		q.Publish("dataServers", ip.String()+":8878")
		time.Sleep(5 * time.Second)
	}
}

var dataServers = make(map[string]time.Time)
var mutex sync.Mutex

//api recevice queue
func ListenHeartbeat() {
	fmt.Println("listen heartbeat")
	q := rabbitmq.New("amqp://test:test888@10.12.35.8:5672/test-host")
	q.Bind("dataServers")
	defer q.Close()
	c := q.Consume() //chan
	for msg := range c {
		ip, _ := strconv.Unquote(string(msg.Body))
		mutex.Lock()
		dataServers[ip] = time.Now()
		mutex.Unlock()
	}
}

func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for s, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}

//api
func GetDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string, 0)
	for s, _ := range dataServers {
		ds = append(ds, s)
	}
	return ds
}
