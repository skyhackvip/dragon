package objects

import (
	"fmt"
	"io"
	"net/http"
	"github.com/skyhackvip/dragon/lib/rabbitmq"
	"strconv"
	"time"
)

func GetStream(object string) (io.Reader, error) {
	server := Locate(object)
	if server == "" {
		return nil, fmt.Errorf("object %s locate fail", object)
	}
	url := "http://" + server + ":8878/datas/" + object
	fmt.Println("get data from:" + url)
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dataServer return http code %d", r.StatusCode)
	}
	return r.Body, nil //r.Body type is io.Reader

}

//send queue 定位数据
func Locate(name string) string {
	q := rabbitmq.New("amqp://test:test888@10.12.35.8:5672/test-host")
	q.Publish("locate", name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	msg := <-c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}
