package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Asset struct {
	ID           uint64
	Collection   string
	Image_Url    string
	Combat       uint64
	Constitution uint64
	Luck         uint64
	Plunder      uint64
	Properties   string
}

type Sale struct {
	Date   time.Time
	Tx     string
	Buyer  string
	Seller string
	Algo   string
	Fiat   string
	Asset  uint64
}

func ParseSaleDate(timestamp string) time.Time {
	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, timestamp)

	if err != nil {
		fmt.Println(err)
	}
	return t
}

func InsertAsset(db *sql.DB, asset Asset) error {
	stmt, err := db.Prepare("INSERT IGNORE INTO asset VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		asset.ID,
		asset.Collection,
		asset.Image_Url,
		asset.Combat,
		asset.Constitution,
		asset.Luck,
		asset.Plunder,
		asset.Properties)

	if err != nil {
		return err
	}
	return nil
}

func InsertSale(db *sql.DB, sale Sale) error {
	stmt, err := db.Prepare("INSERT IGNORE INTO sale(Date,Tx, Buyer, Seller, Algo, Fiat, Asset) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		sale.Date.Format("2006-01-02 15:04:05"),
		sale.Tx,
		sale.Buyer,
		sale.Seller,
		sale.Algo,
		sale.Fiat,
		sale.Asset,
	)

	if err != nil {
		return err
	}
	return nil
}

func createDb(db *sql.DB) {
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS algoseas")
	if err != nil {
		log.Fatal(err)
	}
}

func createAssetTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `asset` ( `ID` INT unsigned, `Collection` VARCHAR(255) NOT NULL, `Image_Url` TEXT NOT NULL, `Combat` INT unsigned NOT NULL, `Constitution` INT unsigned NOT NULL, `Luck` INT unsigned NOT NULL, `Plunder` INT unsigned NOT NULL, `Properties` TEXT NOT NULL, KEY `Collection_ID_IDX` (`Collection`,`ID`) USING BTREE, PRIMARY KEY (`ID`) );")
	if err != nil {
		log.Fatal(err)
	}
}

func createSaleTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `sale` ( `ID` INT unsigned NOT NULL AUTO_INCREMENT, `TX` VARCHAR(52) NOT NULL, `Date` DATETIME NOT NULL, `Buyer` VARCHAR(58) NOT NULL, `Seller` VARCHAR(58) NOT NULL, `Algo` BIGINT NOT NULL, `Fiat` FLOAT, `Asset` INT unsigned, PRIMARY KEY (ID), CONSTRAINT fk_asset FOREIGN KEY (Asset) REFERENCES asset(ID) ON DELETE CASCADE ON UPDATE CASCADE );")
	if err != nil {
		log.Fatal(err)
	}
}

func dbNeedsPopulating(db *sql.DB) bool {
	res, _ := db.Query("SELECT * FROM `asset`")
	return !res.Next()
}

func loadAssetIds(db *sql.DB, seenAssets map[string]bool) {

	id := ""
	rows, _ := db.Query("SELECT Id FROM asset")
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Fatalln(err)
		}
		seenAssets[id] = true
	}

}

func getLatestIngestedSale(db *sql.DB) time.Time {
	date := time.Time{}
	rows, _ := db.Query("SELECT Date FROM sale ORDER BY Date DESC LIMIT 1")
	for rows.Next() {
		err := rows.Scan(&date)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return date
}
