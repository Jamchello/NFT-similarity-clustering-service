package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

type Asset struct {
	ID               uint64 `json:"Id"`
	UpdatedAt        uint   `json:"UpdatedAt,omitempty"`
	Collection       string `json:"Collection"`
	ImageUrl         string `json:"ImageUrl"`
	Combat           uint64 `json:"Combat"`
	Constitution     uint64 `json:"Constitution"`
	Luck             uint64 `json:"Luck"`
	Plunder          uint64 `json:"Plunder"`
	Scenery          string `json:"Scenery,omitempty"`
	LeftArm          string `json:"LeftArm,omitempty"`
	Body             string `json:"Body,omitempty"`
	BackItem         string `json:"BackItem,omitempty"`
	Pants            string `json:"Pants,omitempty"`
	Footwear         string `json:"Footwear,omitempty"`
	RightArm         string `json:"RightArm,omitempty"`
	Shirts           string `json:"Shirts,omitempty"`
	Hat              string `json:"Hat,omitempty"`
	HipItem          string `json:"HipItem,omitempty"`
	Tattoo           string `json:"Tattoo,omitempty"`
	Face             string `json:"Face,omitempty"`
	BackgroundAccent string `json:"BackgroundAccent,omitempty"`
	Necklace         string `json:"Necklace,omitempty"`
	Head             string `json:"Head,omitempty"`
	Background       string `json:"Background,omitempty"`
	FacialHair       string `json:"FacialHair,omitempty"`
	BackHand         string `json:"BackHand,omitempty"`
	FrontHand        string `json:"FrontHand,omitempty"`
	Overcoat         string `json:"Overcoat,omitempty"`
	Pet              string `json:"Pet,omitempty"`
}

type Sale struct {
	Date   time.Time
	Tx     string
	Buyer  string
	Seller string
	Algo   string
	Fiat   float64
	Asset  uint64
}

func CreateAssetFromNote(note AlgoSeasNote, collectionName string, assetId string, updatedAt uint) Asset {
	idInt, _ := strconv.ParseUint(assetId, 10, 64)
	return Asset{
		ID:               idInt,
		UpdatedAt:        updatedAt,
		Collection:       "AlgoSeas Pirates",
		ImageUrl:         note.MediaURL,
		Combat:           uint64(note.Properties.Combat),
		Constitution:     uint64(note.Properties.Constitution),
		Luck:             uint64(note.Properties.Luck),
		Plunder:          uint64(note.Properties.Plunder),
		Scenery:          note.Properties.Scenery,
		LeftArm:          note.Properties.LeftArm,
		Body:             note.Properties.Body,
		BackItem:         note.Properties.BackItem,
		Pants:            note.Properties.Pants,
		Footwear:         note.Properties.Footwear,
		RightArm:         note.Properties.RightArm,
		Shirts:           note.Properties.Shirts,
		Hat:              note.Properties.Hat,
		HipItem:          note.Properties.HipItem,
		Tattoo:           note.Properties.Tattoo,
		Face:             note.Properties.Face,
		BackgroundAccent: note.Properties.BackgroundAccent,
		Necklace:         note.Properties.Necklace,
		Head:             note.Properties.Head,
		Background:       note.Properties.Background,
		FacialHair:       note.Properties.FacialHair,
		BackHand:         note.Properties.BackHand,
		FrontHand:        note.Properties.FrontHand,
		Overcoat:         note.Properties.Overcoat,
		Pet:              note.Properties.Pet,
	}
}

// func CreateSale()

func ParseSaleDate(timestamp string) time.Time {
	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, timestamp)

	if err != nil {
		fmt.Println(err)
	}
	return t
}

func InsertAsset(db *sql.DB, asset Asset) error {
	stmt, err := db.Prepare("REPLACE INTO asset VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		asset.ID,
		asset.UpdatedAt,
		asset.Collection,
		asset.ImageUrl,
		asset.Combat,
		asset.Constitution,
		asset.Luck,
		asset.Plunder,
		asset.Scenery,
		asset.LeftArm,
		asset.Body,
		asset.BackItem,
		asset.Pants,
		asset.Footwear,
		asset.RightArm,
		asset.Shirts,
		asset.HipItem,
		asset.Tattoo,
		asset.Face,
		asset.BackgroundAccent,
		asset.Necklace,
		asset.Hat,
		asset.Head,
		asset.Background,
		asset.FacialHair,
		asset.BackHand,
		asset.FrontHand,
		asset.Overcoat,
		asset.Pet,
	)
	stmt.Close()
	if err != nil {
		return err
	}
	return nil
}

