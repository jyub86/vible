package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

type AppHandler struct {
	http.Handler
}

var rd *render.Render = render.New(render.Options{
	Extensions: []string{".html", ".tmpl"},
})

var initTemplates = template.Must(template.ParseGlob("public/*.html"))

func (a *AppHandler) initHandler(w http.ResponseWriter, r *http.Request) {
	err := initTemplates.ExecuteTemplate(w, "indexPage", nil)
	if err != nil {
		rd.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
}

func (a *AppHandler) searchHandler(w http.ResponseWriter, r *http.Request) {
	data := Item{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
		rd.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	newData, err := data.Search()
	if err != nil {
		log.Println(err)
		rd.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	rd.JSON(w, http.StatusOK, map[string]Item{"data": *newData})
}

func MakeHandler() *AppHandler {
	r := mux.NewRouter()
	n := negroni.Classic()
	n.UseHandler(r)
	a := &AppHandler{
		Handler: n,
	}
	r.HandleFunc("/", a.initHandler)
	r.HandleFunc("/search", a.searchHandler).Methods("POST")
	return a
}
