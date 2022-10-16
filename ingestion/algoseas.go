package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type AlgoSeasHistoryItem struct {
	AssetInformation struct {
		Collection struct {
			ColPropertyMetadata struct {
				NonRankableProps   []interface{} `json:"nonRankableProps"`
				NonSearchableProps []interface{} `json:"nonSearchableProps"`
				NonVisualProps     []interface{} `json:"nonVisualProps"`
				PropertyOrder      []interface{} `json:"propertyOrder"`
				RangeProps         []interface{} `json:"rangeProps"`
			} `json:"colPropertyMetadata"`
			Description string `json:"description"`
			Name        string `json:"name"`
			Verified    bool   `json:"verified"`
		} `json:"collection"`
		Image   string `json:"image"`
		NName   string `json:"nName"`
		NRank   int    `json:"nRank"`
		Sk      string `json:"SK"`
		Reserve string `json:"reserve"`
	} `json:"assetInformation"`
	MarketActivity struct {
		AlgoAmount               int    `json:"algoAmount"`
		AssetID                  int    `json:"assetID"`
		CreationDate             string `json:"creationDate"`
		CurrentAuctionAlgoAmount int    `json:"currentAuctionAlgoAmount"`
		Event                    string `json:"event"`
		ExpirationDate           string `json:"expirationDate"`
		GroupTxnID               string `json:"groupTxnID"`
		InitiatorAddress         string `json:"initiatorAddress"`
		ListedAlgoAmount         int    `json:"listedAlgoAmount"`
		MarketplaceID            string `json:"marketplaceID"`
		PreviousOwner            string `json:"previousOwner"`
		TxnID                    string `json:"txnID"`
	} `json:"marketActivity"`
}

type AlgoSeasNote struct {
	Description string `json:"description"`
	ExternalURL string `json:"external_url"`
	MediaURL    string `json:"media_url"`
	MimeType    string `json:"mime_type"`
	Properties  struct {
		Combat           int    `json:"combat"`
		Constitution     int    `json:"constitution"`
		Luck             int    `json:"luck"`
		Plunder          int    `json:"plunder"`
		Scenery          string `json:"Scenery"`
		LeftArm          string `json:"Left Arm"`
		Body             string `json:"Body"`
		BackItem         string `json:"Back Item"`
		Pants            string `json:"Pants"`
		Footwear         string `json:"Footwear"`
		RightArm         string `json:"Right Arm"`
		Shirts           string `json:"Shirts"`
		Hat              string `json:"Hat"`
		HipItem          string `json:"Hip Item"`
		Tattoo           string `json:"Tattoo"`
		Face             string `json:"Face"`
		BackgroundAccent string `json:"Background Accent"`
		Necklace         string `json:"Necklace"`
		Head             string `json:"Head"`
		Background       string `json:"Background"`
		FacialHair       string `json:"Facial Hair"`
		BackHand         string `json:"Back Hand"`
		FrontHand        string `json:"Front Hand"`
		Overcoat         string `json:"Overcoat"`
		Pet              string `json:"Pet"`
	} `json:"properties"`
	Royalties string `json:"royalties"`
	Standard  string `json:"standard"`
}

type AlgoSeasListingItem struct {
	AssetInformation struct {
		NName   string `json:"nName"`
		Sk      string `json:"SK"`
		Listing struct {
			Date            string `json:"date"`
			EscrowAddress   string `json:"escrowAddress"`
			Expires         string `json:"expires"`
			IsDutch         bool   `json:"isDutch"`
			LogicSig        string `json:"logicSig"`
			Marketplace     string `json:"marketplace"`
			MinBidDelta     int    `json:"minBidDelta"`
			NextPayout      string `json:"nextPayout"`
			Seller          string `json:"seller"`
			SnipeThreshold  int    `json:"snipeThreshold"`
			Price           int    `json:"price"`
			Quantity        int    `json:"quantity"`
			Royalty         int    `json:"royalty"`
			RoyaltyString   string `json:"royaltyString"`
			VerifiedRoyalty bool   `json:"verifiedRoyalty"`
			ListingID       int    `json:"listingID"`
			VariableID      int    `json:"variableID"`
		} `json:"listing"`
	} `json:"assetInformation"`
	MarketActivity struct {
		AlgoAmount               int    `json:"algoAmount"`
		AssetID                  int    `json:"assetID"`
		CreationDate             string `json:"creationDate"`
		CurrentAuctionAlgoAmount int    `json:"currentAuctionAlgoAmount"`
		Event                    string `json:"event"`
		ExpirationDate           string `json:"expirationDate"`
		GroupTxnID               string `json:"groupTxnID"`
		InitiatorAddress         string `json:"initiatorAddress"`
		ListedAlgoAmount         int    `json:"listedAlgoAmount"`
		MarketplaceID            string `json:"marketplaceID"`
		PreviousOwner            string `json:"previousOwner"`
		TxnID                    string `json:"txnID"`
	} `json:"marketActivity"`
}

func ReadDataFromJson() []AlgoSeasHistoryItem {
	// Open our jsonFile
	jsonFile, err := os.Open("historicData.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	data := []AlgoSeasHistoryItem{}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &data)
	return data
}

func GetListings() []AlgoSeasListingItem {
	listings := []AlgoSeasListingItem{}

	res, err := http.Get("https://d3ohz23ah7.execute-api.us-west-2.amazonaws.com/prod/marketplace/listings?type=listing&sortBy=price&sortAscending=false&collectionName=AlgoSeas%20Pirates&limit=500")
	if err != nil {
		fmt.Println("Failed to fetch latest listings")
	} else {
		defer res.Body.Close()
		json.NewDecoder(res.Body).Decode(&listings)
	}
	return listings
}

func GetSales() []AlgoSeasHistoryItem {
	parsed := []AlgoSeasHistoryItem{}
	historyUrl := fmt.Sprintf("https://d3ohz23ah7.execute-api.us-west-2.amazonaws.com/prod/marketplace/sales?collectionName=%s&sortBy=time&sortAscending=false&limit=500", url.QueryEscape("AlgoSeas Pirates"))
	res, err := http.Get(historyUrl)
	if err != nil {
		fmt.Errorf("Failed to fetch %s", historyUrl)
	} else {
		defer res.Body.Close()
		json.NewDecoder(res.Body).Decode(&parsed)
	}

	return parsed
}
