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

	testQuery(client)
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

func testQuery(client *as.Client) error {
	stmt := as.NewStatement("cibucks", "campaigns", "profile")

	rs, err := client.Query(nil, stmt)
	for rec := range rs.Results() {
		if rec.Err != nil {
			log.Println("***** ERROR *****: ", rec.Err)
			// handle error here
			// if you want to exit, cancel the recordset to release the resources
		} else {
			profile := rec.Record.Bins["profile"]
			log.Printf("profile: %v", profile)

			receivedMap := profile.([]interface{})
			for item := range receivedMap {
				log.Printf("item: %v", receivedMap[item])
				internalMap := receivedMap[item].(map[interface{}]interface{})
				log.Printf("interestIds: %v", internalMap["interestIds"])
				log.Printf("groupId: %v", internalMap["groupId"])
			}
		}
	}
	return err
}
