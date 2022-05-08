package gocache

type GoCacaheData struct {
	id        string
	raw       string
	signature string
	timestamp string // timestamp in unix
}

func (d GoCacaheData) GetId() string {
	return d.id
}

func (d GoCacaheData) GetRawData() string {
	return d.raw
}

func (d GoCacaheData) GetSignature() string {
	return d.signature
}

func (d GoCacaheData) GetTimeStamp() string {
	return d.timestamp
}
