// TODO: Active listings
// TODO: USD ammt
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

func initialLoad() *sql.DB {

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	if dbUser == "" || dbPassword == "" {
		log.Fatal("Failed to load DB credentials from environment variables (DB_USER, DB_PASSWORD)")
	}

	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(localhost:3306)/algoseas?parseTime=true", dbUser, dbPassword))
	if err != nil {
		log.Fatal(err)
	}
	createDb(db)
	createAssetTable(db)
	createSaleTable(db)
	seenAssets := make(map[string]bool)
	if dbNeedsPopulating(db) {
		fmt.Println("No Data detected, populating the database")
		mintedAssets := GetAllMintedAssets()
		for _, asset := range mintedAssets {
			err := InsertAsset(db, asset)
			if err != nil {
				fmt.Println("Failed to insert asset", err)
			}
		}
	} else {
		loadAssetIds(db, seenAssets)
	}
	assets := ReadAllAssets(db)
	UpdateAssetsMapping(assets)
	PerformClustering(assets)
	fmt.Println("Finished initial load")
	return db
}

func startPolling() {
	tick := time.Tick(1 * time.Minute)
	for range tick {
		//Update Metadata && Insert newly minted tokens
		newAssets := GetNewMetadatas(Database)
		for _, asset := range newAssets {
			InsertAsset(Database, asset)
		}
		lastIngestedSale := getLatestIngestedSale(Database)
		sales := GetSales()
		newSales := []Sale{}

		prevNumberActive := len(IdToListings)
		for _, item := range sales {
			saleTime := ParseSaleDate(item.MarketActivity.CreationDate)
			if lastIngestedSale.Before(saleTime) {
				newSale := Sale{
					Date:   saleTime,
					Tx:     item.MarketActivity.TxnID,
					Buyer:  item.MarketActivity.InitiatorAddress,
					Seller: item.MarketActivity.PreviousOwner,
					Algo:   fmt.Sprintf("%d", item.MarketActivity.AlgoAmount),
					Fiat:   0,
					Asset:  uint64(item.MarketActivity.AssetID),
				}
				newSales = append(newSales, newSale)
			}
		}

		for _, sale := range newSales {
			InsertSale(Database, sale)
		}

		for k := range IdToListings {
			delete(IdToListings, k)
		}

		activeListings := GetListings()
		for _, listing := range activeListings {
			assetId := listing.AssetInformation.Sk
			currentListing, ok := IdToListings[assetId]
			if !ok {
				IdToListings[assetId] = listing
			} else {
				currentCreationDate := ParseSaleDate(currentListing.MarketActivity.CreationDate)
				comparisonCreationDate := ParseSaleDate(listing.MarketActivity.CreationDate)

				if currentCreationDate.Before(comparisonCreationDate) {
					fmt.Println("Higher Listing found! ", listing.AssetInformation.Listing.ListingID)
					IdToListings[assetId] = listing
				}
			}
		}

		fmt.Printf("Ingested %d sales, %d new asset updates, ActiveListings changed length: %d (change of %d)\n", len(newSales), len(newAssets), len(activeListings), len(activeListings)-prevNumberActive)
		if len(newAssets) > 0 {
			//Re-perform clustering if new assets minted
			assets := ReadAllAssets(Database)
			UpdateAssetsMapping(assets)
			PerformClustering(assets)
		}
	}
}

func main() {
	Database := initialLoad()
	defer Database.Close()
	go startPolling()
	mux := http.NewServeMux()
	mux.HandleFunc("/similar", SimilarAssetsHandler)
	//TODO: asset endpoint to get data on an asset
	http.ListenAndServe(":8080", mux)
}
