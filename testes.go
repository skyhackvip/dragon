package main

import (
	"fmt"
	"github.com/skyhackvip/dragon/lib/config"
	"github.com/skyhackvip/dragon/lib/es"
)

//usage go run testes.go -c env.json
func main() {
	config.LoadConfig()
	es.InitEs()

	a, _ := es.GetMetadata("bb", 2)
	fmt.Println(a)

	e := es.PutMetadata("bb", 1, 20, "dsajkldsajkld")
	fmt.Println(e)

	b, _ := es.SearchLatestVersion("bb")
	fmt.Println(b)

	/*b, _ := es.SearchAllVersions("bb", 0, 10)
	fmt.Println(b)

	es.DelMetadata("bb", 3)
	fmt.Println("-----")
	b, _ = es.SearchAllVersions("bb", 0, 10)
	fmt.Println(b)
	*/

}
