package main

import (
	"github.com/kyroy/kdtree"
	"github.com/kyroy/kdtree/points"
)

// This file initializes the global maps used to store data within the program.
var IdToAsset = map[uint64]Asset{}
var IdToListings = map[uint64][]Listing{}

var MainKdTree = kdtree.New([]kdtree.Point{})
var SecondaryKdTree *kdtree.KDTree
var IdToVector = map[uint64][]float64{}

var CategoriesDict = map[string]float64{}

func AddToMainKdTree(asset Asset) {
	vector := VecoriseAsset(asset)
	IdToVector[asset.ID] = vector
	pt := &PointWithData{
		Point: *points.NewPoint(vector, Data{value: asset.ID}),
		Data:  Data{value: asset.ID},
	}
	MainKdTree.Insert(pt)
}

func BuildListingsKdTree(assets []Asset) {
	SecondaryKdTree = kdtree.New([]kdtree.Point{})
	vectorised := VectoriseAssets(assets)
	for i, item := range vectorised {
		assetId := assets[i].ID
		pt := &PointWithData{
			Point: *points.NewPoint(item, Data{value: assetId}),
			Data:  Data{value: assetId},
		}
		SecondaryKdTree.Insert(pt)
	}
}

func RemoveFromMainKdTree(asset Asset) {
	vector := IdToVector[asset.ID]
	MainKdTree.Remove(points.NewPoint(vector, nil))
}

func GetNMostSimilarIds(assetId uint64, n int) []uint64 {
	vector := IdToVector[assetId]
	closest := MainKdTree.KNN(&points.Point{Coordinates: vector, Data: Data{}}, n)
	closestIds := []uint64{}
	for _, pt := range closest {
		closestIds = append(closestIds, pt.(*PointWithData).Data.getValue())
	}
	return closestIds
}

func GetNMostSimilarListedIds(assetId uint64, n int) []uint64 {
	vector := IdToVector[assetId]
	closest := SecondaryKdTree.KNN(&points.Point{Coordinates: vector, Data: Data{}}, n)
	closestIds := []uint64{}
	for _, pt := range closest {
		closestIds = append(closestIds, pt.(*PointWithData).Data.getValue())
	}
	return closestIds
}

func AssetIdsToAssets(assetIds []uint64) []Asset {
	assets := make([]Asset, len(assetIds))
	for i, assetId := range assetIds {
		assets[i] = IdToAsset[assetId]
	}
	return assets
}

func AssetIdsToListings(assetIds []uint64) [][]Listing {
	listings := make([][]Listing, len(assetIds))
	for i, assetId := range assetIds {
		listings[i] = IdToListings[assetId]
	}
	return listings
}
