package temp

import (
	"fmt"
	"net/http"
	"os"
	"encoding/json"
	"strings"
	"strconv"
)
type tempInfo struct {
	Uuid string
	Name string
	Size int64
}
func (t *tempInfo) hash() string {
	s := strings.Split(t.Name, ".")
	return s[0]
}
func (t *tempInfo) id() int {
	s := strings.Split(t.Name, ".")
	id, _ := strconv.Atoi(s[1])
	return id
}
func (t *tempInfo) writeToFile() error {
	f,e := os.Create("log/temp/"+t.Uuid) //生成uuid文件并写入内容
	if e!= nil {
		return e
	}
	defer f.Close()
	b, _ := json.Marshal(t) //uuid, name, size
	f.Write(b)
	return nil
}

func Handler(w http.ResponseWriter, r * http.Request) {
	m := r.Method
	fmt.Println(m + "handler")
	if m == http.MethodPost {
		post(w, r)
		return
	} else if m == http.MethodPut {
		put(w, r)
		return
	} else if m == http.MethodPatch {
		patch(w, r)
		return
	} else if m == http.MethodDelete {
		del(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
