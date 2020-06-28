package objects
/**
 *  curl -XGET http://localhost:11800/objects/aa
 *  curl -XPOST http://localhost:11800/objects/aa -d "asdjklfsdjkl"
 */

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"strconv"
	"github.com/skyhackvip/dragon/lib/es"
)

func get(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get handler")
	name := strings.Split(r.URL.EscapedPath(), "/")[2]  // server/objects/name
	var e error

	//version
	versionQuery := r.URL.Query()["version"]  // server/objects/name?version=1
	version := 1
	if len(versionQuery) != 0 {
		version, e = strconv.Atoi(versionQuery[0]) // parseint
		if e != nil {
			fmt.Println(e)
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	//根据name找最新版本hash，再请求hash地址
	es.InitEs()
	var meta es.Metadata
	if version > 1 {
		meta, e = es.GetMetadata(name, version)
	} else {
		meta, e = es.SearchLatestVersion(name)
	}
	if e != nil {
		fmt.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if meta.Hash == "" {
		fmt.Println("404")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//根据hash调用data url
	object := url.PathEscape(meta.Hash)
	stream, e := GetStream(object)
	if e != nil {
		fmt.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, stream)
}
