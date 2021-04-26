package main

import (
	"bytes"
	"fmt"
	"github.com/json-iterator/go"
	"log"
	"net/http"
)

const graph = "http://127.0.0.1:8547/graphql"

type GraphQLAccount struct {
	Address string `json:"address,omitzero"`
}

type GraphQLLog struct {
	Data    string         `json:"data,omitzero"`
	Topics  []string       `json:"topics,omitzero"`
	Account GraphQLAccount `json:"account,omitzero"`
}

type GraphQLTransaction struct {
	Logs []GraphQLLog `json:"logs,omitzero"`
}

// 以太坊number为无符整形，heco和bsc为字符串
type GraphQLBlockRsp struct {
	Number       int64                `json:"number,omitzero"`
	Timestamp    string               `json:"timestamp,omitzero"`
	Transactions []GraphQLTransaction `json:"transactions,omitzero"`
}

type GraphQLRsp struct {
	Blocks []GraphQLBlockRsp `json:"blocks,omitzero"`
}

type GraphRequest struct {
	Query string `json:"query,omitempty"`
}

type GraphResponse struct {
	Data  GraphQLRsp         `json:"data,omitempty"`
	Error GraphResponseError `json:"error,omitempty"`
}

type GraphResponseError struct {
	Code    int64  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func GraphQLRun(url, param string) (ret GraphQLRsp, err error) {
	req := GraphRequest{
		Query: param,
	}
	data, err := jsoniter.Marshal(req)
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(data)
	http_req, err := http.NewRequest(http.MethodPost, url, buf)
	http_req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(http_req)
	if nil != err {
		return
	}
	defer resp.Body.Close()

	var rsp GraphResponse
	if err = jsoniter.NewDecoder(resp.Body).Decode(&rsp); nil != err {
		return
	}

	return rsp.Data, nil
}

func main() {
	start := 6698820
	end := 6698830
	param := fmt.Sprintf(
		`{
		blocks(from: %d, to:%d) {
			number
			timestamp
			transactions {
				logs {
            		topics
            		data
					account{
			  			address
					}
				}
			}
		}
		}`, start, end)

	var rsp GraphQLRsp
	rsp, err := GraphQLRun(graph, param)
	if err != nil {
		log.Printf("GraphQLRun err %v\n", err)
		return
	}

	for _, block := range rsp.Blocks {
		for _, trans := range block.Transactions {
			for _, loger := range trans.Logs {

				log.Printf("number:%d,timestamp:%s,transactions:%d,data:%s,address:%s,topics:%s\n", block.Number, block.Timestamp, len(block.Transactions), loger.Data, loger.Account.Address, loger.Topics)
			}
		}

	}

}
