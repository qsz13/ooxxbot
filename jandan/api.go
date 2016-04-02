package jandan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func downloadHomePage() (string, error) {
	res, err := http.Get("http://jandan.net")
	if err != nil {
		fmt.Println(err)
		return "", err

	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	html := string(body)
	return html, nil
}

func parsePage(html string) (map[string][]string, error) {
	result := make(map[string][]string)

	return result, nil
}

func getPic() {

}

func getOOXX() {

}

func getOOXXByAPI() ([]OOXX, error) {
	res, err := http.Get("http://jandan.net/?oxwlxojflwblxbsapi=jandan.get_ooxx_comments&page=1")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var or OOXXResult
	err = json.Unmarshal(body, &or)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return or.Comments, nil
}

func getPicByAPI() ([]Pic, error) {
	res, err := http.Get("http://jandan.net/?oxwlxojflwblxbsapi=jandan.get_pic_comments&page=1")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var pr PicResult
	err = json.Unmarshal(body, &pr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return pr.Comments, nil
}

func GetLatestOOXX() *OOXX {
	ooxxs, err := getOOXXByAPI()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &ooxxs[0]
}

func GetLatestPic() *Pic {
	pics, err := getPicByAPI()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &pics[0]
}
