package jsonsock

import (
	"encoding/json"
	"errors"
)

type JsonRequest struct {
	Id      int           `json:"id"`
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type JsonResponse struct {
	Id      int         `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
}

func MarshalRequest(targetFunc string, funcArgs []interface{}) (string, error) {
	req := &JsonRequest{123, "2.0", targetFunc, funcArgs}
	resp, err := json.Marshal(req)
	if err != nil {
		return "", errors.New("Failed to convert!")
	}
	return string(resp), nil
}

func UnmarshalRequest(incomingRequest string) (JsonRequest, error) {
	var request JsonRequest
	err := json.Unmarshal([]byte(incomingRequest), &request)
	if err != nil {
		return request, errors.New("Failed to unmarshal request")
	}
	return request, nil
}

func MarshalResponse(response interface{}) (string, error) {
	jresp := &JsonResponse{123, "2.0", response}
	resp, err := json.Marshal(jresp)
	if err != nil {
		return "", errors.New("Failed to create response")
	}
	return string(resp), nil
}
