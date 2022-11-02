package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type SimilarAssetsResponse struct {
	SimilarAssets   []Asset   `json:"SimilarAssets"`
	RelatedListings []Listing `json:"RelatedListings"`
}

func SimilarAssetsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		assetIdStr := r.URL.Query().Get("assetId")
		if assetIdStr == "" {
			http.Error(w, "Invalid assetId", http.StatusBadRequest)
			return
		}
		assetId, err := strconv.ParseUint(assetIdStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid assetId", http.StatusBadRequest)
			return
		}
		similar, ok := IdToSimilar[uint64(assetId)]

		if !ok {
			http.Error(w, "No Listing stored for this asset", http.StatusBadRequest)
			return
		}

		similarActive, ok := IdToSimilarListed[assetId]
		if !ok {
			http.Error(w, "No Listing stored for this asset", http.StatusBadRequest)
			return
		}
		// clusterAssets, ok := ClusterToAssetIds[cluster]

		// if !ok {
		// 	http.Error(w, "No Listing stored for this asset", http.StatusBadRequest)
		// 	return
		// }

		similarAssets := AssetIdsToAssets(similar)
		relatedListings := AssetIdsToListings(similarActive)
		relatedListingsFlat := []Listing{}
		for _, a := range relatedListings {
			relatedListingsFlat = append(relatedListingsFlat, a...)
		}
		// activeListingsInCluster := AssetIdsToListings(ClusterToActiveAssetIds[0])

		response := SimilarAssetsResponse{
			SimilarAssets:   similarAssets,
			RelatedListings: relatedListingsFlat,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func AssetHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		assetIdStr := r.URL.Query().Get("assetId")
		if assetIdStr == "" {
			http.Error(w, "Invalid assetId", http.StatusBadRequest)
			return
		}
		assetId, err := strconv.ParseUint(assetIdStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid assetId", http.StatusBadRequest)
			return
		}

		asset, ok := IdToAsset[assetId]
		if !ok {
			http.Error(w, "Asset does not exist / has no metadata", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(asset)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
