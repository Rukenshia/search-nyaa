package main

import (
	"fmt"
	"os"
	"os/exec"
	"encoding/xml"
	"log"
	"net/url"
	"regexp"
	"path/filepath"
)

type Rss2 struct {
	XMLName xml.Name `xmlns:"rss"`
	Version string `xml:"version,attr"`
	Title string `xml:"channel>title"`
	Link string `xml:"channel>link"`
	Description string `xml:"channel>description"`
	ItemList []Item `xml:"channel>item"`
}

type Item struct {
	Title string `xml:"title"`
	Link string `xml:"link"`
}

func parseRss2(content []byte) Rss2 {
	rss2 := Rss2{}
	err := xml.Unmarshal(content, &rss2)
	if err != nil {
		log.Panic(err)
	}

	return rss2
}

func main() {
	fmt.Println(os.Args)
	if len(os.Args) < 3 {
		fmt.Println("nope")
		return
	}

	anime := os.Args[1]
	subgroup := os.Args[2]

	quality := ""
	if len(os.Args) >= 4 {
		quality = " " + os.Args[3]
	}

	fmt.Println("kyou best gril")

	// get the file
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Panic("error getting working dir")
	}
	out, err := exec.Command("/usr/bin/bash", filepath.Join(dir, "nyaa-get"), url.QueryEscape(fmt.Sprintf("%s %s%s", anime, subgroup, quality))).CombinedOutput()
	if err != nil {
		log.Panic(err)
	}
	feed := parseRss2(out)


	episodeReg := regexp.MustCompile(fmt.Sprintf(`\[(.*)\] %s - ([0-9]{1,}(v[0-9]{1,})?)`, anime))
	var episodes []string
	for _, item := range feed.ItemList {
		match := episodeReg.FindStringSubmatch(item.Title)
		if len(match) >= 4 {
			episodes = append(episodes, fmt.Sprintf("%s ||| %s ||| %s\n", match[2], item.Link, item.Title))
		}
	}

	// reverse
	for i := len(episodes) - 1; i >= 0; i-- {
		fmt.Print(episodes[i])
	}
	return
}
