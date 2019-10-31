package model

import "encoding/json"

type Profile struct {
	Name         string
	Style        string
	City         string
	Dealer       string
	Timer        string
	Price        string
	Fuel         string
	Mileage      int
	Space        int
	Power        int
	Manipulation int
	Consumption  int
	Comfort      int
	Exterior     int
	Interior     int
	Cost         int
}

func FromJsonObj(o interface{}) (Profile, error) {
	var profile Profile
	s, err := json.Marshal(o)
	if err != nil {
		return profile, err
	}
	err = json.Unmarshal(s, &profile)
	return profile, err
}
