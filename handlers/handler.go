package handlers

import (
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"news/db"
	"news/db/models"
	"news/repositories"
	"news/useCases"
	"strconv"
)

type MainViewData struct {
	PerPage     int
	Page        int
	PrevPage    int
	NextPage    int
	SearchQuery string
	News        []models.News
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("query")
	rPage := r.URL.Query().Get("page")
	rPerPage := r.URL.Query().Get("perPage")
	page, err := strconv.Atoi(rPage)
	if err != nil {
		page = 1
	}

	perPage, err := strconv.Atoi(rPerPage)
	if err != nil {
		perPage = 15
	}

	newsRepository := repositories.NewNewsRepository(db.DbConnect())
	news, err := newsRepository.GetAllWithPag(page, perPage, searchQuery)

	if err != nil {
		logrus.Error("Something Wrong")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	prevPage := 1
	if page-1 == 0 {
		prevPage = page - 1
	} else {
		prevPage = 1
	}
	data := MainViewData{
		PerPage:     perPage,
		Page:        page,
		PrevPage:    prevPage,
		NextPage:    page + 1,
		SearchQuery: searchQuery,
		News:        news,
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		logrus.Error("Fail parse template")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Fail exec template" + err.Error())
		return
	}
}

type ResourcesListResponse struct {
	Resources []models.NewsResource
}

func ResourcesListHandler(w http.ResponseWriter, r *http.Request) {
	resourcesRepos := repositories.NewResourcesRepository(db.DbConnect())
	resources, err := resourcesRepos.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Error(err.Error())
	}

	tmpl, err := template.ParseFiles("templates/resources.html")
	if err != nil {
		logrus.Error("Fail parse template")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := ResourcesListResponse{Resources: resources}
	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Fail exec template" + err.Error())
		return
	}
}

func ResourcesCreateViewHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/create_resource.html")
	if err != nil {
		logrus.Error("Fail parse template")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Fail exec template" + err.Error())
		return
	}
}

func ResourcesCreateHandler(w http.ResponseWriter, r *http.Request) {
	resourcesRepository := repositories.NewResourcesRepository(db.DbConnect())

	isRss, _ := strconv.ParseBool(r.FormValue("isRss"))
	newResource := models.NewsResource{
		IsRss:            isRss,
		Link:             r.FormValue("link"),
		ItemAttrs:        r.FormValue("itemAttrs"),
		PictureAttrs:     r.FormValue("pictureAttrs"),
		PublishDateAttrs: r.FormValue("publishDateAttrs"),
		TitleAttrs:       r.FormValue("titleAttrs"),
		DescriptionAttrs: r.FormValue("descAttrs"),
		LinkAttrs:        r.FormValue("linkAttrs"),
	}

	err := resourcesRepository.Create(&newResource)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Something Wrong!")
		return
	}

	http.Redirect(w, r, "/resources/list", http.StatusTemporaryRedirect)
}

func DestroyResourceHandler(w http.ResponseWriter, r *http.Request) {
	resourcesRepository := repositories.NewResourcesRepository(db.DbConnect())

	rResourceId := r.FormValue("id")
	if rResourceId == "" {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Undefined Id!")
		return
	}
	resourceId, err := strconv.Atoi(rResourceId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Error("ID must be valid resource id")
		return
	}
	err = resourcesRepository.Destroy(resourceId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Something Wrong!" + err.Error())
		return
	}

	http.Redirect(w, r, "/resources/list", http.StatusTemporaryRedirect)
}

func ActualizeNewsHandler(w http.ResponseWriter, r *http.Request){
	useCases.ActualizeNewsUC()
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
