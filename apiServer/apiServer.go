package main

/**
 *  echo -n "aaabbbccc"|openssl dgst -sha256 -binary | base64
 *  put curl -v -XPUT http://localhost:11800/objects/aa -d "aaabbbccc" -H "Digest: SHA-256=+4SkX2330dFwNvk58c/rhzOf9dvfQRIi83Yt12d5ooc="
 *  get one curl -XGET http://localhost:11800/objects/aa
 *  get all curl -XGET http://localhost:11800/versions/aa
 */
import (
	"github.com/skyhackvip/dragon/apiServer/handler/objects"
	"github.com/skyhackvip/dragon/apiServer/handler/versions"
	"github.com/skyhackvip/dragon/apiServer/service/heartbeat"
	"github.com/skyhackvip/dragon/lib/config"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("apiserver start...")
	//接收心跳包 bind dataServers
	go heartbeat.ListenHeartbeat()

	config.LoadConfig()

	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	listen_address := config.GlobalEnv.Server.ListenAddress //":11800"
	log.Fatal(http.ListenAndServe(listen_address, nil))
}
