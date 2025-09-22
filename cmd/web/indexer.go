package main

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rk1165/pse/internal/models"
	"github.com/rk1165/pse/pkg/logger"
)

func Index(request *models.Request, app *application, ch chan<- int) {
	inChannel := make(chan string)
	outChannel := make(chan models.Post)
	var wg sync.WaitGroup
	wg.Add(1)

	links, err := getAllLinks(request.Url, request.Links)
	if err != nil {
		logger.ErrorLog.Printf("error occurred when getting links %v", err)
		ch <- http.StatusInternalServerError
		return
	}

	logger.InfoLog.Printf("Total_Links=%d", len(links))

	workers := max(len(links)/10, 1)
	logger.InfoLog.Printf("Started %d workers", workers)

	for i := 0; i < workers; i++ {
		go createPost(inChannel, outChannel, request.Title, request.Content)
	}

	go savePost(&wg, app, outChannel, len(links))

	for _, link := range links {
		inChannel <- link
	}
	close(inChannel)
	wg.Wait()
	logger.InfoLog.Printf("finished indexing %d posts", len(links))
	logger.InfoLog.Printf("indexed %s", request.Url)
	ch <- http.StatusOK
	close(ch)
}

func getAllLinks(site, selector string) ([]string, error) {
	logger.InfoLog.Printf("fetching links for site %s", site)
	var links []string
	if len(selector) == 0 {
		links = append(links, site)
		return links, nil
	}
	c := colly.NewCollector()
	c.OnHTML(selector, func(e *colly.HTMLElement) {
		url := e.Attr("href")
		absUrl := e.Request.AbsoluteURL(url)
		links = append(links, absUrl)
	})

	err := c.Visit(site)
	if err != nil {
		logger.ErrorLog.Printf("error fetching links for site %s", site)
		return nil, err
	}
	return links, nil
}

func createPost(in <-chan string, out chan<- models.Post, titleSelector, contentSelector string) {
	// only exits if in channel is closed
	for link := range in {
		logger.InfoLog.Printf("processing link=%s", link)
		c := colly.NewCollector()
		var content string
		post := models.Post{Url: link}

		c.OnHTML(titleSelector, func(e *colly.HTMLElement) {
			title := e.Text
			post.Title = clean(title)
		})

		c.OnHTML(contentSelector, func(e *colly.HTMLElement) {
			content = e.Text
			content = clean(content)
			post.Content = content
		})

		err := c.Visit(link)
		if err != nil {
			logger.ErrorLog.Printf("error occurred when visiting link=%s, error=%v", link, err)
			return
		}
		out <- post
	}
}

func savePost(wg *sync.WaitGroup, app *application, out <-chan models.Post, jobs int) {
	defer wg.Done()
	for i := 0; i < jobs; i++ {
		post := <-out
		err := app.post.Insert(post)
		if err != nil {
			logger.ErrorLog.Printf("error occurred when inserting post in db %v", err)
			return
		}
		logger.InfoLog.Printf("post [%v] inserted in db", post.Title)
	}
}

func clean(input string) string {
	words := strings.Fields(input)
	return strings.Join(words, " ")
}
