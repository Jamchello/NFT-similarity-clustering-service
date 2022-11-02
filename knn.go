package main

import (
	"github.com/kyroy/kdtree/points"
)

func VecoriseAsset(asset Asset) []float64 {
	return []float64{float64(asset.Combat), float64(asset.Constitution), float64(asset.Plunder), float64(asset.Luck)}
}

func VectoriseAssets(assets []Asset) [][]float64 {
	var asArray = make([][]float64, len(assets))
	for i, asset := range assets {
		asArray[i] = VecoriseAsset(asset)
	}
	return asArray
}

type Data struct {
	value uint64
}
type PointWithData struct {
	points.Point
	Data Data
}

func (d *Data) getValue() uint64 {
	return d.value
}
