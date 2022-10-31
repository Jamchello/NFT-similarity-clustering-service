package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
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
		Sk      string `json:"SK,omitempty"`
		Image   string `json:"image,omitempty"`
		NName   string `json:"nName,omitempty"`
		NRank   int    `json:"nRank,omitempty"`
		Total   int    `json:"total,omitempty"`
		Reserve string `json:"reserve,omitempty"`
		NProps  struct {
			Properties struct {
				Constitution     int    `json:"constitution,omitempty"`
				Luck             int    `json:"luck,omitempty"`
				HipItem          bool   `json:"Hip Item,omitempty"`
				Tattoo           string `json:"Tattoo,omitempty"`
				Combat           int    `json:"combat,omitempty"`
				RightArm         string `json:"Right Arm,omitempty"`
				Shirts           string `json:"Shirts,omitempty"`
				Shirt            bool   `json:"Shirt,omitempty"`
				BackItem         bool   `json:"Back Item,omitempty"`
				Face             bool   `json:"Face,omitempty"`
				Body             string `json:"Body,omitempty"`
				BackgroundAccent bool   `json:"Background Accent,omitempty"`
				Necklace         bool   `json:"Necklace,omitempty"`
				Head             bool   `json:"Head,omitempty"`
				Scenery          string `json:"Scenery,omitempty"`
				Background       bool   `json:"Background,omitempty"`
				Footwear         string `json:"Footwear,omitempty"`
				Pants            string `json:"Pants,omitempty"`
				FacialHair       bool   `json:"Facial Hair,omitempty"`
				LeftArm          string `json:"Left Arm,omitempty"`
				BackHand         bool   `json:"Back Hand,omitempty"`
				FrontHand        bool   `json:"Front Hand,omitempty"`
				Plunder          int    `json:"plunder,omitempty"`
				Overcoat         bool   `json:"Overcoat,omitempty"`
				Hat              string `json:"Hat,omitempty"`
				Pet              bool   `json:"Pet,omitempty"`
			} `json:"properties,omitempty"`
		} `json:"nProps,omitempty"`
		Listing struct {
			Date            time.Time `json:"date,omitempty"`
			EscrowAddress   string    `json:"escrowAddress,omitempty"`
			Expires         string    `json:"expires,omitempty"`
			IsDutch         bool      `json:"isDutch,omitempty"`
			LogicSig        string    `json:"logicSig,omitempty"`
			Marketplace     string    `json:"marketplace,omitempty"`
			MinBidDelta     int       `json:"minBidDelta,omitempty"`
			NextPayout      string    `json:"nextPayout,omitempty"`
			Seller          string    `json:"seller,omitempty"`
			SnipeThreshold  int       `json:"snipeThreshold,omitempty"`
			Price           int       `json:"price,omitempty"`
			Quantity        int       `json:"quantity,omitempty"`
			Royalty         int       `json:"royalty,omitempty"`
			RoyaltyString   string    `json:"royaltyString,omitempty"`
			VerifiedRoyalty bool      `json:"verifiedRoyalty,omitempty"`
			ListingID       int       `json:"listingID,omitempty"`
			VariableID      int       `json:"variableID,omitempty"`
		} `json:"listing,omitempty"`
	} `json:"assetInformation,omitempty"`
	MarketActivity struct {
		AlgoAmount               int    `json:"algoAmount,omitempty"`
		AssetID                  int    `json:"assetID,omitempty"`
		CreationDate             string `json:"creationDate,omitempty"`
		CurrentAuctionAlgoAmount int    `json:"currentAuctionAlgoAmount,omitempty"`
		Event                    string `json:"event,omitempty"`
		ExpirationDate           string `json:"expirationDate,omitempty"`
		GroupTxnID               string `json:"groupTxnID,omitempty"`
		InitiatorAddress         string `json:"initiatorAddress,omitempty"`
		ListedAlgoAmount         int    `json:"listedAlgoAmount,omitempty"`
		MarketplaceID            string `json:"marketplaceID,omitempty"`
		PreviousOwner            string `json:"previousOwner,omitempty"`
		TxnID                    string `json:"txnID,omitempty"`
	} `json:"marketActivity,omitempty"`
}

type AlgoSeasListingResponse struct {
	Assets    []AlgoSeasListingItem `json:"assets,omitempty"`
	NextToken string                `json:"nextToken,omitempty"`
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

func GetListings(listings []AlgoSeasListingItem, token string) []AlgoSeasListingItem {
	resp := AlgoSeasListingResponse{}
	url := fmt.Sprintf("https://d3ohz23ah7.execute-api.us-west-2.amazonaws.com/prod/marketplace/v2/assetsByCollection/AlgoSeas%%20Pirates?type=listing&sortBy=price&sortAscending=true&limit=200&nextToken=%s", token)
	fmt.Println(url)
	res, err := http.Post(url, "application/json", &bytes.Buffer{})
	if err != nil {
		fmt.Println("Failed to fetch latest listings")
	} else {
		defer res.Body.Close()
		json.NewDecoder(res.Body).Decode(&resp)
		listings = append(listings, resp.Assets...)
		if resp.NextToken != "" {
			return GetListings(listings, resp.NextToken)
		}
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
