package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	as "github.com/aerospike/aerospike-client-go"
)

func main() {
	client, err := as.NewClient("aerospike", 3000)

	if err != nil {
		fmt.Println("error connecting to as!")
		return
	}
	fmt.Println("listening...")
	http.HandleFunc("/", serverHandler(client))
	http.ListenAndServe(":8080", nil)

}

func serverHandler(client *as.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				http.Error(response, "userID not found or else", http.StatusBadRequest)
			}
		}()

		fmt.Println("serving: ")
		query := request.URL.Query()
		userID, _ := strconv.Atoi(query.Get("userID"))
		log.Println("userID: ", userID)

		if userID == 0 {
			http.Error(response, "userID is required parameter", http.StatusBadRequest)
			return
		}
		risultato, _ := profile(client, userID)
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
	rec, err := client.Get(readPolicy, key)

	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", rec.Bins), err
}
