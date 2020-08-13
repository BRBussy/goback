package jsonrpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type rpcRequest struct {
	Method string `json:"method"`
}

func getRPCMethod(r *http.Request) (string, error) {
	// confirm body not nil
	if r.Body == nil {
		return "", errors.New("no body")
	}

	// read body bytes
	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(r.Body)

	// restore the io.ReadCloser to its original state for reading in
	// subsequent middleware
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// parse rpc request
	var request rpcRequest
	if err := json.Unmarshal(bodyBytes, &request); err != nil {
		return "", err
	}

	return request.Method, nil
}
