package objectstream

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//调用dataserver临时上传
type TempPutStream struct {
	Server string
	Uuid string
}
func newTempPutStream(server, object string, size int64) (*TempPutStream,error) {
	url := fmt.Sprintf("http://%s/temp/%s", server, object)
	fmt.Println("put temp server:" + url)
	request, e := http.NewRequest("POST", url, nil)  //post 创建uuid和uuid.dat，并返回uuid
	if e!= nil {
		return nil, e
	}

	request.Header.Set("Size", fmt.Sprintf("%d", size))
	client := http.Client{}
	r, e := client.Do(request)
	if e == nil && r.StatusCode != http.StatusOK {
		e = fmt.Errorf("dataserver return http code %d", r.StatusCode)
	}

	uuid, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return nil , e
	}
	return &TempPutStream{server, string(uuid)} ,nil
}

func (w *TempPutStream) Write(p []byte)(n int, err error) {
	url := fmt.Sprintf("http://%s/temp/%s", w.Server, w.Uuid)
	request, e := http.NewRequest("PATCH", url, strings.NewReader(string(p))) //patch 实际写入内容
	if e != nil {
		return 0, e
	}
	client := http.Client{}
	r, e := client.Do(request)
	if e!= nil {
		return 0, e
	}
	if r.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("dataserver return http code %d", r.StatusCode)
	}
	return len(p), nil
}

func (w *TempPutStream) Commit(good bool) {
	method := "DELETE" //删除
	if good {
		method = "PUT" //转正
	}
	url := fmt.Sprintf("http://%s/temp/%s", w.Server, w.Uuid)
	request, _ := http.NewRequest(method, url, nil)
	client := http.Client{}
	client.Do(request)
}
