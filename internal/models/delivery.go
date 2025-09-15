package models

import "errors"

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

func (d *Delivery) Validate() error {
	if d.Name == "" {
		return errors.New("name is empty")
	}
	if d.Phone == "" {
		return errors.New("phone is empty")
	}
	if d.Zip == "" {
		return errors.New("zip is empty")
	}
	if d.City == "" {
		return errors.New("city is empty")
	}
	if d.Address == "" {
		return errors.New("address is empty")
	}
	if d.Region == "" {
		return errors.New("region is empty")
	}
	if d.Email == "" {
		return errors.New("email is empty")
	}
	return nil
}
