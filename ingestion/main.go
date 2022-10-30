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
func generate2DArray(x int) [][]float64 {
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
			num := float64(rand.Intn(max+1-min) + min)
			newArr = append(newArr,num)
			
		}
	
		//then add that new 4 num array to the 2d array
		testArr = append(testArr, newArr)
		
	}

	fmt.Println(testArr)
	return testArr

}

func testClusters() ([]int, [][]float64){

	data := generate2DArray(10)
	copydata := generate2DArray(10)


	var observation []float64

	c,e := clusters.KMeans(20,5, clusters.EuclideanDistance)
	if e!= nil{
		panic(e)
	}

// Use the data to train the clusterer
	if e = c.Learn(copydata); e != nil {
		panic(e)
	}

	fmt.Printf("Clustered data set into %d\n", c.Sizes())

	fmt.Printf("Assigned observation %v to cluster %d\n", observation, c.Predict(observation))
	

	for index, number := range c.Guesses() {
		fmt.Println("Index ",index,"Assigned data point ",data[index]," to cluster ",number)
	}

	fmt.Println(c.Guesses())
	//c.Guesses() contains list of cluster for each asset (index) in original indexing order. So c.Guesses()[0] tells us the cluster that data[0] belongs to and so on
	//can use to create hashmap
	return c.Guesses(), data



}

//TODO: Implement function to take data -> cluster mapping, create dedicated hashmap from asset -> cluster
//Figure out how to link back each asset characteristic array to its ID (possibly ignored fields)
//possibly two hashmaps
//original asset -> asset characteristic array then asset characteristic array -> cluster


//Currently implements a map from cluster -> list of assets for that cluster using the outputs from the testClusters method
func compare(cluster []int, data [][]float64 ){
	ClusterToAsset := make(map[int][][]float64)
	for index, number := range cluster{
		//fetch current asset from data list using the fact that data[i] corresponds with cluster[i]
		currentAsset := data[index]
		//fetch assets for the current cluster that are already in the map
		current_clusterAssets := ClusterToAsset[number]
		//append new asset to asset list
		new_clusterAssets:= append(current_clusterAssets, currentAsset)
		//Put the new asset list (value) in for the current cluster (key)
		ClusterToAsset[number] = new_clusterAssets 
	}

	//now for each cluster we can fetch all assets within the same cluster
	// could be useful for comparison ??
	
}

func deepCopyArr(origarr [][]float64) [][]float64{
	copyarr:= [][]float64{}
	for i:=0; i<len(origarr); i++{
		copyarr = append(copyarr, origarr[i])
	}

	return copyarr
}

func testEffect(){
	origarr := generate2DArray(4)
	copyarr:= generate2DArray(4)

	fmt.Println("before")
	fmt.Println(origarr)
	fmt.Println(copyarr)	
	c,e := clusters.KMeans(20,5, clusters.EuclideanDistance)
	if e!= nil{
		panic(e)
	}
	if e = c.Learn(copyarr); e != nil {
		panic(e)
	}
	fmt.Println("after")
	fmt.Println(origarr)
	fmt.Println(copyarr)	
}

func main() {

	testClusters()


}
