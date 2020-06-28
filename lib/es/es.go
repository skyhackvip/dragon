package es

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"github.com/skyhackvip/dragon/lib/config"
	"strings"
)

type Metadata struct {
	Name    string
	Version int
	Size    int64
	Hash    string
}

type hit struct {
	Source Metadata `json:"_source"`
}

type searchResult struct {
	Hits struct {
		Total int
		Hits  []hit
	}
}

var server string
var Server string

func InitEs() {
	server = fmt.Sprintf("http://%s:%d/%s/%s",
		config.GlobalEnv.Es.Server,
		config.GlobalEnv.Es.Port,
		config.GlobalEnv.Es.Index,
		config.GlobalEnv.Es.Type)
}

//"10.12.35.8:9200/metadata/objects"

//name_version as id

//get
func GetMetadata(name string, version int) (Metadata, error) {
	var meta Metadata
	url := fmt.Sprintf("%s/%s_%d/_source", server, name, version)
	fmt.Println(fmt.Sprintf("es: %s", url))

	r, e := http.Get(url)
	if e != nil {
		return meta, e
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to get %s_%d :%d", name, version, r.StatusCode)
		return meta, e
	}
	result, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(result, &meta)
	return meta, nil
}

//insert
func InsertMetadata(name string, version int, size int64, hash string) error {
	doc := fmt.Sprintf(`{"name":"%s", "version":%d, "size":%d, "hash":"%s"}`, name, version, size, hash)
	client := http.Client{}
	url := fmt.Sprintf("%s/%s_%d?op_type=create", server, name, version)
	fmt.Println(fmt.Sprintf("es: %s", url))
	request, _ := http.NewRequest("PUT", url, strings.NewReader(doc))
	request.Header.Add("Content-Type", "application/json")
	r, e := client.Do(request)
	if e != nil {
		return e
	}
	if r.StatusCode == http.StatusConflict {
		return InsertMetadata(name, version+1, size, hash)
	}
	if r.StatusCode != http.StatusCreated {
		result, _ := ioutil.ReadAll(r.Body)
		return fmt.Errorf("fail to put metadata:%d %s", r.StatusCode, string(result))
	}
	return nil
}

//insert or update
func PutMetadata(name, hash string, size int64) error {
	metadata, e := SearchLatestVersion(name)
	var version int
	if e != nil || metadata.Hash == "" { //未找到，从v1开始
		version = 1
		return InsertMetadata(name, version, size, hash)
	} else {
		version = metadata.Version
		return InsertMetadata(name, version, size, hash)
	}
}

//get new one
func SearchLatestVersion(name string) (Metadata, error) {
	var meta Metadata
	url := fmt.Sprintf("%s/_search?q=name:%s&size=1&sort=version:desc", server, url.PathEscape(name))
	fmt.Println(fmt.Sprintf("es: %s", url))
	r, e := http.Get(url)
	if e != nil {
		return meta, e
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to search latest metadata:%d", r.StatusCode)
		return meta, e
	}
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		meta = sr.Hits.Hits[0].Source
	}
	return meta, nil
}

//get all
func SearchAllVersions(name string, from, size int) ([]Metadata, error) {
	url := fmt.Sprintf("%s/_search?sort=name,version&from=%d&size=%d", server, from, size)
	fmt.Println(fmt.Sprintf("es: %s", url))
	if name != "" {
		url += "&q=name:" + name
	}
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	metas := make([]Metadata, 0)
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	for i := range sr.Hits.Hits {
		metas = append(metas, sr.Hits.Hits[i].Source)
	}
	return metas, nil

}

//delete
func DelMetadata(name string, version int) {
	client := http.Client{}
	url := fmt.Sprintf("%s/%s_%d", server, name, version)
	fmt.Println(fmt.Sprintf("es: %s", url))
	request, _ := http.NewRequest("DELETE", url, nil)
	client.Do(request)
}
