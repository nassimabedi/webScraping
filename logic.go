package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type pageInfo struct {
	HTMLVersion             string
	pageTitle               string
	Heading1Count           int
	Heading2Count           int
	Heading3Count           int
	Heading4Count           int
	Heading5Count           int
	Heading6Count           int
	AmountInternalLinks     int
	AmountExternalLinks     int
	AmountInaccessibleLinks int
	LoginForm               bool
	Error                   error
}

var linkInfo pageInfo

const (
	HTML5            = "<!DOCTYPE HTML>"
	HTML4Strict      = "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\" \n   \"http://www.w3.org/TR/html4/strict.dtd\">"
	HTML4Traditional = "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01 Transitional//EN\" \n   \"http://www.w3.org/TR/html4/loose.dtd\">"
	HTML4Frameset    = "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01 Frameset//EN\" \n   \"http://www.w3.org/TR/html4/frameset.dtd\">"
	XHTMLStrict      = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \n   \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\">"
	XHTMLTraditional = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \n   \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">"
	XHTMLFrameset    = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Frameset//EN\" \n   \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-frameset.dtd\">"
)

// crawl godoc
// @Summary crawl links of an HTML source page to find external and internal and inaccessible links
// @Description crawl links of an HTML source page to find external and internal and inaccessible links
func crawl(link, host string, wg *sync.WaitGroup) {
	time.Sleep(1 * time.Second)
	defer wg.Done()
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()

	checkExternal := strings.HasPrefix(link, "http")

	if checkExternal {
		linkInfo.AmountExternalLinks = linkInfo.AmountExternalLinks + 1

	} else {
		linkInfo.AmountInternalLinks = linkInfo.AmountInternalLinks + 1
		link = "https://" + host + link
	}

	res, err := http.Get(link)
	if err != nil {
		linkInfo.AmountInaccessibleLinks = linkInfo.AmountInaccessibleLinks + 1

	}
	if res.StatusCode != 200 {
		linkInfo.AmountInaccessibleLinks = linkInfo.AmountInaccessibleLinks + 1
	}

}

// getHTMLVersion godoc
// @Summary GetHTMLVersion of the source of html
// @Description GetHTMLVersion of the source of html
func getHTMLVersion(body string) string {
	var result string
	switch {
	case strings.Contains(body, HTML5):
		result = "HTML 5"
	case strings.Contains(body, HTML4Strict):
		result = "HTML 4.01 (Strict)"
	case strings.Contains(body, HTML4Traditional):
		result = "HTML 4.01 (Transitional)"
	case strings.Contains(body, HTML4Frameset):
		result = "HTML 4.01 (Frameset)"
	case strings.Contains(body, XHTMLStrict):
		result = "XHTML 1.0 (Strict (quick reference))"
	case strings.Contains(body, XHTMLTraditional):
		result = "XHTML 1.0 (Transitional)"
	case strings.Contains(body, XHTMLFrameset):
		result = "XHTML 1.0 (Frameset)"

	default:
		result = "HTML 5"
	}
	return result
}

// getHeadings godoc
// @Summary get Heading by level source html
// @Description get Heading by level in source html
func getHeadings(doc *goquery.Document) {
	doc.Find("h1").Each(func(i int, s *goquery.Selection) {
		linkInfo.Heading1Count = linkInfo.Heading1Count + 1
	})

	doc.Find("h2").Each(func(i int, s *goquery.Selection) {
		linkInfo.Heading2Count = linkInfo.Heading2Count + 1
	})

	doc.Find("h3").Each(func(i int, s *goquery.Selection) {
		linkInfo.Heading3Count = linkInfo.Heading3Count + 1
	})

	doc.Find("h4").Each(func(i int, s *goquery.Selection) {
		linkInfo.Heading4Count = linkInfo.Heading4Count + 1
	})

	doc.Find("h5").Each(func(i int, s *goquery.Selection) {
		linkInfo.Heading5Count = linkInfo.Heading5Count + 1
	})

	doc.Find("h6").Each(func(i int, s *goquery.Selection) {
		linkInfo.Heading6Count = linkInfo.Heading6Count + 1
	})

}

// hasLoginForm godoc
// @Summary check if the page source has login form
// @Description check if the page source has login form
func hasLoginForm(doc *goquery.Document) bool {
	var loginForm bool
	doc.Find("body input").Each(func(_ int, item *goquery.Selection) {
		itemId, _ := item.Attr("id")
		if itemId == "password" {
			loginForm = true
		}
	})

	login, _ := doc.Find("form").Attr("id")
	checkLogin := strings.Contains(login, "login")
	if checkLogin {
		loginForm = true
	}

	return loginForm

}

// analyse godoc
// @Summary analyse a web page for scraping
// @Description analyse a web page for scraping
func analyse(c *gin.Context) error {
	pageURL := c.Query("q")
	res, err := http.Get(pageURL)
	if err != nil {
		log.Println(err)
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Println("status code error: %d %s", res.StatusCode, res.Status)
		return errors.New("not found : Page can not found")
	}

	body := res.Body

	u, err := url.Parse(pageURL)
	if err != nil {
		log.Println(err)
		return err
	}
	host := u.Host

	doc, err := goquery.NewDocument(pageURL)
	if err != nil {
		log.Println(err)
		return err
	}

	//get HTML version
	html, err := ioutil.ReadAll(body)
	if err != nil {
		log.Println(err)
		return err
	}
	linkInfo.HTMLVersion = getHTMLVersion(string(html[:]))

	// get page title
	linkInfo.pageTitle = doc.Find("title").Contents().Text()

	// get heading
	getHeadings(doc)

	// check if page has Login form
	linkInfo.LoginForm = hasLoginForm(doc)

	// crawl links to find external and internal and inaccessible links
	var wg sync.WaitGroup
	doc.Find("body a").Each(func(index int, item *goquery.Selection) {
		wg.Add(1)
		linkTag := item
		link, _ := linkTag.Attr("href")
		go crawl(link, host, &wg)
	})

	wg.Wait()

	return nil
}

// Search godoc
// @Summary search a link to scrap it
// @Description search a link to scrap it
func Search(c *gin.Context) {
	err := analyse(c)
	if err != nil {
		linkInfo.Error = err
	}
	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"index.html",
		// Pass the data that the page uses
		gin.H{
			"title":    "Home Page",
			"linkInfo": linkInfo,
		},
	)
}
