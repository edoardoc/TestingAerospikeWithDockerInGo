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
		record, _ := userProfiles(client, userID)
		binsThisUser := record["profile"].([]interface{})

		output := []int{}
		stmt := as.NewStatement("cibucks", "campaigns", "key", "profile")
		rsCampaigns, _ := client.Query(nil, stmt)
		for recCampaign := range rsCampaigns.Results() {
			if recCampaign.Err != nil {
				log.Println("***** ERROR *****: ", recCampaign.Err)
			} else {
				binsCampaigns := recCampaign.Record.Bins["profile"].([]interface{})
				log.Printf("%v", binsCampaigns)
				// User should have all the groupIds that the campaign targets
				numMatch := 0
				for n := range binsCampaigns {
					binCampaign := binsCampaigns[n].(map[interface{}]interface{})
					for k := range binsThisUser {
						userProfile := binsThisUser[k].(map[interface{}]interface{})
						if binCampaign["groupId"] == userProfile["groupId"] {
							numMatch++
						}
					}
				}
				if numMatch == len(binsCampaigns) {
					// User should have at least one same interestId per groupId
					// log.Printf("OK %d %d", numMatch, len(binsCampaigns))
					numMatch = 0
					for n := range binsCampaigns {
						binCampaign := binsCampaigns[n].(map[interface{}]interface{})
						for k := range binsThisUser {
							userProfile := binsThisUser[k].(map[interface{}]interface{})
							if checkInterest(userProfile["interestIds"].([]interface{}), binCampaign["interestIds"].([]interface{})) {
								numMatch++
							}
						}
					}
					if numMatch == len(binsCampaigns) {
						foundOne := recCampaign.Record.Bins["key"].(int)
						log.Printf("MATCH!!! %v", foundOne)
						output = append(output, foundOne)
					}
				}
			}
		}

		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusOK)
		out, _ := json.MarshalIndent(output, "", "  ")
		fmt.Fprintln(response, string(out))
	}
}

// campaign[] of interestId contains at least one item of user[]?
func checkInterest(user, campaign []interface{}) bool {
	//set := make(map[int]int)
	for _, valueCampaign := range campaign {
		for _, valueUser := range user {
			if valueCampaign == valueUser {
				return true
			}
		}
	}
	return false
}
