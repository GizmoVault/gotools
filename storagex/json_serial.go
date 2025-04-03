package storagex

import "encoding/json"

type JSONSerial struct {
	MarshalIndent bool
}

func (serial *JSONSerial) Marshal(t any) ([]byte, error) {
	if serial.MarshalIndent {
		return json.MarshalIndent(t, "", "\t")
	}

	return json.Marshal(t)
}

func (*JSONSerial) Unmarshal(d []byte, t any) error {
	return json.Unmarshal(d, t)
}
