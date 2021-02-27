package aggregator

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"news/db/models"
	"time"
)

type IAggregator interface {
	AggregateFromHtml(res *http.Response, resource models.NewsResource) ([]models.News, error)
	AggregateFromRss(res *http.Response , resource models.NewsResource) ([]models.News, error)
	AggregateNews(resource models.NewsResource) ([]models.News, error)
}

type Aggregator struct {
	Client http.Client
}

type RssResource struct {
	Channel RssChanel `xml:"channel"`
}

type RssChanel struct {
	Items []RssItem `xml:"item"`
}

type RssItem struct {
	Guid        string `xml:"guid"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Category    string `xml:"category"`
	Author      string `xml:"author"`
}

func NewAggregator() IAggregator {
	client := http.Client{Timeout: time.Second * 10}
	var aggregator = Aggregator{
		Client: client,
	}
	return aggregator
}

func (a Aggregator) AggregateNews(resource models.NewsResource) ([]models.News, error) {
	res, err := a.Client.Get(resource.Link)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	if res.StatusCode != 200 {
		err := errors.New(fmt.Sprintf("Failed request, status %d", res.StatusCode))
		logrus.Error(err.Error())
		return nil, err
	}

	var news []models.News
	if resource.IsRss {
		news, err = a.AggregateFromRss(res, resource)
	} else {
		news, err = a.AggregateFromHtml(res, resource)
	}
	return news, err
}

func (a Aggregator) AggregateFromHtml(res *http.Response, resource models.NewsResource) ([]models.News, error) {
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}

	var news []models.News
	doc.Find(resource.ItemAttrs).Each(func(i int, selection *goquery.Selection) {
		newItem := models.News{
			Title:       selection.Find(resource.TitleAttrs).Text(),
			Description: selection.Find(resource.DescriptionAttrs).Text(),
			PublishDate: selection.Find(resource.PublishDateAttrs).Text(),
			Resource: resource,
		}

		newsLinkSelection := selection.Find(resource.LinkAttrs)
		link, exists := newsLinkSelection.Attr("href")
		if !exists {
			return
		}
		newItem.Link = link
		newsPictSelection := selection.Find(resource.PictureAttrs)
		picture, exists := newsPictSelection.Attr("src")
		if exists {
			newItem.Picture = picture
		}

		news = append(news, newItem)
	})

	return news, nil
}

func (a Aggregator) AggregateFromRss(res *http.Response, resource models.NewsResource) ([]models.News, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	rssResource := &RssResource{}
	err = xml.Unmarshal(body, rssResource)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}

	var news []models.News
	for _, item := range rssResource.Channel.Items {
		news = append(news, models.News{
			Title:       item.Title,
			Author:      item.Author,
			Category:    item.Category,
			Description: item.Description,
			PublishDate: item.PubDate,
			Link:        item.Link,
			Resource: resource,
		})
	}

	return news, nil
}
