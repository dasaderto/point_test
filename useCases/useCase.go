package useCases

import (
	"github.com/sirupsen/logrus"
	"news/aggregator"
	"news/db"
	"news/db/models"
	"news/repositories"
	"sync"
)

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func ActualizeNewsUC(){
	connect := db.DbConnect()
	resourcesRepos := repositories.NewResourcesRepository(connect)
	newsRepos := repositories.NewNewsRepository(connect)

	resources, err := resourcesRepos.GetAll()
	if err != nil{
		logrus.Error(err.Error())
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(resources))

	workChan := make(chan []models.News, len(resources))
	for _, resource := range resources{
		go FetchResourceNews(resource, workChan, &wg)
	}
	wg.Wait()
	close(workChan)
	for news := range workChan{
		if len(news) == 0{
			continue
		}
		//Получение ссылок для поиска уникальных
		var links []string
		for _, item := range news{
			links = append(links, item.Link)
		}
		//Конвертация найденных в список ссылок
		existsNews, err := newsRepos.GetByLinks(links)
		var existsNewsLinks []string
		for _, item := range existsNews{
			existsNewsLinks = append(existsNewsLinks, item.Link)
		}

		var newsToCreate []models.News
		for _, item := range news{
			_, isFound := Find(existsNewsLinks, item.Link)
			if isFound{
				continue
			}
			newsToCreate = append(newsToCreate, item)
		}
		if len(newsToCreate) == 0{
			continue
		}
		err = newsRepos.Create(&newsToCreate)
		if err!=nil{
			logrus.Error("Fail news pack create")
			continue
		}
	}
}

func FetchResourceNews(resource models.NewsResource,workChan chan []models.News, wg *sync.WaitGroup){
	defer wg.Done()
	newsAggregator := aggregator.NewAggregator()
	news, err := newsAggregator.AggregateNews(resource)
	if err != nil{
		workChan <- []models.News{}
		return
	}

	workChan <- news
}