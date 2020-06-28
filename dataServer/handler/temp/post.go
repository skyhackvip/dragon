package temp
/**
 * 生成uuid文件，並存入信息
 * 生成uuid.dat文件（空文件）等待patch方法上傳內容
 */
import(
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"strconv"
)

func post(w http.ResponseWriter, r *http.Request) {
	output, _ := exec.Command("uuidgen").Output()

	uuid := strings.TrimSuffix(string(output), "\n")
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, e := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if e != nil {
		fmt.Println(e)
		return
	}
	t := tempInfo{uuid, name, size}
	fmt.Println(t)
	e = t.writeToFile()
	if e != nil {
		fmt.Println(e)
		return
	}
	os.Create("log/temp/"+t.Uuid+".dat") //生成uuid.dat
	fmt.Println("ok")

	w.Write([]byte(uuid))
}
