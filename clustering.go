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


func EuclideanDistance(asset1 Asset, asset2 Asset) float64 {
	combat:= math.Pow(float64(asset1.Combat-asset2.Combat),2)
	constitution:= math.Pow(float64(asset1.Constitution-asset2.Constitution),2)
	luck := math.Pow(float64(asset1.Luck-asset2.Luck),2)
	plunder:= math.Pow(float64(asset1.Plunder-asset2.Plunder),2)

	totalDistance:= math.Sqrt(combat + constitution + luck+ plunder)
	return totalDistance


}


func contains(list []uint64, element uint64) bool{
	for _, value := range(list){
		if(value == element) {
			return true
		}
	}
	return false
}

//Takes in an asset, and list of assets, then inserts Asset ID and list of similar assets to IdToAsset map
func findSimilarAssets(asset Asset) SimilarAssetsReturn{

	similarAssetIDs := []uint64{}
	similarAsset_distance := map[uint64]float64{}
	similarAsset_Listings :=[]uint64{}
	//iterate over all assets map
	for _, current_asset:= range IdToAsset{
		//when we encounter the asset we are checking against (as its a map) we ignore it
		if current_asset.ID != asset.ID{
			//if the similarAssetID list is <5 then we can just fill it
			if len(similarAssetIDs)<5{
				similarAssetIDs = append(similarAssetIDs, current_asset.ID)
				similarAsset_distance[current_asset.ID] = EuclideanDistance(asset, current_asset)
			//once it contains 5 ids we need to start checking for ones with smaller distances
			}else{
				//keep track of similarAsset with highest distance to our passed in asset (asset)
				highest_distance:= float64(0)
				var replace_index int
				lowerDistance := false
				for i:=0; i<5; i++{
					current_distance := EuclideanDistance(asset, IdToAsset[similarAssetIDs[i]]) //gets distance between passed in asset and similar asset
					if (current_distance > highest_distance){ //if the distance between them is the highest so far then we might want to replace it
						highest_distance = current_distance //replace highest distance with current distance
						new_distance :=EuclideanDistance(asset, current_asset) //gets distance between passed in asset and the fetched one from the hashmap
						//if the distance between the passed asset and the current one from the map is less than the previous current_distance, then we have found a potential replacement
						if (new_distance < current_distance){
							replace_index = i 
							lowerDistance = true
						}
					}
				}
				//if lowerDistance = true then we replaced one of the similar assets
				if lowerDistance == true {
					similarAssetIDs[replace_index] = current_asset.ID
					similarAsset_distance[current_asset.ID] = EuclideanDistance(asset, current_asset)
				//need to implement logic for finding similar asset Listings
			
				//could make similarAsset_distance a map of assetID -> distance from passed in asset 
				//we iterate through the similarAsset_Listing array and do similarAssetDistance[i] and see if its lower than that of the new asset
				//first of all, check if the new similarAsset has a listing (if it doesn't then we don't want to put it into similarListings anyway)
				_, ok := IdToListings[current_asset.ID]
				if ok == true{
					highestDistance:=float64(0)
					replace_index:= 0
					for index, assetID := range(similarAsset_Listings){
						currentAssetDistance := similarAsset_distance[assetID]  //distance from passed asset to the current iterated similarAsset
						if(currentAssetDistance > highestDistance){
							highestDistance = currentAssetDistance
							replace_index = index
						}
					}
						newAssetDistance := similarAsset_distance[current_asset.ID]  //distance from passed asset to the current iterated (top level for loop) asset
						//if the distance of our new asset is smaller then we can replace it 
						if(newAssetDistance < highestDistance){
							similarAsset_Listings[replace_index] = current_asset.ID
							delete(similarAsset_distance, similarAsset_Listings[replace_index])
						}

					}
				}

				
				

			}	
		}
	}
		//IdToSimilarAssets[asset.ID] = similarAssetIDs
		returnObject := SimilarAssetsReturn{similarAssetIDs, similarAsset_Listings}
		return returnObject

			
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