func InsertSale(db *sql.DB, sale Sale) error {
	stmt, err := db.Prepare("INSERT INTO sale(Date,Tx, Buyer, Seller, Algo, Fiat, Asset) VALUES(?, ?, ?, ?, ?, ?, ?)")
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
	stmt.Close()
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
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `asset` (`ID` INT unsigned, `UpdatedAt` INT unsigned, `Collection` VARCHAR(255) NOT NULL,`ImageUrl` TEXT NOT NULL,`Combat` INT unsigned NOT NULL,`Constitution` INT unsigned NOT NULL,`Luck` INT unsigned NOT NULL,`Plunder` INT unsigned NOT NULL,`Scenery` VARCHAR(255) NOT NULL,`LeftArm` VARCHAR(255) NOT NULL,`Body` VARCHAR(255) NOT NULL,`BackItem` VARCHAR(255) NOT NULL,`Pants` VARCHAR(255) NOT NULL,`Footwear` VARCHAR(255) NOT NULL,`RightArm` VARCHAR(255) NOT NULL,`Shirts` VARCHAR(255) NOT NULL,`HipItem` VARCHAR(255) NOT NULL,`Tattoo` VARCHAR(255) NOT NULL,`Face` VARCHAR(255) NOT NULL,`BackgroundAccent` VARCHAR(255) NOT NULL,`Necklace` VARCHAR(255) NOT NULL,`Hat` VARCHAR(255) NOT NULL,`Head` VARCHAR(255) NOT NULL,`Background` VARCHAR(255) NOT NULL,`FacialHair` VARCHAR(255) NOT NULL,`BackHand` VARCHAR(255) NOT NULL,`FrontHand` VARCHAR(255) NOT NULL,`Overcoat` VARCHAR(255) NOT NULL,`Pet` VARCHAR(255) NOT NULL,KEY `Collection_ID_IDX` (`Collection`,`ID`) USING BTREE,PRIMARY KEY (`ID`));")
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

func GetLastAssetUpdate(db *sql.DB) uint {
	latestIngestedRound := uint(0)
	rows, _ := db.Query("SELECT UpdatedAt FROM asset ORDER BY UpdatedAt DESC LIMIT 1")
	for rows.Next() {
		err := rows.Scan(&latestIngestedRound)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return latestIngestedRound
}

func ReadAllAssets(db *sql.DB) []Asset {
	assets := []Asset{}
	rows, _ := db.Query("SELECT `ID`, `UpdatedAt`, `Collection`, `ImageUrl`, `Combat`, `Constitution`, `Luck`, `Plunder`, `Scenery`, `LeftArm`, `Body`, `BackItem`, `Pants`, `Footwear`, `RightArm`, `Shirts`, `HipItem`, `Tattoo`, `Face`, `BackgroundAccent`, `Necklace`, `Hat`, `Head`, `Background`, `FacialHair`, `BackHand`, `FrontHand`, `Overcoat`, `Pet` FROM asset")
	for rows.Next() {
		asset := Asset{}
		err := rows.Scan(&asset.ID, &asset.UpdatedAt, &asset.Collection, &asset.ImageUrl, &asset.Combat, &asset.Constitution, &asset.Luck, &asset.Plunder, &asset.Scenery, &asset.LeftArm, &asset.Body, &asset.BackItem, &asset.Pants, &asset.Footwear, &asset.RightArm, &asset.Shirts, &asset.HipItem, &asset.Tattoo, &asset.Face, &asset.BackgroundAccent, &asset.Necklace, &asset.Hat, &asset.Head, &asset.BackHand, &asset.FacialHair, &asset.BackHand, &asset.FrontHand, &asset.Overcoat, &asset.Pet)
		if err != nil {
			log.Fatalln(err)
		}
		assets = append(assets, asset)
	}
	return assets
}
