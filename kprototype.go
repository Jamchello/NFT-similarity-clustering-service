package main

import (
	"fmt"
	"math"

	"gopkg.in/e-XpertSolutions/go-cluster.v1/cluster"
)

//takes in an asset and returns it in vectorized form.
//uses the CategoriesDict dictionary to get encoded values for each piece of clothing

func KPVectorizeAsset(asset Asset) []float64 {
	return []float64{float64(asset.Combat), float64(asset.Constitution), float64(asset.Plunder), float64(asset.Luck),
	float64(CategoriesDict[asset.Scenery]),	float64(CategoriesDict[asset.LeftArm]),
	float64(CategoriesDict[asset.Body])	,float64(CategoriesDict[asset.BackItem]),
	float64(CategoriesDict[asset.Pants]),	float64(CategoriesDict[asset.Footwear]),
	float64(CategoriesDict[asset.RightArm]),	float64(CategoriesDict[asset.Shirts]),
	float64(CategoriesDict[asset.Hat]), 	float64(CategoriesDict[asset.HipItem]),
	float64(CategoriesDict[asset.Tattoo]), 	float64(CategoriesDict[asset.Face]),
	float64(CategoriesDict[asset.BackgroundAccent]), 	float64(CategoriesDict[asset.Necklace]),
	float64(CategoriesDict[asset.Head]), 	float64(CategoriesDict[asset.Background]),
	float64(CategoriesDict[asset.FacialHair]), 	float64(CategoriesDict[asset.BackHand]),
	float64(CategoriesDict[asset.FrontHand]), 	float64(CategoriesDict[asset.Overcoat]),
	float64(CategoriesDict[asset.Pet])}

}



func KPTestVecorizeAssets(assets []Asset) []float64{
	vector_assets := make([]float64, len(assets)*25)
	for _, asset :=range(assets){
		vector_assets = append(vector_assets, KPVectorizeAsset(asset)...)
	}

	return vector_assets
}


func KPvectorizeAssets(assets []Asset) [][]float64{
	vector_assets := make([][]float64, len(assets))
	for _, asset := range(assets){
		vector_assets = append(vector_assets, KPVectorizeAsset(asset))
	}

	return vector_assets

}




func setup(asset Asset, assets []Asset){
	vectorized_asset := VecoriseAsset(asset)
	VectorizedAssets := KPTestVecorizeAssets(assets)

	data := cluster.NewDenseMatrix(len(assets), 25, VectorizedAssets)
	newData := cluster.NewDenseMatrix(1, 25, vectorized_asset)

	
	clusters := math.Sqrt(float64(len(assets)/2))
	max_iterations := 20 

	distanceFunction := cluster.WeightedHammingDistance
	initializationFunction := cluster.InitCao

	//we weight the first 4 (numeric stats) thrice as much as we do clothing
	weights := []float64{3,3,3,3,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1}
	wvec := [][]float64{weights}

	//5-25 of fields we are assessing
	categorical_columns := []int{5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25}
	gamma := 0.2

	//initalise the KPrototype object
	kp := cluster.NewKPrototypes(distanceFunction, initializationFunction, categorical_columns, int(clusters), 1, max_iterations, wvec, gamma, "")

	//now to train


	err := kp.FitModel(data)
	if err != nil{
		fmt.Println(err)
	}

	//predict for the new data
	newLabelsP, err := kp.Predict(newData)
	if err!= nil{
		fmt.Println(err)
	}

	fmt.Println(newLabelsP)
	

}