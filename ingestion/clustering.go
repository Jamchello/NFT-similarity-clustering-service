package main

import (
	"fmt"

	"github.com/mpraski/clusters"
)

func arrayifyAssets(assets []Asset) [][]float64 {
	var asArray = make([][]float64, len(assets))
	for i, asset := range assets {
		asArray[i] = []float64{float64(asset.Combat), float64(asset.Constitution), float64(asset.Plunder), float64(asset.Luck)}
	}
	return asArray
}

func PerformClustering(assetList []Asset) {
	data := arrayifyAssets(assetList)
	numOfClusters := 5

	c, e := clusters.KMeans(20, numOfClusters, clusters.EuclideanDistance)
	if e != nil {
		fmt.Println("Error in clustering", e)
		return
	}

	// Use the data to train the clusterer
	if e = c.Learn(data); e != nil {
		fmt.Println("Error in clustering", e)
		return
	}

	fmt.Printf("Clustered data set into %d\n", c.Sizes())

	//Clearing out existing assets
	for i := 1; i <= numOfClusters; i++ {
		ClusterToAssetIds[i] = []uint64{}
		ClusterToActiveAssetIds[i] = []uint64{}
	}

	for index, number := range c.Guesses() {
		assetCluster := number
		asset := assetList[index]
		//insert into hashmap Asset ID as key and then the cluster number for the value
		IdToCluster[asset.ID] = assetCluster

		//insert cluster number as key and value as the given array with givenAsset appended to it
		ClusterToAssetIds[number] = append(ClusterToAssetIds[number], asset.ID)
		_, isActive := IdToListings[string(asset.ID)]
		if isActive {
			ClusterToActiveAssetIds[number] = append(ClusterToActiveAssetIds[number], asset.ID)
		}
	}

	fmt.Println(len(ClusterToAssetIds))

	for k := range ClusterToAssetIds {
		fmt.Println(k)
	}

}
