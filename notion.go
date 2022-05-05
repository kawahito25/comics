package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const NOTION_VERSION = "2022-02-22"

func getComicTitles() []string {

	url := "https://api.notion.com/v1/databases/" + os.Getenv("TITLES_DATABASE_ID") + "/query"

	payload := strings.NewReader("{\"page_size\":100}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Notion-Version", NOTION_VERSION)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + os.Getenv("NOTION_API_KEY"))

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var resBody interface{}

	err := json.Unmarshal(body, &resBody)

	if err != nil {
		fmt.Println("error:", err)
	}

	var titles []string

	results := resBody.(map[string]interface{})["results"].([]interface{})

	for _, result := range results {
		propperties := result.(map[string]interface{})["properties"]
		title := propperties.(map[string]interface{})["名前"].(map[string]interface{})["title"].([]interface{})
		plainText := title[0].(map[string]interface{})["plain_text"]
		titles = append(titles, plainText.(string))
	}

	return titles
}


func insert(comic *comicInfo) {
	url := "https://api.notion.com/v1/pages"

	jsonBody, err := ioutil.ReadFile("page_template.json")
	if err != nil {
		fmt.Println("error:", err)
	}

	jsonPayload := string(jsonBody)
	jsonPayload = strings.Replace(jsonPayload, "<<database_id>>", os.Getenv("COMICS_TO_READ_DATABASE_ID"), 1)
	jsonPayload = strings.Replace(jsonPayload, "<<comic_title>>", comic.title, 1)
	jsonPayload = strings.Replace(jsonPayload, "<<shelf_num>>", comic.shelfNumber, 1)
	jsonPayload = strings.Replace(jsonPayload, "<<arrival_date>>", comic.arrivalDate, 1)

	payload := strings.NewReader(jsonPayload)
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Notion-Version", NOTION_VERSION)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + os.Getenv("NOTION_API_KEY"))

	_, ok := http.DefaultClient.Do(req)
	if ok != nil {
		fmt.Println("error:", ok)
	}

}


func insertComics(comics comics) {
	for _, c := range comics {
		insert(c)
	}
}


