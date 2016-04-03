package jandan

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strings"
)

func parsePage() ([]Hot, error) {
	hots := make([]Hot, 0)
	doc, err := goquery.NewDocument("http://jandan.net")
	if err != nil {
		fmt.Println(err)
		return hots, err
	}
	parsePicHot(doc, &hots)
	parseOOXXHot(doc, &hots)
	//fmt.Println(hots)

	return hots, nil
}

func parseOOXXHot(doc *goquery.Document, hots *[]Hot) {
	doc.Find("img").Remove()
	doc.Find("div#list-girl div.in").Each(func(i int, s *goquery.Selection) {
		hot := Hot{}
		s.Find("div.acv_author").Each(func(i int, s *goquery.Selection) {
			url, exist := s.Find("a").Attr("href")
			if !exist {
				fmt.Println("Link not exists")
				return
			}
			url = strings.Replace(url, "http://jandan.net", "", -1)
			hot.URL = url
			content, err := s.Next().Html()
			if err != nil {
				fmt.Println(err)
				return
			}
			content = dataCleaning(content)

			hot.Content = content
			hot.Type = OOXX_TYPE
			*hots = append(*hots, hot)
		})
	})

}

func parsePicHot(doc *goquery.Document, hots *[]Hot) {
	doc.Find("img").Remove()
	doc.Find("div#list-pic div.in").Each(func(i int, s *goquery.Selection) {
		hot := Hot{}
		s.Find("div.acv_author").Each(func(i int, s *goquery.Selection) {
			url, exist := s.Find("a").Attr("href")
			if !exist {
				fmt.Println("Link not exists")
				return
			}
			url = strings.Replace(url, "http://jandan.net", "", -1)
			hot.URL = url
			content, err := s.Next().Html()
			if err != nil {
				fmt.Println(err)
				return
			}
			content = dataCleaning(content)

			hot.Content = content
			hot.Type = PIC_TYPE
			*hots = append(*hots, hot)
		})
	})

}

func dataCleaning(content string) string {
	content = strings.TrimSpace(content)
	content = strings.Replace(content, "<p>", "", -1)
	content = strings.Replace(content, "</p>", "", -1)
	content = strings.Replace(content, "<br/>", "\n", -1)
	content = strings.Replace(content, " target=\"_blank\" class=\"view_img_link\"", "", -1)
	return content
}

func GetHot() ([]Hot, error) {
	hots, err := parsePage()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return hots, nil

}

func GetAllComment() ([]Comment, error) {
	comments := []Comment{}
	ooxx, err := getCommentByAPI(OOXX_TYPE)
	pic, err := getCommentByAPI(PIC_TYPE)

	comments = append(comments, ooxx...)
	comments = append(comments, pic...)
	return comments, err
}

func getCommentByAPI(jdType JandanType) ([]Comment, error) {
	var (
		res *http.Response
		err error
	)
	if jdType == OOXX_TYPE {
		res, err = http.Get("http://jandan.net/?oxwlxojflwblxbsapi=jandan.get_ooxx_comments&page=1")
	} else if jdType == PIC_TYPE {
		res, err = http.Get("http://jandan.net/?oxwlxojflwblxbsapi=jandan.get_pic_comments&page=1")
	}

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
	var cr CommentResult
	err = json.Unmarshal(body, &cr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for i, comment := range cr.Comments {
		cr.Comments[i].Type = jdType
		cr.Comments[i].Content = cleanComment(comment.Content)
	}

	return cr.Comments, nil
}

func cleanComment(content string) string {
	content = strings.Replace(content, "<img src", "<a href", -1)
	content = strings.Replace(content, "/>", ">查看原图</a>", -1)
	return content
}

func GetLatestOOXX() *Comment {
	ooxxs, err := getCommentByAPI(OOXX_TYPE)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &ooxxs[0]
}

func GetLatestPic() *Comment {
	pics, err := getCommentByAPI(PIC_TYPE)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &pics[0]
}
