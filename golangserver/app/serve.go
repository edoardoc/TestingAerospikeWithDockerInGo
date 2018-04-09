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
	http.HandleFunc("/", serverHandler(client))
	http.HandleFunc("/validCampaigns", validCampaigns(client))
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
			http.Error(response, "userID is a required parameter", http.StatusBadRequest)
			return
		}
		result, _ := userProfiles(client, userID)
		response.WriteHeader(http.StatusOK)
		fmt.Fprintln(response, string(fmt.Sprintf("%v", result["profile"])))
	}
}

func userProfiles(client *as.Client, userID int) (as.BinMap, error) {
	readPolicy := as.NewPolicy()
	namespace := "cibucks"
	setName := "userProfiles"
	key, _ := as.NewKey(namespace, setName, as.IntegerValue(userID))
	rec, err := client.Get(readPolicy, key)

	if err != nil {
		return nil, err
	}
	return rec.Bins, err
}

func validCampaigns(client *as.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				http.Error(response, "userID not found or else", http.StatusBadRequest)
			}
		}()

		fmt.Println("serving validCampaigns: ")
		query := request.URL.Query()
		userID, _ := strconv.Atoi(query.Get("userID"))
		log.Println("userID: ", userID)

		if userID == 0 {
			http.Error(response, "userID is a required parameter", http.StatusBadRequest)
			return
		}
		result, _ := userProfiles(client, userID)

		stmt := as.NewStatement("cibucks", "campaigns", "profile")
		rs, _ := client.Query(nil, stmt)
		for rec := range rs.Results() {
			if rec.Err != nil {
				log.Println("***** ERROR *****: ", rec.Err)
			} else {
				// User should have all the groupIds that the campaign targets
				// User should have at least one same interestId per groupId
				campaigns := rec.Record.Bins["profile"].([]interface{})
				for n := range campaigns {
					campaign := campaigns[n].(map[interface{}]interface{})
					log.Printf("interestIds: %v", campaign["interestIds"])
					log.Printf("groupId: %v", campaign["groupId"])
					for k := range result["profile"].([]interface{}) {
						userProfile := result["profile"].([]interface{})[k].(map[interface{}]interface{})
						log.Printf("userProfile interestIds: %v", userProfile["interestIds"])
						log.Printf("userProfile groupId: %v", userProfile["groupId"])
					}
				}

			}
		}

		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusOK)
		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Fprintln(response, string(out))
	}
}

func testQuery(client *as.Client) (string, error) {
	stmt := as.NewStatement("cibucks", "campaigns", "profile")

	rs, err := client.Query(nil, stmt)
	for rec := range rs.Results() {
		if rec.Err != nil {
			log.Println("***** ERROR *****: ", rec.Err)
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
	return "", err
}
