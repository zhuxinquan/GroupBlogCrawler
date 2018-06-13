package models

import "encoding/xml"

type Wordpress struct {
	Rss     xml.Name         `xml: "rss"`
	Channel WordpressChannel `xml:"channel"`
}

type WordpressChannel struct {
	Channel xml.Name        `xml:"channel"`
	Item    []WordpressItem `xml:"item"`
}

type WordpressItem struct {
	Item        xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Author      string   `xml:"creator"`
	PubDate     string   `xml:"pubDate"`
	Summary     string   `xml:"description"`
	Category    []string `xml:"category"`
	Description string   `xml:"encoded"`
}
