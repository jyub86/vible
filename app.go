package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Item struct {
	Word     string `json:"word"`
	Version  string `json:"version"`
	Version2 string `json:"version2"`
	Clip     bool   `json:"clip"`
	Book     string `json:"book"`
	Chapter  string `json:"chapter"`
	Verse    string `json:"verse"`
	Result   string `json:"result"`
	Result2  string `json:"result2"`
	Title    string `json:"title"`
	Title2   string `json:"title2"`
	Option   string `json:"option"`
}

func (i Item) Search() (*Item, error) {
	re, _ := regexp.Compile(`(^[ㄱ-ㅎㅏ-ㅣ가-힣]+)([0-9]+):([0-9]+$)`)
	results := re.FindStringSubmatch(strings.Replace(i.Word, " ", "", -1))
	if results == nil {
		re, _ = regexp.Compile(`(^.+[ㄱ-ㅎㅏ-ㅣ가-힣])([0-9]+):([0-9]+$)`)
		results = re.FindStringSubmatch(strings.Replace(i.Word, " ", "", -1))
		if results == nil {
			return &i, errors.New(fmt.Sprintf("%s, 검색어 에러. ex) 창 1:1", i.Word))
		}
	}
	book, err := Books(results[1])
	if err != nil {
		return &i, err
	}
	i.Book = book
	i.Chapter = results[2]
	i.Verse = results[3]
	// option
	switch i.Option {
	case "prev":
		verse, _ := strconv.Atoi(i.Verse)
		if verse != 1 {
			verse = verse - 1
		}
		i.Verse = strconv.Itoa(verse)
	case "next":
		verse, _ := strconv.Atoi(i.Verse)
		i.Verse = strconv.Itoa(verse + 1)
	}
	// search
	title, result := crawling(i.Version, i.Book, i.Chapter, i.Verse)
	if result == "" {
		return &i, errors.New(fmt.Sprintf("%s, %s, %s, %s, search failed", i.Version, i.Book, i.Chapter, i.Verse))
	}
	i.Title, i.Result = title, result
	if i.Version2 != "" {
		title, result = crawling(i.Version2, i.Book, i.Chapter, i.Verse)
		if result == "" {
			return &i, errors.New(fmt.Sprintf("%s, %s, %s, %s, search failed.", i.Version2, i.Book, i.Chapter, i.Verse))
		}
		i.Title2, i.Result2 = title, result
	}
	return &i, nil
}

func crawling(version, book, chapter, verse string) (title, result string) {
	url := "http://bible.godpia.com/read/reading.asp?" + fmt.Sprintf("ver=%s&ver2=&vol=%s&chap=%s&sec=%s", version, book, chapter, verse)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	keyword := fmt.Sprintf("#%s_%s_%s_%s", version, book, chapter, verse)
	doc.Find(keyword).Each(func(_ int, s *goquery.Selection) {
		result = s.Text()
	})
	re, _ := regexp.Compile(`^([0-9]+)(.+)$`)
	results := re.FindStringSubmatch(result)
	if results == nil {
		return title, result
	}
	result = results[2]
	titleStr := doc.Find("#selectBible1").Text()
	title = fmt.Sprintf("%s %s:%s", titleStr, chapter, verse)
	return title, result
}
