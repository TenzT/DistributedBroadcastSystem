package gocache

import "testing"

func TestGoCacaheData_GetId(t *testing.T) {
	type fields struct {
		id        string
		raw       string
		signature string
		timestamp string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := GoCacaheData{
				id:        tt.fields.id,
				raw:       tt.fields.raw,
				signature: tt.fields.signature,
				timestamp: tt.fields.timestamp,
			}
			if got := d.GetId(); got != tt.want {
				t.Errorf("GoCacaheData.GetId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoCacaheData_GetRawData(t *testing.T) {
	type fields struct {
		id        string
		raw       string
		signature string
		timestamp string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := GoCacaheData{
				id:        tt.fields.id,
				raw:       tt.fields.raw,
				signature: tt.fields.signature,
				timestamp: tt.fields.timestamp,
			}
			if got := d.GetRawData(); got != tt.want {
				t.Errorf("GoCacaheData.GetRawData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoCacaheData_GetSignature(t *testing.T) {
	type fields struct {
		id        string
		raw       string
		signature string
		timestamp string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := GoCacaheData{
				id:        tt.fields.id,
				raw:       tt.fields.raw,
				signature: tt.fields.signature,
				timestamp: tt.fields.timestamp,
			}
			if got := d.GetSignature(); got != tt.want {
				t.Errorf("GoCacaheData.GetSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoCacaheData_GetTimeStamp(t *testing.T) {
	type fields struct {
		id        string
		raw       string
		signature string
		timestamp string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := GoCacaheData{
				id:        tt.fields.id,
				raw:       tt.fields.raw,
				signature: tt.fields.signature,
				timestamp: tt.fields.timestamp,
			}
			if got := d.GetTimeStamp(); got != tt.want {
				t.Errorf("GoCacaheData.GetTimeStamp() = %v, want %v", got, tt.want)
			}
		})
	}
}
