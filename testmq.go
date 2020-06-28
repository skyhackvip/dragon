package main

import (
	"fmt"
	"github.com/skyhackvip/dragon/lib"
	"strconv"
	"time"
)

func main() {
	mq := lib.New("amqp://test:test888@10.12.35.8:5672/test-host")
	fmt.Println(mq)
	defer mq.Close()
	//mq.Bind("apiServers")

	/* send
	i, err := lib.ExternalIP()
	if err != nil {
		fmt.Println(err)
	}
	ipStr := i.String()
	go func() {
		mq.Publish("apiServer", "http://"+ipStr+":8878")
		time.Sleep(5 * time.Second)

	}()
	*/

	//get
	go func() {
		mq.Bind("dataServers")
		c := mq.Consume()
		for msg := range c {
			s, e := strconv.Unquote(string(msg.Body))
			fmt.Println(s)
		}
		time.Sleep(1 * time.Second)
	}()

	done := make(chan bool)
	<-done

}
