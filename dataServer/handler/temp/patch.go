package temp
/**
 * uuid.dat 文件上傳對象內容
 */
import (
	"fmt"
	"strings"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"io"
	"os"
)
func patch(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	tempinfo, e := readFromFile(uuid)
	if e != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := "/home/rong/data/temp/" + uuid
	datFile := infoFile + ".dat"

	f, e := os.OpenFile(datFile, os.O_WRONLY|os.O_APPEND, 0) //每次向後追加內容
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, e = io.Copy(f, r.Body) //傳入對象內容
	if e !=nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	info, e := f.Stat()
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("size:" + string(info.Size()) + "/" + string(tempinfo.Size))
	/*if info.Size() > tempinfo.Size {
		os.Remove(datFile)
		os.Remove(infoFile)
		fmt.Println("size not match")
		w.WriteHeader(http.StatusInternalServerError)
	}*/

}

func readFromFile(uuid string)(*tempInfo, error) {
	f,e := os.Open("/home/rong/data/temp/"+ uuid)
	if e!=nil {
		return nil, e
	}
	defer f.Close()
	b, _ := ioutil.ReadAll(f)
	var info tempInfo
	json.Unmarshal(b, &info)
	return &info, nil
}
