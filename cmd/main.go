package main

import (
	"github.com/cbotte21/microservice-common/pkg/datastore"
	"github.com/cbotte21/microservice-common/pkg/environment"
	"github.com/cbotte21/microservice-common/pkg/jwtParser"
	"github.com/cbotte21/username-go/internal"
	"github.com/cbotte21/username-go/internal/schema"
	"log"
	"os"
	"strconv"
)

func main() {
	// Verify environment variables exist
	environment.VerifyEnvVariable("port")
	environment.VerifyEnvVariable("mongo_uri")
	environment.VerifyEnvVariable("mongo_db")
	environment.VerifyEnvVariable("jwt_secret")

	port, err := strconv.Atoi(environment.GetEnvVariable("port"))
	if err != nil {
		log.Fatalf("Could not parse port variable!")
	}

	usernameClient := datastore.MongoClient[schema.Username]{}
	err = usernameClient.Init()
	if err != nil {
		panic(err)
	}

	jwtSecret := jwtParser.JwtSecret(os.Getenv("jwt_secret"))

	api, res := service.NewApi(port, &usernameClient, &jwtSecret)
	if !res || api.Start() != nil { //Start API Listener
		log.Fatal("Failed to initialize API.")
	}
}
