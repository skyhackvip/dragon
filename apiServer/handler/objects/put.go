package objects

import (
	"github.com/skyhackvip/dragon/lib/es"
	"github.com/skyhackvip/dragon/lib/utils"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	fmt.Println("put handler")
	hash := utils.GetHashFromHeader(r.Header)
	size := utils.GetSizeFromHeader(r.Header)
	if hash == "" {
		fmt.Println("hash error")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//curl data url
	name := strings.Split(r.URL.EscapedPath(), "/")[2] //文件名
	//c, e := storeObject(r, name) 存储文件name 改为存储文件hash
	c, e := storeObject(r, url.PathEscape(hash), size)
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
	if c != http.StatusOK {
		fmt.Println("status error")
		w.WriteHeader(c)
		return
	}

	//write es
	es.InitEs()
	e = es.PutMetadata(name, hash, size)
	if e != nil {
		fmt.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(c)
}
