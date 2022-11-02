package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/exp/maps"
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
	//Assign initial IdToAsset mapping
	assets := ReadAllAssets(db)
	for _, asset := range assets {
		IdToAsset[asset.ID] = asset
	}
	// PerformClustering(assets)
	PerformKnnSearch(assets, IdToSimilar)
	fmt.Println("Finished initial load")
	return db
}

func processActiveListings(db *sql.DB) []uint64 {
	activeListings := GetListings()
	prevNumberActive := len(IdToListings)
	for k := range IdToListings {
		delete(IdToListings, k)
	}

	for _, listing := range activeListings {
		assetId, err := strconv.ParseUint(listing.AssetInformation.Sk, 10, 64)
		if err != nil {
			fmt.Printf("Failed to convert assetId %s into a Uint", listing.AssetInformation.Sk)
		}
		currentListing, ok := IdToListings[assetId]
		if !ok {
			IdToListings[assetId] = listing
		} else {
			currentCreationDate := ParseDate(currentListing.MarketActivity.CreationDate)
			comparisonCreationDate := ParseDate(listing.MarketActivity.CreationDate)

			if currentCreationDate.Before(comparisonCreationDate) {
				fmt.Println("Higher Listing found! ", listing.AssetInformation.Listing.ListingID)
				IdToListings[assetId] = listing
			}
		}
	}

	fmt.Printf("ActiveListings changed length: %d (change of %d)\n", len(activeListings), len(activeListings)-prevNumberActive)
	keys := maps.Keys(IdToListings)
	return keys
}

func startPolling(db *sql.DB) {
	tick := time.Tick(1 * time.Minute)
	for range tick {
		//Update Metadata && Insert newly minted tokens
		newAssets := GetNewMetadata(db)
		for _, asset := range newAssets {
			InsertAsset(db, asset)
			IdToAsset[asset.ID] = asset
		}
		fmt.Printf("Ingested %d new asset updates!\n", len(newAssets))
		assetsWithListings := processActiveListings(db)
		assets := ReadAllAssets(db)
		PerformKnnSearch(assets, IdToSimilar)
		PerformKnnSearch(AssetIdsToAssets(assetsWithListings), IdToSimilarActive)
	}

	fmt.Println()

}

func main() {
	db := initialLoad()
	defer db.Close()
	go startPolling(db)
	mux := http.NewServeMux()
	mux.HandleFunc("/similar", SimilarAssetsHandler)
	mux.HandleFunc("/assets", AssetHandler)
	http.ListenAndServe(":8080", mux)
	fmt.Println("Server listening on port 8080")
}
