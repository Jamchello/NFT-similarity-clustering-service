// TODO: Active listings
// TODO: USD ammt
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

var listingsCache = map[string]AlgoSeasListingItem{}

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

		historicalSalesData := ReadDataFromJson()
		mintedAssets := GetAllMintedAssets()
		for _, asset := range mintedAssets {
			err := InsertAsset(db, asset)
			if err != nil {
				fmt.Println("Failed to insert asset", err)
			}
		}
		for _, sale := range historicalSalesData {
			sale := Sale{
				Date:   ParseSaleDate(sale.MarketActivity.CreationDate),
				Tx:     sale.MarketActivity.TxnID,
				Buyer:  sale.MarketActivity.InitiatorAddress,
				Seller: sale.MarketActivity.PreviousOwner,
				Algo:   fmt.Sprintf("%d", sale.MarketActivity.AlgoAmount),
				Fiat:   "",
				Asset:  uint64(sale.MarketActivity.AssetID),
			}

			err := InsertSale(db, sale)
			if err != nil {
				fmt.Println("Failed to insert sale", err)
			}
		}
	} else {
		loadAssetIds(db, seenAssets)
	}
	fmt.Println("Finished initial load")
	return db
}

func startPolling(db *sql.DB) {
	tick := time.Tick(1 * time.Minute)
	for range tick {
		//Update Metadata && Insert newly minted tokens
		newAssets := GetNewMetadatas(db)
		for _, asset := range newAssets {
			InsertAsset(db, asset)
		}
		lastIngestedSale := getLatestIngestedSale(db)
		sales := GetSales()
		newSales := []Sale{}
		for _, item := range sales {
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
			}
		}

		for _, sale := range newSales {
			InsertSale(db, sale)
		}

		for k := range listingsCache {
			delete(listingsCache, k)
		}

		activeListings := GetListings()
		for _, listing := range activeListings {
			assetId := listing.AssetInformation.Sk
			currentListing, ok := listingsCache[assetId]
			if !ok {
				listingsCache[assetId] = listing
			} else {
				currentCreationDate := ParseSaleDate(currentListing.MarketActivity.CreationDate)
				comparisonCreationDate := ParseSaleDate(listing.MarketActivity.CreationDate)

				if currentCreationDate.Before(comparisonCreationDate) {
					fmt.Println("Higher Listing found! ", listing.AssetInformation.Listing.ListingID)
					listingsCache[assetId] = listing
				}
			}
		}

		fmt.Printf("Ingested %d sales, %d new asset updates, listingsCache replenished to %d\n", len(newSales), len(newAssets), len(activeListings))
	}
}

// Temporary handler to debug, need to flesh out the actual handler once we have the data...
func ListingsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		assetId := r.URL.Query().Get("assetId")
		if assetId == "" {
			http.Error(w, "Invalid assetId", http.StatusBadRequest)
			return
		}

		listing, ok := listingsCache[assetId]

		if !ok {
			http.Error(w, "No Listing stored for this asset", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(listing)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	db := initialLoad()
	defer db.Close()

	go startPolling(db)

	mux := http.NewServeMux()

	mux.HandleFunc("/listing", ListingsHandler)

	http.ListenAndServe(":8080", mux)

}
