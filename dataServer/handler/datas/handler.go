package datas

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

//"github.com/skyhackvip/dragon/lib/utils"
//"net/url"

func Handler(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]

	storage_root := "/data"
	if r.Method == http.MethodPut { //put
		fmt.Println("put request")
		f, err := os.Create(storage_root + "/objects/" + object)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		io.Copy(f, r.Body)
	} else if r.Method == http.MethodGet { //get
		fmt.Println("get request")
		path := storage_root + "/objects/" + object
		fmt.Println(path)
		f, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()

		/*realHash := url.PathEscape(utils.CalculateHash(f))
		fmt.Println(realHash)

		if realHash != object { //object is hash
			fmt.Println("hash mismatch:" + realHash + "/" + object)
			//os.Remove() locate.Del()
			w.WriteHeader(http.StatusNotFound)
		}*/

		io.Copy(w, f)
	}
}
