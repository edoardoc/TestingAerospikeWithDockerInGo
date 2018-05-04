package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"time"

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

// curl -L http://127.0.0.1:8080/validCampaigns?userID=1
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
		defer timeTrack(time.Now())
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

		output := []int{}

		for k := 0; k <= 10; k++ {
			go func() {
				someBigRandomProcessing()
			}()
		}
		log.Printf("all completed...?")

		response.WriteHeader(http.StatusOK)
		out, _ := json.MarshalIndent(output, "", "  ")
		fmt.Fprintln(response, string(out))
	}
}

func someBigRandomProcessing() {
	toWhere := rand.Intn(10000000000)
	log.Printf("someBigRandomProcessing toWhere = %v", toWhere)

	j := 200000000000
	for k := 0; k <= toWhere; k++ {
		j = j * j
	}
	log.Printf("someBigRandomProcessing completed")

}

// all user groupid in this campaign?
func allGroupIDs(binsCampaigns, binsThisUser []interface{}) bool {
	log.Printf("%v", binsCampaigns)
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
	return len(binsCampaigns) == numMatch
}

// User should have at least one same interestId per groupId
func oneInterestPerGroupID(binsCampaigns, binsThisUser []interface{}) bool {
	numMatch := 0
	for n := range binsCampaigns {
		binCampaign := binsCampaigns[n].(map[interface{}]interface{})
		for k := range binsThisUser {
			userProfile := binsThisUser[k].(map[interface{}]interface{})
			if checkInterest(userProfile["interestIds"].([]interface{}), binCampaign["interestIds"].([]interface{})) {
				numMatch++
			}
		}
	}
	return len(binsCampaigns) == numMatch
}

func match(binsCampaigns, binsThisUser []interface{}) bool {
	numMatchGID := 0
	numMatchINT := 0
	for n := range binsCampaigns {
		binCampaign := binsCampaigns[n].(map[interface{}]interface{})
		for k := range binsThisUser {
			userProfile := binsThisUser[k].(map[interface{}]interface{})
			if binCampaign["groupId"] == userProfile["groupId"] {
				numMatchGID++
			}
			if checkInterest(userProfile["interestIds"].([]interface{}), binCampaign["interestIds"].([]interface{})) {
				numMatchINT++
			}
		}
	}
	return len(binsCampaigns) == numMatchGID && len(binsCampaigns) == numMatchINT
}

// campaign[] of interestId contains at least one item of user[]?
func checkInterest(user, campaign []interface{}) bool {
	for _, valueCampaign := range campaign {
		for _, valueUser := range user {
			if valueCampaign == valueUser {
				return true
			}
		}
	}
	return false
}

func timeTrack(start time.Time) {
	elapsed := time.Since(start)

	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path).
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	log.Println(fmt.Sprintf("%s took %s", name, elapsed))
}
