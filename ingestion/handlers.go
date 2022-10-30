package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type SimilarAssetsResponse struct {
	SimilarAssets   []Asset               `json:"SimilarAssets"`
	RelatedListings []AlgoSeasListingItem `json:"RelatedListings"`
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
		cluster, ok := IdToCluster[uint64(assetId)]

		if !ok {
			http.Error(w, "No Listing stored for this asset", http.StatusBadRequest)
			return
		}

		// clusterAssets, ok := ClusterToAssetIds[cluster]

		// if !ok {
		// 	http.Error(w, "No Listing stored for this asset", http.StatusBadRequest)
		// 	return
		// }

		allInCluster := AssetIdsToAssets(ClusterToAssetIds[cluster])
		activeListingsInCluster := AssetIdsToListings(ClusterToActiveAssetIds[cluster])

		response := SimilarAssetsResponse{
			SimilarAssets:   allInCluster,
			RelatedListings: activeListingsInCluster,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// TODO: Return an assets data given its ID
func AssetHandler(w http.ResponseWriter, r *http.Request) {
}
