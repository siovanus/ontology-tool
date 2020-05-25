package side_chain_governance

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type jsonError struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    interface{}     `json:"data,omitempty"`
}

type blockReq struct {
	JsonRpc string               `json:"jsonrpc"`
	Method  string               `json:"method"`
	Params  []interface{}       `json:"params"`
	Id      uint                 `json:"id"`
}

type blockRsp struct {
	JsonRPC string               `json:"jsonrpc"`
	Result  *types.Header        `json:"result,omitempty"`
	Error   *jsonError           `json:"error,omitempty"`
	Id      uint                 `json:"id"`
}

type heightReq struct {
	JsonRpc string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  []string        `json:"params"`
	Id      uint            `json:"id"`
}

type heightRsp struct {
	JsonRpc string          `json:"jsonrpc"`
	Result  string          `json:"result,omitempty"`
	Error   *jsonError      `json:"error,omitempty"`
	Id      uint            `json:"id"`
}

func GetNodeHeader(restClient *RestClient, height uint64) ([]byte, error) {
	params := []interface{} {fmt.Sprintf("0x%x", height), true}
	req := &blockReq{
		JsonRpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  params,
		Id:      1,
	}
	reqdata, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("GetNodeHeight: marshal req err: %s", err)
	}
	rspdata, err := restClient.SendRestRequest(reqdata)
	if err != nil {
		return nil, fmt.Errorf("GetNodeHeight err: %s", err)
	}
	rsp := &blockRsp{}
	err = json.Unmarshal(rspdata, rsp)
	if err != nil {
		return nil, fmt.Errorf("GetNodeHeight, unmarshal resp err: %s", err)
	}
	if rsp.Error != nil {
		return nil, fmt.Errorf("GetNodeHeight, unmarshal resp err: %s", rsp.Error.Message)
	}
	block, err := json.Marshal(rsp.Result)
	return block, nil
}

func GetNodeHeight(restClient *RestClient) (uint64, error) {
	req := &heightReq{
		JsonRpc: "2.0",
		Method:  "eth_blockNumber",
		Params:  make([]string, 0),
		Id:      1,
	}
	reqData, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("GetNodeHeight: marshal req err: %s", err)
	}
	rspData, err := restClient.SendRestRequest(reqData)
	if err != nil {
		return 0, fmt.Errorf("GetNodeHeight err: %s", err)
	}
	rsp := &heightRsp{}
	err = json.Unmarshal(rspData, rsp)
	if err != nil {
		return 0, fmt.Errorf("GetNodeHeight, unmarshal resp err: %s", err)
	}
	if rsp.Error != nil {
		return 0, fmt.Errorf("GetNodeHeight, unmarshal resp err: %s", rsp.Error.Message)
	}
	height, err := strconv.ParseUint(rsp.Result, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("GetNodeHeight, parse resp height %s failed", rsp.Result)
	} else {
		return height, nil
	}
}

type RestClient struct {
	Addr       string
	restClient *http.Client
}

func NewRestClient() *RestClient {
	return &RestClient{
		restClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   5,
				DisableKeepAlives:     false,
				IdleConnTimeout:       time.Second * 300,
				ResponseHeaderTimeout: time.Second * 300,
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: time.Second * 300,
		},
	}
}

func (self *RestClient) SetAddr(addr string) *RestClient {
	self.Addr = addr
	return self
}

func (self *RestClient) SendRestRequest(data []byte) ([]byte, error) {
	resp, err := self.restClient.Post(self.Addr, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("http post request:%s error:%s", data, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read rest response body error:%s", err)
	}
	return body, nil
}
