package dbshttp

import "BroadcastService/common"

type ListAllRecentDataRsp struct {
	DataList []*HttpData `json:"data_list"`
}

type PostNewDataRsp struct {
	Result bool   `json:"result"`
	Msg    string `json:"msg"`
}

type HttpData struct {
	Id        string `json:"id"`
	Raw       string `json:"raw"`
	Signature string `json:"signature"`
}

func NewHttpData(data *common.Data) *HttpData {
	return &HttpData{
		Id:        data.Id,
		Raw:       data.Raw,
		Signature: data.Signature,
	}
}
