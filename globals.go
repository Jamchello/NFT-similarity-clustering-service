package main

import "database/sql"

var NumOfClusters = 5

// TODO: Alias the database, make the insertions etc a function of this...
var Database *sql.DB

// This file initializes the global maps used to store data within the program.
var IdToAsset = map[uint64]Asset{}
var IdToListings = map[string]AlgoSeasListingItem{}

var IdToCluster = map[uint64]int{}
var ClusterToAssetIds = make([][]uint64, NumOfClusters)
var ClusterToActiveAssetIds = make([][]uint64, NumOfClusters)

func UpdateAssetsMapping(assets []Asset) {
	for _, asset := range assets {
		IdToAsset[asset.ID] = asset
	}
}

func AssetIdsToAssets(assetIds []uint64) []Asset {
	assets := make([]Asset, len(assetIds))
	for i, assetId := range assetIds {
		assets[i] = IdToAsset[assetId]
	}
	return assets
}

func AssetIdsToListings(assetIds []uint64) []AlgoSeasListingItem {
	listings := make([]AlgoSeasListingItem, len(assetIds))
	for i, assetId := range assetIds {
		listings[i] = IdToListings[string(assetId)]
	}
	return listings
}
