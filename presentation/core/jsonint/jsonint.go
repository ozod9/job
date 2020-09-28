package jsonint

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type JSONInt struct {
	Value int64
	Valid bool
	Set   bool
}

type JSONString struct {
	Value string
	Valid bool
	Set   bool
}

type TransactionJSON struct {
	FromId JSONInt    `json: "fromId"`
	ToId   JSONInt    `json: "toId"`
	Amount JSONString `json: "amount"`
	Reason JSONString `json: "reason"`
	Type   JSONString
}

type AllRatesJSON struct {
	Rates map[string]float64
}

func (i *JSONInt) UnmarshalJSON(data []byte) error {
	i.Set = true

	if string(data) == "null" {
		i.Valid = false
		return nil
	}

	var temp int64
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	i.Value = temp
	i.Valid = true
	return nil
}

func (s *JSONString) UnmarshalJSON(data []byte) error {
	s.Set = true
	if string(data) == "null" {
		s.Valid = false
		return nil
	}

	var temp string
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	s.Value = temp
	s.Valid = true
	return nil
}

func BodyToJSON(body io.ReadCloser, note interface{}) error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &note)
	return err
}
