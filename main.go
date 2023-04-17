package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/imroc/req/v3"
	"log"
	"os"
//	"time"
)

type StringDict map[string]string

type AccountNumber StringDict

type AccountPost struct {
	AsOfDate     string          `json:"as_of_date"`
	BalanceTypes []string        `json:"balance_types"`
	Accounts     []AccountNumber `json:"accounts"`
	BankId       *string         `json:"bank_id,omitempty"`
}

type Result struct {
	Data string `json:"data"`
}

const (
	YYYYMMDD = "2006-01-02"
	URL      = "https://api-sandbox.wellsfargo.com/treasury/balance-reporting/v1/balances/report"
	FILENAME = "account.json"
)

func main() {

	// create variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// grab environment variables
	consumerKey := os.Getenv("Consumerkey")
//	consumerSecret := os.Getenv("Consumersecret")
	apiKey := os.Getenv("APIkey")

	// Open our jsonFile
	var accountData AccountPost
	fileBytes, _ := os.ReadFile(FILENAME)
	err = json.Unmarshal(fileBytes, &accountData)
	if err != nil {
		log.Fatal(err)
	}

	accountData.AsOfDate = "2023-04-14" // time.Now().UTC().Format(YYYYMMDD)
	fmt.Print(accountData)

	//
	client := req.C().DevMode()
	var result Result

	resp, err := client.R().
		SetHeader("Authorization",  fmt.Sprintf("Bearer %s", apiKey)).
		SetHeader("Content-Type", "application/json").
		SetHeader("gateway-entity-id", consumerKey).
		SetHeader("client-request-id", uuid.New().String()).
		SetBody(accountData).
		SetSuccessResult(&result).
		Post(URL)

	if err != nil {
		log.Fatal(err)
	}

	if !resp.IsSuccessState() {
		fmt.Println("bad response status:", resp.Status)
		return
	}
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println("data:", result.Data)
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++")
}
