package handler

import (
	"fmt"
	"github.com/cbotte21/microservice-common/pkg/datastore"
	"github.com/cbotte21/microservice-common/pkg/jwtParser"
	"github.com/cbotte21/username-go/internal/schema"
	"net/http"
)

func SetHandler(w http.ResponseWriter, r *http.Request, usernameClient *datastore.MongoClient[schema.Username], jwtSecret *jwtParser.JwtSecret) {
	err := r.ParseForm() //Populate PostForm
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Please try again later.\n"))
		return
	}

	credentials := r.PostForm

	if !credentials.Has("jwt") || !credentials.Has("username") { //HAS email and password
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request must contain a valid jwt and the desired username.\n"))
		return
	}

	//TODO: Validate JWT, get ID
	res, err := jwtSecret.Redeem(credentials.Get("jwt"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse jwt.\n"))
	}

	_ = usernameClient.Delete(schema.Username{Id: res.Id})
	err = usernameClient.Create(schema.Username{
		Id:       res.Id,
		Username: credentials.Get("username"),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to store username.\n"))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{ "status": true }`)))
}
