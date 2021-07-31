package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/atotto/clipboard"
)

type Item struct {
	Version string
	Book    string
	Chapter string
	Verse   string
	Args    []string
}

func search(version, book, chapter, verse string) string {
	url := "https://www.bskorea.or.kr/bible/korbibReadpage.php?" + fmt.Sprintf("version=%s&book=%s&chap=%s&sec=%s", version, book, chapter, verse)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	str := ""
	doc.Find("span").Each(func(i int, s *goquery.Selection) {
		st, ok := s.Attr("style")
		if ok && strings.Contains(st, "#376BCB") {
			num := s.Find("span.number").Text()
			str = strings.TrimLeft(s.Text(), num)
		}
	})
	return str
}

func (i Item) parser() (*Item, error) {
	re, _ := regexp.Compile(`([ㄱ-ㅎㅏ-ㅣ가-힣]+)([0-9]+):([0-9]+$)`)
	results := re.FindStringSubmatch(strings.Join(i.Args, ""))
	if results == nil {
		return &i, errors.New(errStr)
	}
	name, err := Books(results[1])
	if err != nil {
		return &i, err
	}
	i.Book = name
	i.Chapter = results[2]
	i.Verse = results[3]
	return &i, nil
}

var errStr string = "검색어 에러. ex) 창1:1"

func main() {
	data := Item{}
	if len(os.Args) < 2 {
		log.Fatal(errStr)
	}
	flag.StringVar(&data.Version, "version", "GAE", "개역개정")
	flag.Parse()

	userInput := os.Args[1:]
	for _, arg := range userInput {
		if arg == "-version" {
			continue
		} else if arg == data.Version {
			continue
		} else {
			data.Args = append(data.Args, arg)
		}
	}

	d, err := data.parser()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Search : ", d.Book, d.Chapter, d.Verse, " at ", d.Version)
	out := search(d.Version, d.Book, d.Chapter, d.Verse)
	clipboard.WriteAll(out)
	fmt.Println(out)
}
