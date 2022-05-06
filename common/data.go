package common

type Data struct {
	Id        string
	Raw       string
	Signature string
	TimeStamp string
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
