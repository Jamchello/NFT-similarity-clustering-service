package main

import (
	"fmt"
	"math"

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
	// Calculating the number of clusters dynamically
	numOfClusters := int(math.Sqrt(float64(len(assetList)) / 2))
	// re-assigning the 2d arrays which map cluster => Asset Ids
	ClusterToActiveAssetIds = make([][]uint64, numOfClusters)
	ClusterToAssetIds = make([][]uint64, numOfClusters)
	data := arrayifyAssets(assetList)

	c, e := clusters.KMeans(50, numOfClusters, clusters.EuclideanDistance)
	if e != nil {
		fmt.Println("Error in clustering", e)
		return
	}

	// Training
	if e = c.Learn(data); e != nil {
		fmt.Println("Error in clustering", e)
		return
	}

	//Resetting existing mappings of Cluster:Assets
	for i := 0; i < numOfClusters; i++ {
		ClusterToAssetIds[i] = []uint64{}
		ClusterToActiveAssetIds[i] = []uint64{}
	}

	for index, clusterNumber := range c.Guesses() {
		clusterIndex := clusterNumber - 1
		asset := assetList[index]
		//insert into hashmap Asset ID as key and then the cluster number for the value
		IdToCluster[asset.ID] = clusterIndex

		//insert cluster number as key and value as the given array with givenAsset appended to it
		ClusterToAssetIds[clusterIndex] = append(ClusterToAssetIds[clusterIndex], asset.ID)
		_, isActive := IdToListings[asset.ID]
		if isActive {
			ClusterToActiveAssetIds[clusterIndex] = append(ClusterToActiveAssetIds[clusterIndex], asset.ID)
		}
	}
}
