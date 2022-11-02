package main

import (
	"database/sql"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type IndexerSearchResp struct {
	CurrentRound int    `json:"current-round"`
	NextToken    string `json:"next-token"`
	Transactions []struct {
		ApplicationTransaction struct {
			Accounts          []interface{} `json:"accounts"`
			ApplicationArgs   []string      `json:"application-args"`
			ApplicationID     int           `json:"application-id"`
			ForeignApps       []interface{} `json:"foreign-apps"`
			ForeignAssets     []interface{} `json:"foreign-assets"`
			GlobalStateSchema struct {
				NumByteSlice int `json:"num-byte-slice"`
				NumUint      int `json:"num-uint"`
			} `json:"global-state-schema"`
			LocalStateSchema struct {
				NumByteSlice int `json:"num-byte-slice"`
				NumUint      int `json:"num-uint"`
			} `json:"local-state-schema"`
			OnCompletion string `json:"on-completion"`
		} `json:"application-transaction,omitempty"`
		CloseRewards     int    `json:"close-rewards"`
		ClosingAmount    int    `json:"closing-amount"`
		ConfirmedRound   int    `json:"confirmed-round"`
		Fee              int    `json:"fee"`
		FirstValid       int    `json:"first-valid"`
		GenesisHash      string `json:"genesis-hash"`
		GenesisID        string `json:"genesis-id"`
		GlobalStateDelta []struct {
			Key   string `json:"key"`
			Value struct {
				Action int `json:"action"`
				Uint   int `json:"uint"`
			} `json:"value"`
		} `json:"global-state-delta,omitempty"`
		Group     string `json:"group"`
		ID        string `json:"id"`
		InnerTxns []struct {
			AssetConfigTransaction struct {
				AssetID int `json:"asset-id"`
				Params  struct {
					Creator       string `json:"creator"`
					Decimals      int    `json:"decimals"`
					DefaultFrozen bool   `json:"default-frozen"`
					Manager       string `json:"manager"`
					Name          string `json:"name"`
					NameB64       string `json:"name-b64"`
					Reserve       string `json:"reserve"`
					Total         int    `json:"total"`
					UnitName      string `json:"unit-name"`
					UnitNameB64   string `json:"unit-name-b64"`
					URL           string `json:"url"`
					URLB64        string `json:"url-b64"`
				} `json:"params"`
			} `json:"asset-config-transaction"`
			CloseRewards      int    `json:"close-rewards"`
			ClosingAmount     int    `json:"closing-amount"`
			ConfirmedRound    int    `json:"confirmed-round"`
			CreatedAssetIndex int    `json:"created-asset-index"`
			Fee               int    `json:"fee"`
			FirstValid        int    `json:"first-valid"`
			IntraRoundOffset  int    `json:"intra-round-offset"`
			LastValid         int    `json:"last-valid"`
			ReceiverRewards   int    `json:"receiver-rewards"`
			RoundTime         int    `json:"round-time"`
			Sender            string `json:"sender"`
			SenderRewards     int    `json:"sender-rewards"`
			TxType            string `json:"tx-type"`
		} `json:"inner-txns,omitempty"`
		IntraRoundOffset int    `json:"intra-round-offset"`
		LastValid        int    `json:"last-valid"`
		ReceiverRewards  int    `json:"receiver-rewards"`
		RoundTime        int    `json:"round-time"`
		Sender           string `json:"sender"`
		SenderRewards    int    `json:"sender-rewards"`
		Signature        struct {
			Sig string `json:"sig"`
		} `json:"signature"`
		TxType                 string `json:"tx-type"`
		AssetConfigTransaction struct {
			AssetID int `json:"asset-id"`
			Params  struct {
				Creator       string `json:"creator"`
				Decimals      int    `json:"decimals"`
				DefaultFrozen bool   `json:"default-frozen"`
				Manager       string `json:"manager"`
				Reserve       string `json:"reserve"`
				Total         int    `json:"total"`
			} `json:"params"`
		} `json:"asset-config-transaction,omitempty"`
		Note string `json:"note,omitempty"`
	} `json:"transactions"`
}

type Transaction struct {
	ID                     string `json:"id"`
	RoundTime              int    `json:"round-time"`
	Sender                 string `json:"sender"`
	Note                   string `json:"note"`
	TxType                 string `json:"tx-type"`
	ConfirmedRound         uint   `json:"confirmed-round"`
	AssetConfigTransaction struct {
		AssetId int `json:"asset-id"`
	} `json:"asset-config-transaction"`
}

type TxHistoryResponse struct {
	NextToken    string        `json:"next-token"`
	Transactions []Transaction `json:"transactions"`
}

func GetMetadataForAsset(assetId string) AlgoSeasNote {
	note := AlgoSeasNote{}
	indexerUrl := fmt.Sprintf("https://algoindexer.algoexplorerapi.io/v2/transactions?asset-id=%s&tx-type=acfg", assetId)
	res, err := http.Get(indexerUrl)
	if err != nil {
		fmt.Printf("Retrying fetch %s in 3 seconds\n", assetId)
		time.Sleep(3 * time.Second)
		return GetMetadataForAsset(assetId)
	} else {
		defer res.Body.Close()
		parsed := IndexerSearchResp{}
		json.NewDecoder(res.Body).Decode(&parsed)
		transactions := parsed.Transactions
		// sort.Slice(transactions, func(i, j int) bool {
		// 	return transactions[i].RoundTime > transactions[j].RoundTime
		// })
		for _, tx := range transactions {
			if tx.Note != "" {
				noteBytes, _ := b64.StdEncoding.DecodeString(tx.Note)
				err := json.Unmarshal(noteBytes, &note)
				if err != nil {
					return note
				}
				if note.Standard == "arc69" {
					return note
				}
			}
		}
	}
	return note
}

//	func GetMintedAssetIds() []string {
//		mintedAssets := []string{}
//
// TODO: consolidate into single func
func GetNewMetadata(db *sql.DB) []Asset {
	latestUpdate := GetLastAssetUpdate(db)
	seenAssets := map[int]bool{}
	assets := []Asset{}
	hasMore := true
	nextToken := ""

	for hasMore {
		indexerUrl := fmt.Sprintf("https://algoindexer.algoexplorerapi.io/v2/transactions?address=SEASZVO4B4DC3F2SQKQVTQ5WXNVQWMCIPFPWTNQT3KMUX2JEGJ5K76ZC4Q&address-role=sender&tx-type=acfg&next=%s&min-round=%d", nextToken, latestUpdate+1)
		res, err := http.Get(indexerUrl)
		if err != nil {
			fmt.Println("Retrying fetch transactions in 3 seconds")
			time.Sleep(3 * time.Second)
		} else {
			defer res.Body.Close()
			parsed := TxHistoryResponse{}
			json.NewDecoder(res.Body).Decode(&parsed)
			for _, tx := range parsed.Transactions {
				assetId := tx.AssetConfigTransaction.AssetId
				if !seenAssets[assetId] && tx.Note != "" {
					note := AlgoSeasNote{}
					noteBytes, _ := b64.StdEncoding.DecodeString(tx.Note)
					err := json.Unmarshal(noteBytes, &note)
					if err != nil {
						break
					}
					if note.Standard == "arc69" {
						asset := CreateAssetFromNote(note, "AlgoSeas Pirates", strconv.Itoa(assetId), tx.ConfirmedRound)
						assets = append(assets, asset)
						seenAssets[assetId] = true
					}

				}
			}
			if parsed.NextToken == "" {
				hasMore = false
			}
			nextToken = parsed.NextToken
		}
	}

	return assets

}

func GetAllMintedAssets() []Asset {
	seenAssets := map[int]bool{}
	assets := []Asset{}
	hasMore := true
	nextToken := ""
	i := 0
	for hasMore {
		i++
		if i%10 == 0 {
			fmt.Printf("Making request number %d to AlgoIndexer, please be patient whilst we scrape all existing NFT's\n", i)
		}

		indexerUrl := fmt.Sprintf("https://algoindexer.algoexplorerapi.io/v2/transactions?address=SEASZVO4B4DC3F2SQKQVTQ5WXNVQWMCIPFPWTNQT3KMUX2JEGJ5K76ZC4Q&address-role=sender&tx-type=acfg&next=%s", nextToken)
		res, err := http.Get(indexerUrl)
		if err != nil {
			fmt.Println("Retrying fetch transactions in 3 seconds")
			time.Sleep(3 * time.Second)
		} else {
			defer res.Body.Close()
			parsed := TxHistoryResponse{}
			json.NewDecoder(res.Body).Decode(&parsed)
			for _, tx := range parsed.Transactions {
				assetId := tx.AssetConfigTransaction.AssetId
				if !seenAssets[assetId] && tx.Note != "" {
					note := AlgoSeasNote{}
					noteBytes, _ := b64.StdEncoding.DecodeString(tx.Note)
					err := json.Unmarshal(noteBytes, &note)
					if err != nil {
						break
					}
					if note.Standard == "arc69" {
						asset := CreateAssetFromNote(note, "AlgoSeas Pirates", strconv.Itoa(assetId), tx.ConfirmedRound)
						assets = append(assets, asset)
						seenAssets[assetId] = true
					}

				}
			}
			if parsed.NextToken == "" {
				hasMore = false
			}
			nextToken = parsed.NextToken
		}
	}

	return assets
}
