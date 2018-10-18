package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"syscall"

	simplejson "github.com/bitly/go-simplejson"
	"golang.org/x/crypto/ssh/terminal"
)

type Pattern struct {
	Id     string
	Title  string
	Schema []byte
}

var (
	url          string
	username     string
	password     string
	version      string
	indexPattern string
	check        bool
)

func init() {
	flag.StringVar(&url, "uri", "http://localhost:5601", "Url to access Kibana. e.g. http://localhost:5601")
	flag.StringVar(&username, "u", "", "Username (if authentication is required)")
	flag.StringVar(&version, "v", "6.4.0", "Kibana version")
	flag.StringVar(&indexPattern, "i", "", "Specify index pattern as regular expression")
	flag.BoolVar(&check, "c", false, "Whether to check the index-pattern matching the regular expression")
}

func main() {
	flag.Parse()
	if indexPattern == "" {
		log.Fatal("indexPattern is not defined.")
	}
	if username != "" {
		fmt.Print("Password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print("\n")
		password = string(bytePassword)
	}

	// get infomation of index-pattern
	patterns, err := getAllIndexPattern(indexPattern)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Println("# Update Targets")
	for _, p := range patterns {
		fmt.Println("Title: " + p.Title + ", Id: " + p.Id)
	}

	if check {
		os.Exit(0)
	}

	fmt.Println("# Update Schema")
	for _, p := range patterns {
		schema, err := getSchema(p.Title)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		p.Schema = schema

		// set information of index-pattern
		_, err = setSchema(p)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		} else {
			fmt.Println(p.Title + ": Update Success")
		}
	}
}

func getAllIndexPattern(indexPattern string) (patterns []Pattern, err error) {
	url := url + "/api/saved_objects/_find?type=index-pattern&fields=title&per_page=10000"
	result, err := doGet(url)
	if err != nil {
		log.Fatal(err)
		return
	}
	jsonData, err := simplejson.NewJson(result)
	if err != nil {
		log.Fatal(err)
		return
	}
	savedObjects := jsonData.Get("saved_objects").MustArray()
	r := regexp.MustCompile(indexPattern)
	for _, value := range savedObjects {
		id := value.(map[string]interface{})["id"].(string)
		attributes := value.(map[string]interface{})["attributes"]
		title := attributes.(map[string]interface{})["title"].(string)
		if r.MatchString(title) {
			patterns = append(patterns, Pattern{Id: id, Title: title})
		}
	}
	return
}

func getSchema(indexPattern string) (schema []byte, err error) {
	getUrl := url +
		"/api/index_patterns/_fields_for_wildcard?pattern=" +
		indexPattern +
		"&meta_fields=%5B%22_source%22%2C%22_id%22%2C%22_type%22%2C%22_index%22%2C%22_score%22%5D"
	result, err := doGet(getUrl)
	schema = result
	return
}

func setSchema(pattern Pattern) (result []byte, err error) {
	setUrl := url +
		"/api/saved_objects/index-pattern/" +
		pattern.Id
	schemaJsonData, err := simplejson.NewJson(pattern.Schema)
	if err != nil {
		log.Fatal(err)
		return
	}
	jsonFields, err := json.Marshal(schemaJsonData.Get("fields"))
	if err != nil {
		log.Fatal(err)
		return
	}
	queryJsonData := simplejson.New()
	queryJsonData.SetPath([]string{"attributes", "title"}, pattern.Title)
	queryJsonData.SetPath([]string{"attributes", "fields"}, string(jsonFields))
	body, _ := queryJsonData.Encode()
	req, err := http.NewRequest("PUT", setUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("kbn-version", version)
	if err != nil {
		log.Fatal(err)
		return
	}
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	result = byteArray
	return
}

func doGet(url string) (result []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	result = byteArray
	return
}
