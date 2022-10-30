// TODO: Active listings
// TODO: USD ammt
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mpraski/clusters"
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
					Fiat:   0,
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




//have to use float64 for KMeans
func generateTest2DArray(x int) [][]float64 {
	testArr := [][]float64{}
	min:= 0
	max:= 100
	seedNum:= int64(1)
	//generate x amount of 4  num arrs
	for i:=0; i< x; i++{
		//GENERATE ARRAY OF 4 NUMBERS
		newArr:= []float64{}
		for j:=0; j<4; j++{
			seedNum ++
			rand.Seed(time.Now().UnixNano() + seedNum)
			//rand.Intn gens numbers from (0,n) exclusive so we +1 on Intn and +min after so we have min as our minimum value (rather than 0)
			newArr = append(newArr,float64(rand.Intn(max+1-min) + min))
		}
		//then add that new 4 num array to the 2d array
		testArr = append(testArr, newArr)
	}

	return testArr

}

func testClusters(){

	data := generateTest2DArray(200)
	var observation []float64

	c,e := clusters.KMeans(200000,8, clusters.EuclideanDistance)
	if e!= nil{
		panic(e)
	}

	c.Learn(data);

	fmt.Printf("Clustered data set into %d\n", c.Sizes())

	fmt.Printf("Assigned observation %v to cluster %d\n", observation, c.Predict(observation))
	
	for index, number := range c.Guesses() {
		fmt.Printf("Assigned data point %v to cluster %d\n", data[index], number)
	}
}


// func compareListings(){

// 	testArr := [][]int32{{20,43,52,70},{10,30,55,24},{42,63,13,78}}
// 	c, e := clusters.KMeans(1000, 4, clusters.EuclideanDistance)
// 	if e !=nil{
// 		panic(e)
// 	}

// 	//otherwise
// 	if e = c.learn

// }

func main() {

	testClusters()

}
