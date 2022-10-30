package main

// This file initializes the global maps used to store data within the program.
var IdToAsset = map[uint64]Asset{}
var IdToListings = map[string]AlgoSeasListingItem{}

var IdToCluster = map[uint64]int{}
var ClusterToAssetIds = map[int][]uint64{}
var ClusterToActiveAssetIds = map[int][]uint64{}

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
