package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

type Asset struct {
	ID               uint64 `db:"ID" json:"Id"`
	UpdatedAt        uint   `db:"UpdatedAt" json:"UpdatedAt,omitempty"`
	Collection       string `db:"Collection" json:"Collection"`
	ImageUrl         string `db:"ImageUrl" json:"ImageUrl"`
	Combat           uint64 `db:"Combat" json:"Combat"`
	Constitution     uint64 `db:"Constitution" json:"Constitution"`
	Luck             uint64 `db:"Luck" json:"Luck"`
	Plunder          uint64 `db:"Plunder" json:"Plunder"`
	Scenery          string `db:"Scenery" json:"Scenery,omitempty"`
	LeftArm          string `db:"LeftArm" json:"LeftArm,omitempty"`
	Body             string `db:"Body" json:"Body,omitempty"`
	BackItem         string `db:"BackItem" json:"BackItem,omitempty"`
	Pants            string `db:"Pants" json:"Pants,omitempty"`
	Footwear         string `db:"Footwear" json:"Footwear,omitempty"`
	RightArm         string `db:"RightArm" json:"RightArm,omitempty"`
	Shirts           string `db:"Shirts" json:"Shirts,omitempty"`
	Hat              string `db:"Hat" json:"Hat,omitempty"`
	HipItem          string `db:"HipItem" json:"HipItem,omitempty"`
	Tattoo           string `db:"Tattoo" json:"Tattoo,omitempty"`
	Face             string `db:"Face" json:"Face,omitempty"`
	BackgroundAccent string `db:"BackgroundAccent" json:"BackgroundAccent,omitempty"`
	Necklace         string `db:"Necklace" json:"Necklace,omitempty"`
	Head             string `db:"Head" json:"Head,omitempty"`
	Background       string `db:"Background" json:"Background,omitempty"`
	FacialHair       string `db:"FacialHair" json:"FacialHair,omitempty"`
	BackHand         string `db:"BackHand" json:"BackHand,omitempty"`
	FrontHand        string `db:"FrontHand" json:"FrontHand,omitempty"`
	Overcoat         string `db:"Overcoat" json:"Overcoat,omitempty"`
	Pet              string `db:"Pet" json:"Pet,omitempty"`
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

func ParseDate(timestamp string) time.Time {
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
