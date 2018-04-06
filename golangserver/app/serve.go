package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	as "github.com/aerospike/aerospike-client-go"
)

func main() {
	client, err := as.NewClient("aerospike", 3000)

	if err != nil {
		fmt.Println("error connecting to as!")
		return
	}
	http.HandleFunc("/", serverHandler(client))
	http.ListenAndServe(":8080", nil)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in profile", r)
		}
		client.Close()
	}()
}

func serverHandler(client *as.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		userID := query.Get("userID")
		log.Println("userID: ", userID)

		if userID == "" {
			http.Error(response, "userID is required parameter", http.StatusBadRequest)
			return
		}
		risultato, err := profile(client, 10)
		if err != nil {
			panicOnError(err)
		}
		// ridondanti?
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusOK)
		out, _ := json.MarshalIndent(risultato, "", "  ")
		fmt.Fprintln(response, string(out))
	}
}

func profile(client *as.Client, userID int) (string, error) {
	readPolicy := as.NewPolicy()
	namespace := "cibucks"
	setName := "userProfiles"
	key, _ := as.NewKey(namespace, setName, as.IntegerValue(userID))
	fmt.Println("... accessing key " + key.String())
	rec, err := client.Get(readPolicy, key)
	panicOnError(err)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return fmt.Sprintf(" %v", rec.Bins), err
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
