// TODO: Active listings
// TODO: USD ammt
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

func initialLoad() (*sql.DB, map[string]bool) {

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
	return db, seenAssets
}

func pollNewData(db *sql.DB) ([]Sale, []Asset) {

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
	if len(newSales) == 0 {
		fmt.Println("no new sales...")
		return newSales, newAssets
	}
	for _, sale := range newSales {
		InsertSale(db, sale)
	}
	return newSales, newAssets
}

func main() {
	db, _ := initialLoad()
	defer db.Close()

	tick := time.Tick(1 * time.Minute)
	for range tick {
		sales, assets := pollNewData(db)
		fmt.Printf("Ingested %d sales, %d new asset updates\n", len(sales), len(assets))
	}

}
