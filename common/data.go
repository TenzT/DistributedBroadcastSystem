package common

type Data struct {
	Id        string `json:"id"`
	Raw       string `json:"raw""`
	Signature string `json:"signature"`
	TimeStamp string `json:"time_stamp"` // unix timestamp
}

func (d *Data) GetId() string {
	return d.Id
}

func (d *Data) GetRawData() string {
	return d.Raw
}

func (d *Data) GetSignature() string {
	return d.Signature
}

func (d *Data) GetTimeStamp() string {
	return d.TimeStamp
}
