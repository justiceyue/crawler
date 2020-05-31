package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	body, err := fetch("https://movie.douban.com/chart")
	if err != nil {
		fmt.Println(err)
	}

	result := parseBody(body)
	for _, item := range result {
		fmt.Println(strings.Trim(item, " <"))
	}
}

func fetch(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.3")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New("failed to get response:" + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func parseBody(body string) []string {
	body = strings.Replace(body, "\n", "", -1)
	var result = make([]string, 0)

	classArticle := regexp.MustCompile(`<div class="article">(.*)</div>`)
	titleRe := regexp.MustCompile(`<h2>(.*?)</h2>`)
	filmNameRe := regexp.MustCompile(`<a href="https://movie.douban.com/subject/(\d+)/"  class="">(.*?)/`)

	items := classArticle.FindAllStringSubmatch(body, -1)
	for _, item := range items {
		result = append(result, titleRe.FindStringSubmatch(item[1])[1])
		fileNames := filmNameRe.FindAllStringSubmatch(item[1], -1)
		for _, fileName := range fileNames {
			result = append(result, fileName[2])
		}
	}

	return result
}
