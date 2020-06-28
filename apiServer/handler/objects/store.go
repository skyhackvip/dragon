package objects

import (
	"fmt"
	"io"
	"strings"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"math/rand"
	"github.com/skyhackvip/dragon/service/heartbeat"
	"github.com/skyhackvip/dragon/service/locate"
	"github.com/skyhackvip/dragon/lib/utils"
)

func storeObject(r *http.Request, object string, size int64) (int, error) {
	//找到数据直接返回，不存储
	if locate.Exists(object) {
		return http.StatusOK, nil
	}

	//临时上传
	//r.Body io.Reader

	//随机找到data server
	server := GetRandomServer()
	if server == "" {
		return 0,fmt.Errorf("cannot found dataServer")
	}
	hash := url.PathEscape(object)
	stream, e := newTempPutStream(server, hash, size)
	if e != nil {
		return 0, e
	}

	reader := io.TeeReader(r.Body, stream)
	//hash重复校验(object为hash值）
	d := utils.CalculateHash(reader)
	if d != hash {
		stream.Commit(false) //删除
		return 0, fmt.Errorf("hash dismatch")
	}
	//转正
	stream.Commit(true)

	//store by http stream
	/*stream := newputStream("", object)
	io.Copy(stream, r.Body)
	stream.Close()*/

	return http.StatusOK, nil

}

//调用data server上传
type PutStream struct {
	writer *io.PipeWriter
	c      chan error
}
func newputStream(server, object string) *PutStream {
	c := make(chan error)
	reader, writer := io.Pipe()
	//随机找到data server
	server = GetRandomServer()
	if server == "" {
		c <- fmt.Errorf("cannot found dataServer")
	} else {
		go func() {
			url := fmt.Sprintf("http://%s/datas/%s", server, object)
			fmt.Println("put server:" + url)
			request, _ := http.NewRequest("PUT", url, reader)
			client := http.Client{}
			r, e := client.Do(request)
			if e == nil && r.StatusCode != http.StatusOK {
				e = fmt.Errorf("dataserver return http code %d", r.StatusCode)
			}
			c <- e

		}()
	}
	return &PutStream{writer, c}
}

func (w *PutStream) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

func (w *PutStream) Close() error {
	w.writer.Close()
	return <-w.c
}


/*
func newTempGetStream(server, uuid string) (*GetStream, error) {
	return new
}*/



func GetRandomServer() string {
	dataServers := heartbeat.GetDataServers()
	n := len(dataServers)
	if n == 0 {
		return ""
	}
	return dataServers[rand.Intn(n)]
}
