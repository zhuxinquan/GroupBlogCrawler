package models

import "encoding/xml"

type Csdn struct {
	Rss xml.Name `xml: "rss"`
	Version string `xml:"version,attr"`
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Channel xml.Name `xml:"channel"`
	Title string `xml:"title"`
	Image Image `xml:"iamge"`
	Description string `xml:"description"`
	Link string `xml:"link"`
	Language string `xml:"language"`
	Generator string `xml:"generator"`
	TTL string `xml:"ttl"`
	Copyright string `xml:"copyright"`
	PubDate string `xml:"pubDate"`
	Item []Item `xml:"item"`
}

type Image struct {
	Image xml.Name `xml:"image"`
	Link string `xml:"link"`
	Url string `xml:"url"`
}

type Item struct {
	Item xml.Name `xml:"item"`
	Title string `xml:"title"`
	Link string `xml:"link"`
	Guid string `xml:"guid"`
	Author string `xml:"author"`
	PubDate string `xml:"pubDate"`
	Description string `xml:"description"`
	Category string `xml:"category"`
}