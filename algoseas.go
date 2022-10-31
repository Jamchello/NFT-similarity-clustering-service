package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

type AlgoSeasListingData struct {
	Date            string `json:"date"`
	EscrowAddress   string `json:"escrowAddress,omitempty"`
	Expires         string `json:"expires,omitempty"`
	IsDutch         bool   `json:"isDutch"`
	LogicSig        string `json:"logicSig,omitempty"`
	Marketplace     string `json:"marketplace,omitempty"`
	MinBidDelta     int    `json:"minBidDelta,omitempty"`
	NextPayout      string `json:"nextPayout,omitempty"`
	Seller          string `json:"seller,omitempty"`
	SnipeThreshold  int    `json:"snipeThreshold,omitempty"`
	Price           int    `json:"price,omitempty"`
	Quantity        int    `json:"quantity,omitempty"`
	Royalty         int    `json:"royalty,omitempty"`
	RoyaltyString   string `json:"royaltyString,omitempty"`
	VerifiedRoyalty bool   `json:"verifiedRoyalty,omitempty"`
	ListingID       int    `json:"listingID"`
	VariableID      int    `json:"variableID,omitempty"`
}

type AlgoSeasListingsAsset struct {
	AssetInformation struct {
		NName   string              `json:"nName"`
		Sk      string              `json:"SK"`
		Listing AlgoSeasListingData `json:"listing"`
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

func GetListings() []AlgoSeasListingsAsset {
	listings := []AlgoSeasListingsAsset{}
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
		fmt.Printf("Failed to fetch %s\n", historyUrl)
	} else {
		defer res.Body.Close()
		json.NewDecoder(res.Body).Decode(&parsed)
	}

	return parsed
}
