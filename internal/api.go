package service

import (
	"github.com/cbotte21/microservice-common/pkg/datastore"
	"github.com/cbotte21/microservice-common/pkg/jwtParser"
	"github.com/cbotte21/username-go/internal/handler"
	"github.com/cbotte21/username-go/internal/schema"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Api struct {
	port           int
	router         *mux.Router
	jwtSecret      *jwtParser.JwtSecret
	usernameClient *datastore.MongoClient[schema.Username]
}

func NewApi(port int, usernameClient *datastore.MongoClient[schema.Username], jwtSecret *jwtParser.JwtSecret) (*Api, bool) {
	api := &Api{}
	api.port = port
	api.usernameClient = usernameClient
	api.jwtSecret = jwtSecret
	api.router = mux.NewRouter()
	api.RegisterHandlers()
	return api, true
}

func (api *Api) Start() error { //maybe change return to bool
	return http.ListenAndServe(":"+strconv.Itoa(api.port), api.router)
}

func (api *Api) RegisterHandlers() { //Add all API handlers here
	prefix := "/api"

	api.router.HandleFunc(prefix+"/", func(w http.ResponseWriter, r *http.Request) {
		handler.GetHandler(w, r, api.usernameClient)
	}).Methods("GET")

	api.router.HandleFunc(prefix+"/", func(w http.ResponseWriter, r *http.Request) {
		handler.SetHandler(w, r, api.usernameClient, api.jwtSecret)
	}).Methods("POST")
}
