package main

import (
	"github.com/kyroy/kdtree"
	"github.com/kyroy/kdtree/points"
)

func arrayifyAssets(assets []Asset) [][]float64 {
	var asArray = make([][]float64, len(assets))
	for i, asset := range assets {
		asArray[i] = []float64{float64(asset.Combat), float64(asset.Constitution), float64(asset.Plunder), float64(asset.Luck)}
	}
	return asArray
}

type Data struct {
	value uint64
}
type mypoint struct {
	points.Point
	Data Data
}

func (d *Data) getValue() uint64 {
	return d.value
}

func PerformKnnSearch(assets []Asset, mapping map[uint64][]uint64) {
	tree := kdtree.New([]kdtree.Point{})

	arrayified := arrayifyAssets(assets)

	for i, item := range arrayified {
		assetId := assets[i].ID
		pt := &mypoint{
			Point: *points.NewPoint(item, Data{value: assetId}),
			Data:  Data{value: assetId},
		}
		tree.Insert(pt)
	}

	for i, vector := range arrayified {
		closest := tree.KNN(&points.Point{Coordinates: vector, Data: Data{}}, 5)
		closestIds := []uint64{}
		for _, pt := range closest {
			closestIds = append(closestIds, pt.(*mypoint).Data.getValue())
		}
		assetId := assets[i].ID
		mapping[assetId] = closestIds
	}

}
