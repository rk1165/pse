package main

import (
	"github.com/gocolly/colly"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rk1165/pse/internal/models"
	"strings"
	"sync"
)

type Post struct {
	title   string
	url     string
	content string
}

func Index(request *models.Request, app *application) error {
	inChannel := make(chan string)
	outChannel := make(chan Post)
	var wg sync.WaitGroup
	wg.Add(1)

	links, err := getAllLinks(app, request.Url, request.Links)
	if err != nil {
		app.errorLog.Printf("error occurred when getting links %v", err)
		return err
	}

	app.infoLog.Printf("total links %v", len(links))

	workers := max(len(links)/10, 1)
	app.infoLog.Printf("starting %v workers", workers)

	for i := 0; i < workers; i++ {
		go createPost(app, inChannel, outChannel, request.Title, request.Content)
	}

	go savePost(&wg, app, outChannel, len(links))

	for _, link := range links {
		inChannel <- link
	}
	close(inChannel)
	wg.Wait()
	app.infoLog.Printf("finished indexing %v posts", len(links))
	return nil
}

func getAllLinks(app *application, site, selector string) ([]string, error) {
	app.infoLog.Printf("fetching links for site %s", site)
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
		app.errorLog.Printf("error fetching links for site %s", site)
		return nil, err
	}
	return links, nil
}

func createPost(app *application, in <-chan string, out chan<- Post, titleSelector, contentSelector string) {
	// only exits if in channel is closed
	for link := range in {
		app.infoLog.Printf("processing link %s", link)
		c := colly.NewCollector()
		var content string
		post := Post{url: link}

		c.OnHTML(titleSelector, func(e *colly.HTMLElement) {
			title := e.Text
			post.title = clean(title)
		})

		c.OnHTML(contentSelector, func(e *colly.HTMLElement) {
			content = e.Text
			content = clean(content)
			post.content = content
		})

		err := c.Visit(link)
		if err != nil {
			app.errorLog.Printf("error occurred when visiting link %s", link)
			return
		}
		out <- post
	}
}

func savePost(wg *sync.WaitGroup, app *application, out <-chan Post, jobs int) {
	defer wg.Done()
	stmt, err := app.db.Prepare("INSERT INTO posts(title, url, content) values(?,?,?)")
	if err != nil {
		app.errorLog.Printf("error occurred when preparing statement %v", err)
		return
	}
	defer stmt.Close()
	for i := 0; i < jobs; i++ {
		post := <-out
		_, err = stmt.Exec(post.title, post.url, post.content)
		if err != nil {
			app.errorLog.Printf("error occurred when inserting post in db %v", err)
			return
		}
		app.infoLog.Printf("post [%v] inserted in db", post.title)
	}
}

func clean(input string) string {
	words := strings.Fields(input)
	return strings.Join(words, " ")
}
