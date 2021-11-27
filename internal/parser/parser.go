package parser

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/opesun/goquery"
	"github.com/opesun/goquery/exp/html"
)

var allowedTags = map[string]struct{}{
	"b":      {},
	"strong": {},
	"i":      {},
	"em":     {},
	"u":      {},
	"ins":    {},
	"s":      {},
	"strike": {},
	"del":    {},
	"a":      {},
}

func Parse() (map[CityCode]*CityEvents, error) {
	url := getArticleUrl()

	if !checkIfPageExists(url) {
		return make(map[CityCode]*CityEvents), nil
	}

	nodes, err := goquery.ParseUrl(url)
	if err != nil {
		return nil, fmt.Errorf("parse: parse url '%s': %s", url, err)
	}

	contentNode := nodes.Find("td.entry-content-right")[0]
	if contentNode == nil {
		return make(map[CityCode]*CityEvents), nil
	}

	contentNodes := contentNode.Child

	result := make(map[CityCode]*CityEvents, 0)
	for _, city := range knownCities {
		events := getCityEvents(contentNodes, city)
		if len(events) > 0 {
			result[city.Code] = &CityEvents{
				City:   city,
				Events: events,
			}
		}
	}

	return result, err
}

func getCityEvents(nodes []*html.Node, city *City) []*Event {
	result := make([]*Event, 0)
	start, end := getCityIndex(nodes, city.Name)
	if start == 0 {
		return result
	}

	for {
		event, nextEventIndex := getEvent(start, end, nodes)
		start = nextEventIndex
		if event.Title != "" && event.DescriptionHTML != "" {
			result = append(result, event)
		}

		if nextEventIndex >= end {
			break
		}
	}

	return result
}

func getCityIndex(nodes []*html.Node, city string) (int, int) {
	var (
		start int
		end   int
	)

	for i, node := range nodes {
		if node.Data == "h1" && getNodeText(node) == city {
			start = i
			continue
		}

		if start > 0 && (node.Data == "h1" || i == len(nodes)-1) {
			end = i
			continue
		}

		if start > 0 && end > 0 {
			break
		}
	}

	return start, end
}

func getEvent(start int, end int, nodes []*html.Node) (*Event, int) {
	var (
		title           string
		descriptionHTML string
		imageURL        string
		nextEvent       int
	)

	for i := start; i <= end; i++ {
		node := nodes[i]

		if node.Data == "h2" && title != "" {
			nextEvent = i
			break
		}

		if node.Data == "p" && descriptionHTML != "" {
			nextEvent = end
			break
		}

		if node.Data == "h2" {
			title = getNodeText(node)
		}

		if node.Data == "p" {
			descriptionHTML = cleanupAndRenderMarkup(node)
		}

		if val, _ := getAttribute(node, "class"); strings.Contains(val, "wp-block-image") {
			img := searchImageNode(node)
			if img != nil {
				imageURL, _ = getAttribute(img, "src")
			}
		}
	}

	if nextEvent == 0 {
		nextEvent = end
	}

	return &Event{
		Title:           strings.TrimSpace(title),
		DescriptionHTML: descriptionHTML,
		ImageURL:        imageURL,
	}, nextEvent
}

func renderNode(node *html.Node) string {
	buf := bytes.NewBufferString("")
	if err := html.Render(buf, node); err != nil {
		log.Println(err)
		return ""
	}

	return buf.String()
}

func cleanupAndRenderMarkup(root *html.Node) string {
	for i, node := range root.Child {
		if containsNotAllowedTag(node) {
			content := renderNode(node)
			root.Child[i] = &html.Node{
				Parent: root,
				Type:   html.TextNode,
				Data:   content,
			}
		}
	}

	result := renderNode(root)
	result = strings.ReplaceAll(result, "<p>", "")
	result = strings.ReplaceAll(result, "</p>", "")

	return result
}

func containsNotAllowedTag(node *html.Node) bool {
	if _, ok := allowedTags[node.Data]; node.Type != html.TextNode && !ok {
		return false
	}

	for _, n := range node.Child {
		return containsNotAllowedTag(n)
	}

	return true
}

func searchImageNode(node *html.Node) *html.Node {
	if node.Data == "img" {
		return node
	}

	for _, child := range node.Child {
		return searchImageNode(child)
	}

	return nil
}

func getAttribute(node *html.Node, key string) (string, bool) {
	for _, attribute := range node.Attr {
		if attribute.Key == key {
			return attribute.Val, true
		}
	}

	return "", false
}

func getArticleUrl() string {
	return "https://sadwave.com/2021/11/16/"
	year, month, day := time.Now().Date()
	return fmt.Sprintf("%s/%d/%d/%d/", sadwaveURL, year, int(month), day)
}

func getNodeText(node *html.Node) string {
	text := &bytes.Buffer{}
	collectText(node, text)
	return text.String()
}

func collectText(node *html.Node, buf *bytes.Buffer) {
	if node.Type == html.TextNode {
		buf.WriteString(node.Data)
	}

	for _, n := range node.Child {
		collectText(n, buf)
	}
}

func checkIfPageExists(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return false
	}

	if resp.StatusCode > 199 && resp.StatusCode < 300 {
		return true
	}

	if resp.StatusCode == 404 {
		return false
	}

	log.Println(fmt.Sprintf("uexpected status '%s'", url))
	return false
}
