// File: ./chat.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	pusher "github.com/pusher/pusher-http-go"
)

// use os package to get the env variable which is already set
func envVariable(key string) string {

	// set env variable using os package
	os.Setenv(key, "gopher")

	// return the env variable using os package
	return os.Getenv(key)
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

var client = pusher.Client{
	AppID:   "PUSHER_APP_ID",
	Key:     "PUSHER_APP_KEY",
	Secret:  "PUSHER_APP_SECRET",
	Cluster: "PUSHER_APP_CLUSTER",
	Secure:  true,
}

type user struct {
	Name  string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email"`
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./public")))

	http.HandleFunc("/new/user", registerNewUser)
	http.HandleFunc("/pusher/auth", pusherAuth)

	log.Fatal(http.ListenAndServe(":8090", nil))
}

func registerNewUser(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	var newUser user

	err = json.Unmarshal(body, &newUser)
	if err != nil {
		panic(err)
	}

	client.Trigger("update", "new-user", newUser)

	json.NewEncoder(rw).Encode(newUser)
}

func pusherAuth(res http.ResponseWriter, req *http.Request) {
	params, _ := ioutil.ReadAll(req.Body)
	response, err := client.AuthenticatePrivateChannel(params)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(res, string(response))
}
