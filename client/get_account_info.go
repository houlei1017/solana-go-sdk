package client

import (
	"context"
	"encoding/json"
	"errors"
)

type GetAccountInfoConfigEncoding string

const (
	GetAccountInfoConfigEncodingBase58     GetAccountInfoConfigEncoding = "base58" // limited to Account data of less than 128 bytes
	GetAccountInfoConfigEncodingBase64     GetAccountInfoConfigEncoding = "base64"
	GetAccountInfoConfigEncodingBase64Zstd GetAccountInfoConfigEncoding = "base64+zstd"
)

type GetAccountInfoConfig struct {
	Encoding  GetAccountInfoConfigEncoding
	DataSlice GetAccountInfoConfigDataSlice
}

type getAccountInfo struct {
	Encoding  GetAccountInfoConfigEncoding   `json:"encoding"`
	DataSlice *GetAccountInfoConfigDataSlice `json:"dataSlice,omitempty"`
}

func (cfg GetAccountInfoConfig) MarshalJSON() ([]byte, error) {
	var dataSlice *GetAccountInfoConfigDataSlice = nil
	if cfg.DataSlice != (GetAccountInfoConfigDataSlice{}) {
		dataSlice = &cfg.DataSlice
	}
	return json.Marshal(getAccountInfo{
		Encoding:  cfg.Encoding,
		DataSlice: dataSlice,
	})
}

type GetAccountInfoConfigDataSlice struct {
	Offset uint64 `json:"offset"`
	Length uint64 `json:"length"`
}

type GetAccountInfoResponse struct {
	Lamports  uint64      `json:"lamports"`
	Owner     string      `json:"owner"`
	Excutable bool        `json:"excutable"`
	RentEpoch uint64      `json:"rentEpoch"`
	Data      interface{} `json:"data"`
}

//type GetAccountInfoResponse struct {
//	Data struct {
//		Parsed struct {
//			Info struct {
//				IsNative    bool   `json:"isNative"`
//				Mint        string `json:"mint"`
//				Owner       string `json:"owner"`
//				State       string `json:"state"`
//				TokenAmount struct {
//					Amount         string  `json:"amount"`
//					Decimals       int     `json:"decimals"`
//					UiAmount       float64 `json:"uiAmount"`
//					UiAmountString string  `json:"uiAmountString"`
//				} `json:"tokenAmount"`
//			} `json:"info"`
//			Type string `json:"type"`
//		} `json:"parsed"`
//		Program string `json:"program"`
//		Space   int    `json:"space"`
//	} `json:"data"`
//	Executable bool   `json:"executable"`
//	Lamports   int    `json:"lamports"`
//	Owner      string `json:"owner"`
//	RentEpoch  int    `json:"rentEpoch"`
//}

func (s *Client) GetAccountInfo(ctx context.Context, account string, cfg GetAccountInfoConfig) (*GetAccountInfoResponse, error) {
	res := struct {
		GeneralResponse
		Result struct {
			Context Context                `json:"context"`
			Value   GetAccountInfoResponse `json:"value"`
		} `json:"result"`
	}{}
	err := s.request(ctx, "getAccountInfo", []interface{}{account, cfg}, &res)
	if err != nil {
		return &GetAccountInfoResponse{}, err
	}
	if res.Error != (ErrorResponse{}) {
		return &GetAccountInfoResponse{}, errors.New(res.Error.Message)
	}
	return &res.Result.Value, nil
}
