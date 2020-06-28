package temp
/**
 * 根據數據校驗結果，決定臨時文件轉正或刪除
 */
import (
	"net/http"
	"strings"
	"os"
	"fmt"
	"github.com/skyhackvip/dragon/service/locate"
)

func put(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	tempinfo, e := readFromFile(uuid)
	if e!=nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := "log/temp/" + uuid
	datFile := infoFile + ".dat"
	f, e := os.Open(datFile)
	if e!=nil {
		return
	}
	defer f.Close()
	info, e := f.Stat()
	if e != nil {
		return
	}
	fmt.Println(info.Size())
	os.Remove(infoFile)
	/*if info.Size() != tempinfo.Size {
		os.Remove(datFile)
		return
	}*/
	commitTempObject(datFile, tempinfo)
}

func commitTempObject(datFile string, tempinfo *tempInfo) {
	fmt.Println("rename")
	//os.Rename(datFile, "log/objects/"+tempinfo.Name)
	//locate.Add(tempinfo.Name) //加入定位内容
	f, _ := os.Open(datFile)
	d := url.PathEscape(utils.CaculateHash(f))
	f.Close()
	os.Rename(datFile, "log/objects/"+tempinfo.Name + "."+d)
	locate.Add(tempinfo.hash(), tempinfo.id())
}
