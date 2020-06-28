package locate

import (
	"fmt"
	"os"
	"time"
	"strings"
	"path/filepath"
	"sync"
	"github.com/skyhackvip/dragon/lib/rabbitmq"
	"github.com/skyhackvip/dragon/lib/utils"
	"github.com/skyhackvip/dragon/lib/types"
	"strconv"
)

//data 定位回复消息 locate
func StartLocate() {
	fmt.Println("start locate")
	q := rabbitmq.New("amqp://test:test888@10.12.35.8:5672/test-host")
	defer q.Close()

	//接收消息并回复
	q.Bind("locate")
	c := q.Consume()
	ip, _ := utils.ExternalIP()
	for msg := range c {
		object, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		/*file_path := "log/objects/" + object
		if localExists(file_path) {
			q.Send(msg.ReplyTo, ip.String())
		}*/
		id := localLocate(object)
		if id != -1 {
			q.Send(msg.ReplyTo, types.LocateMessage{Addr:ip.String(), Id:id})
		}
	}
}

var objects = make(map[string]int)
var mutex sync.Mutex
//data
func CollectObjects() {
	files, _ := filepath.Glob("log/objects/*")
	for i := range files {
		file := strings.Split(filepath.Base(file[i]), ".")
		if len(file) != 3 { // xxx.1.xxx
			panic(file[i])
		}
		hash := file[0]
		id, e := strconv.Atoi(file[1])
		if e!= nil {
			panic(e)
		}
		objects[hash] = id
	}
}
func Add(hash string, id int) {
	mutex.Lock()
	objects[hash] = id
	mutex.Unlock()
}
func Del(hash string) {
	mutex.Lock()
	delete(objects, hash)
	mutex.Unlock()
}
//data 查找系统真实数据
func localLocate(name string) int {
	//_, err := os.Stat(name)
	//return !os.IsNotExist(err) io操作过多
	mutex.Lock()
	id, ok := objects[hash]
	mutex.Unlock()
	if !ok {
		return -1
	}
	return id
}

//api 向data发送name(hash值）寻找
func Locate(name string) string {
	q := rabbitmq.New("amqp://test:test888@10.12.35.8:5672/test-host")
	defer q.Close()
	q.Publish("locate", name)

	//等1s接收响应
	c := q.Consume()
	go func(){
		time.Sleep(time.Second)
		q.Close()
	}()
	msg := <- c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}

//api 查找是否存在
func Exists(name string) bool {
	return Locate(name) != ""
}