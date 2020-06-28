package versions
import (
	"fmt"
	"net/http"
	"encoding/json"
	"strings"
	"github.com/skyhackvip/dragon/lib/es"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("versions")
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	from := 0
	size := 100
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	es.InitEs()
	for  {
		metas, e := es.SearchAllVersions(name, from, size)
		if e != nil {
			fmt.Println(e)
			return
		}
		for i := range metas {
			b, _ := json.Marshal(metas[i])
			w.Write(b)
			w.Write([]byte("\n"))
		}
		if len(metas) != size {
			return
		}
		from += size
	}
}

