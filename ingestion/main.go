package main

import (
	"database/sql"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type AlgoSeasHistoryResp []struct {
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

type IndexerSearchResp struct {
	CurrentRound int    `json:"current-round"`
	NextToken    string `json:"next-token"`
	Transactions []struct {
		ApplicationTransaction struct {
			Accounts          []interface{} `json:"accounts"`
			ApplicationArgs   []string      `json:"application-args"`
			ApplicationID     int           `json:"application-id"`
			ForeignApps       []interface{} `json:"foreign-apps"`
			ForeignAssets     []interface{} `json:"foreign-assets"`
			GlobalStateSchema struct {
				NumByteSlice int `json:"num-byte-slice"`
				NumUint      int `json:"num-uint"`
			} `json:"global-state-schema"`
			LocalStateSchema struct {
				NumByteSlice int `json:"num-byte-slice"`
				NumUint      int `json:"num-uint"`
			} `json:"local-state-schema"`
			OnCompletion string `json:"on-completion"`
		} `json:"application-transaction,omitempty"`
		CloseRewards     int    `json:"close-rewards"`
		ClosingAmount    int    `json:"closing-amount"`
		ConfirmedRound   int    `json:"confirmed-round"`
		Fee              int    `json:"fee"`
		FirstValid       int    `json:"first-valid"`
		GenesisHash      string `json:"genesis-hash"`
		GenesisID        string `json:"genesis-id"`
		GlobalStateDelta []struct {
			Key   string `json:"key"`
			Value struct {
				Action int `json:"action"`
				Uint   int `json:"uint"`
			} `json:"value"`
		} `json:"global-state-delta,omitempty"`
		Group     string `json:"group"`
		ID        string `json:"id"`
		InnerTxns []struct {
			AssetConfigTransaction struct {
				AssetID int `json:"asset-id"`
				Params  struct {
					Creator       string `json:"creator"`
					Decimals      int    `json:"decimals"`
					DefaultFrozen bool   `json:"default-frozen"`
					Manager       string `json:"manager"`
					Name          string `json:"name"`
					NameB64       string `json:"name-b64"`
					Reserve       string `json:"reserve"`
					Total         int    `json:"total"`
					UnitName      string `json:"unit-name"`
					UnitNameB64   string `json:"unit-name-b64"`
					URL           string `json:"url"`
					URLB64        string `json:"url-b64"`
				} `json:"params"`
			} `json:"asset-config-transaction"`
			CloseRewards      int    `json:"close-rewards"`
			ClosingAmount     int    `json:"closing-amount"`
			ConfirmedRound    int    `json:"confirmed-round"`
			CreatedAssetIndex int    `json:"created-asset-index"`
			Fee               int    `json:"fee"`
			FirstValid        int    `json:"first-valid"`
			IntraRoundOffset  int    `json:"intra-round-offset"`
			LastValid         int    `json:"last-valid"`
			ReceiverRewards   int    `json:"receiver-rewards"`
			RoundTime         int    `json:"round-time"`
			Sender            string `json:"sender"`
			SenderRewards     int    `json:"sender-rewards"`
			TxType            string `json:"tx-type"`
		} `json:"inner-txns,omitempty"`
		IntraRoundOffset int    `json:"intra-round-offset"`
		LastValid        int    `json:"last-valid"`
		ReceiverRewards  int    `json:"receiver-rewards"`
		RoundTime        int    `json:"round-time"`
		Sender           string `json:"sender"`
		SenderRewards    int    `json:"sender-rewards"`
		Signature        struct {
			Sig string `json:"sig"`
		} `json:"signature"`
		TxType                 string `json:"tx-type"`
		AssetConfigTransaction struct {
			AssetID int `json:"asset-id"`
			Params  struct {
				Creator       string `json:"creator"`
				Decimals      int    `json:"decimals"`
				DefaultFrozen bool   `json:"default-frozen"`
				Manager       string `json:"manager"`
				Reserve       string `json:"reserve"`
				Total         int    `json:"total"`
			} `json:"params"`
		} `json:"asset-config-transaction,omitempty"`
		Note string `json:"note,omitempty"`
	} `json:"transactions"`
}

type AlgoSeasNote struct {
	Description string `json:"description"`
	ExternalURL string `json:"external_url"`
	MediaURL    string `json:"media_url"`
	MimeType    string `json:"mime_type"`
	Properties  struct {
		Scenery      string `json:"Scenery"`
		LeftArm      string `json:"Left Arm"`
		Body         string `json:"Body"`
		BackItem     string `json:"Back Item"`
		Pants        string `json:"Pants"`
		Footwear     string `json:"Footwear"`
		RightArm     string `json:"Right Arm"`
		Shirts       string `json:"Shirts"`
		Hat          string `json:"Hat"`
		Combat       int    `json:"combat"`
		Constitution int    `json:"constitution"`
		Luck         int    `json:"luck"`
		Plunder      int    `json:"plunder"`
	} `json:"properties"`
	Royalties string `json:"royalties"`
	Standard  string `json:"standard"`
}

func getMetadataForAsset(assetId string) AlgoSeasNote {
	note := AlgoSeasNote{}
	indexerUrl := fmt.Sprintf("https://algoindexer.algoexplorerapi.io/v2/transactions?asset-id=%s&tx-type=acfg", assetId)
	res, err := http.Get(indexerUrl)
	if err != nil {
		fmt.Printf("Retrying fetch %s in 3 seconds\n", assetId)
		time.Sleep(3 * time.Second)
		return getMetadataForAsset(assetId)
	} else {
		defer res.Body.Close()
		parsed := IndexerSearchResp{}
		json.NewDecoder(res.Body).Decode(&parsed)
		transactions := parsed.Transactions
		// sort.Slice(transactions, func(i, j int) bool {
		// 	return transactions[i].RoundTime > transactions[j].RoundTime
		// })
		for _, tx := range transactions {
			if tx.Note != "" {
				noteBytes, _ := b64.StdEncoding.DecodeString(tx.Note)
				err := json.Unmarshal(noteBytes, &note)
				if err != nil {
					return note
				}
				if note.Standard == "arc69" {
					return note
				}
			}
		}
	}
	return note
}

func readHistoricalDataFromJson() AlgoSeasHistoryResp {
	// Open our jsonFile
	jsonFile, err := os.Open("historicData.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened historicData.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	data := AlgoSeasHistoryResp{}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &data)
	return data
}

func initialLoad() (*sql.DB, map[string]bool) {
	db, err := sql.Open("mysql",
		"pirate:test_password@tcp(localhost:3306)/algoseas?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	createDb(db)
	createAssetTable(db)
	createSaleTable(db)
	seenAssets := make(map[string]bool)
	if dbNeedsPopulating(db) {
		metadataNeeded := []string{}
		historicalSalesData := readHistoricalDataFromJson()
		for _, item := range historicalSalesData {
			uniqueId := item.AssetInformation.Sk
			if !seenAssets[uniqueId] {
				seenAssets[uniqueId] = true
				metadataNeeded = append(metadataNeeded, uniqueId)
			}
		}

		// assets := []Asset{}
		assetChannel := make(chan Asset, len(metadataNeeded))
		wg := &sync.WaitGroup{}
		wg.Add(len(metadataNeeded))

		for _, assetId := range metadataNeeded {
			go func(assetId string) {
				note := getMetadataForAsset(assetId)
				idInt, _ := strconv.ParseUint(assetId, 10, 64)
				asset := Asset{
					ID:           idInt,
					Collection:   "AlgoSeas Pirates", //TODO: Make dynamic by passing in
					Image_Url:    note.MediaURL,
					Combat:       uint64(note.Properties.Combat),
					Constitution: uint64(note.Properties.Constitution),
					Luck:         uint64(note.Properties.Luck),
					Plunder:      uint64(note.Properties.Plunder),
					Properties:   "{}", //TODO: Extract properties correctly
				}
				assetChannel <- asset

			}(assetId)
		}

		for j := 0; j < len(metadataNeeded); j++ {
			// assets = append(assets, <-assetChannel)
			err := InsertAsset(db, <-assetChannel)
			if err != nil {
				fmt.Println("Failed to insert", err)
			}
		}

		for _, item := range historicalSalesData {
			sale := Sale{
				Date:   ParseSaleDate(item.MarketActivity.CreationDate),
				Tx:     item.MarketActivity.TxnID,
				Buyer:  item.MarketActivity.InitiatorAddress,
				Seller: item.MarketActivity.PreviousOwner,
				Algo:   fmt.Sprintf("%d", item.MarketActivity.AlgoAmount),
				Fiat:   "",
				Asset:  uint64(item.MarketActivity.AssetID),
			}

			err := InsertSale(db, sale)
			if err != nil {
				fmt.Println("Failed to insert sale", err)
			}
		}
	} else {
		loadAssetIds(db, seenAssets)
	}
	return db, seenAssets
}

func getSales() AlgoSeasHistoryResp {
	parsed := AlgoSeasHistoryResp{}
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

func pollNewData(db *sql.DB, ingestedAssets map[string]bool) ([]Sale, []Asset) {
	lastIngestedSale := getLatestIngestedSale(db)
	sales := getSales()
	newSales := []Sale{}
	needsMetaData := []string{}
	newAssets := []Asset{}
	for _, item := range sales {
		assetId := item.AssetInformation.Sk
		saleTime := ParseSaleDate(item.MarketActivity.CreationDate)
		if lastIngestedSale.Before(saleTime) {
			newSale := Sale{
				Date:   saleTime,
				Tx:     item.MarketActivity.TxnID,
				Buyer:  item.MarketActivity.InitiatorAddress,
				Seller: item.MarketActivity.PreviousOwner,
				Algo:   fmt.Sprintf("%d", item.MarketActivity.AlgoAmount),
				Fiat:   "",
				Asset:  uint64(item.MarketActivity.AssetID),
			}
			newSales = append(newSales, newSale)
			fmt.Println(newSale.Date)
			if !ingestedAssets[assetId] {
				needsMetaData = append(needsMetaData, assetId)
				ingestedAssets[assetId] = true
			}
		}
	}
	if len(newSales) == 0 {
		fmt.Println("no new sales...")
		return newSales, newAssets
	}

	assetChannel := make(chan Asset)
	for _, assetId := range needsMetaData {
		go func(assetId string) {
			note := getMetadataForAsset(assetId)
			idInt, _ := strconv.ParseUint(assetId, 10, 64)
			asset := Asset{
				ID:           idInt,
				Collection:   "AlgoSeas Pirates", //TODO: Make dynamic by passing in
				Image_Url:    note.MediaURL,
				Combat:       uint64(note.Properties.Combat),
				Constitution: uint64(note.Properties.Constitution),
				Luck:         uint64(note.Properties.Luck),
				Plunder:      uint64(note.Properties.Plunder),
				Properties:   "{}", //TODO: Extract properties correctly
			}
			assetChannel <- asset
			err := InsertAsset(db, asset)
			if err != nil {
				fmt.Println(err)
			}
		}(assetId)
	}

	for j := 0; j < len(needsMetaData); j++ {
		newAssets = append(newAssets, <-assetChannel)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(newSales))
	for _, sale := range newSales {
		go func(sale Sale) {
			InsertSale(db, sale)
			wg.Done()
		}(sale)
	}
	wg.Wait()
	return newSales, newAssets
}

func main() {
	db, ingestedAssets := initialLoad()
	defer db.Close()

	tick := time.Tick(1 * time.Minute)
	for range tick {
		sales, assets := pollNewData(db, ingestedAssets)
		fmt.Printf("Ingested %x sales, and %x new assets", len(sales), len(assets))
	}
}
