package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var router *gin.Engine

//type pageInfo struct {
//	HTMLVersion             string
//	pageTitle               string
//	Heading1Count           int
//	Heading2Count           int
//	Heading3Count           int
//	Heading4Count           int
//	Heading5Count           int
//	Heading6Count           int
//	AmountInternalLinks     int
//	AmountExternalLinks     int
//	AmountInaccessibleLinks int
//	LoginForm               bool
//}
//
//var linkInfo pageInfo
//
//const (
//	HTML5            = "<!DOCTYPE HTML>"
//	HTML4Strict      = "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\" \n   \"http://www.w3.org/TR/html4/strict.dtd\">"
//	HTML4Traditional = "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01 Transitional//EN\" \n   \"http://www.w3.org/TR/html4/loose.dtd\">"
//	HTML4Frameset    = "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01 Frameset//EN\" \n   \"http://www.w3.org/TR/html4/frameset.dtd\">"
//	XHTMLStrict      = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \n   \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\">"
//	XHTMLTraditional = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \n   \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">"
//	XHTMLFrameset    = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Frameset//EN\" \n   \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-frameset.dtd\">"
//)
//
//func crawl(link, host string, wg *sync.WaitGroup) {
//	time.Sleep(1 * time.Second)
//	defer wg.Done()
//	defer func() {
//		if err := recover(); err != nil {
//			log.Println("panic occurred:", err)
//		}
//	}()
//
//	checkExternal := strings.HasPrefix(link, "http")
//
//	if checkExternal {
//		linkInfo.AmountExternalLinks = linkInfo.AmountExternalLinks + 1
//
//	} else {
//		linkInfo.AmountInternalLinks = linkInfo.AmountInternalLinks + 1
//		link = "https://" + host + link
//	}
//
//	res, err := http.Get(link)
//	if err != nil {
//		linkInfo.AmountInaccessibleLinks = linkInfo.AmountInaccessibleLinks + 1
//
//	}
//	if res.StatusCode != 200 {
//		linkInfo.AmountInaccessibleLinks = linkInfo.AmountInaccessibleLinks + 1
//	}
//
//}
//
//func getHTMLVersion(body string) {
//
//	switch {
//	case strings.Contains(body, HTML5):
//		linkInfo.HTMLVersion = "HTML 5"
//	case strings.Contains(body, HTML4Strict):
//		linkInfo.HTMLVersion = "HTML 4.01 (Strict)"
//	case strings.Contains(body, HTML4Traditional):
//		linkInfo.HTMLVersion = "HTML 4.01 (Transitional)"
//	case strings.Contains(body, HTML4Frameset):
//		linkInfo.HTMLVersion = "HTML 4.01 (Frameset)"
//	case strings.Contains(body, XHTMLStrict):
//		linkInfo.HTMLVersion = "XHTML 1.0 (Strict (quick reference))"
//	case strings.Contains(body, XHTMLTraditional):
//		linkInfo.HTMLVersion = "XHTML 1.0 (Transitional)"
//	case strings.Contains(body, XHTMLFrameset):
//		linkInfo.HTMLVersion = "XHTML 1.0 (Frameset)"
//
//	default:
//		linkInfo.HTMLVersion = "HTML 5"
//	}
//}
//
//func search(c *gin.Context) {
//	pageURL := c.Query("q")
//	res, err := http.Get(pageURL)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer res.Body.Close()
//	if res.StatusCode != 200 {
//		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
//	}
//
//	body := res.Body
//
//	u, err := url.Parse(pageURL)
//	if err != nil {
//		panic(err)
//	}
//	host := u.Host
//
//	//doc, err := goquery.NewDocumentFromReader(body)
//	doc, err := goquery.NewDocument(pageURL)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	html, err := ioutil.ReadAll(body)
//	if err != nil {
//
//		panic(err)
//	}
//
//	getHTMLVersion(string(html[:]))
//
//	linkInfo.pageTitle = doc.Find("title").Contents().Text()
//
//	doc.Find("h1").Each(func(i int, s *goquery.Selection) {
//		linkInfo.Heading1Count = linkInfo.Heading1Count + 1
//	})
//
//	doc.Find("h2").Each(func(i int, s *goquery.Selection) {
//		linkInfo.Heading2Count = linkInfo.Heading2Count + 1
//	})
//
//	doc.Find("h3").Each(func(i int, s *goquery.Selection) {
//		linkInfo.Heading3Count = linkInfo.Heading3Count + 1
//	})
//
//	doc.Find("h4").Each(func(i int, s *goquery.Selection) {
//		linkInfo.Heading4Count = linkInfo.Heading4Count + 1
//	})
//
//	doc.Find("h5").Each(func(i int, s *goquery.Selection) {
//		linkInfo.Heading5Count = linkInfo.Heading5Count + 1
//	})
//
//	doc.Find("h6").Each(func(i int, s *goquery.Selection) {
//		linkInfo.Heading6Count = linkInfo.Heading6Count + 1
//	})
//
//	doc.Find("body input").Each(func(_ int, item *goquery.Selection) {
//		itemId, _ := item.Attr("id")
//		if itemId == "password" {
//			linkInfo.LoginForm = true
//		}
//	})
//
//	login, _ := doc.Find("form").Attr("id")
//	checkLogin := strings.Contains(login, "login")
//	if checkLogin {
//		linkInfo.LoginForm = true
//
//	}
//
//	var wg sync.WaitGroup
//	doc.Find("body a").Each(func(index int, item *goquery.Selection) {
//		wg.Add(1)
//		linkTag := item
//		link, _ := linkTag.Attr("href")
//		go crawl(link, host, &wg)
//	})
//
//	wg.Wait()
//
//	// Call the HTML method of the Context to render a template
//	c.HTML(
//		// Set the HTTP status to 200 (OK)
//		http.StatusOK,
//		// Use the index.html template
//		"index.html",
//		// Pass the data that the page uses
//		gin.H{
//			"title":    "Home Page",
//			"linkInfo": linkInfo,
//		},
//	)
//}

func main() {

	// Set the router as the default one provided by Gin
	router = gin.Default()

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	router.LoadHTMLGlob("templates/*")

	// Define the route for the index page and display the index.html template
	// To start with, we'll use an inline route handler. Later on, we'll create
	// standalone functions that will be used as route handlers.
	router.GET("/", func(c *gin.Context) {

		// Call the HTML method of the Context to render a template
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"index.html",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"title": "Home Page",
			},
		)

	})

	router.GET("/search", Search)

	// Start serving the application
	router.Run()

}
