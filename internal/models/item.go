package models

import "errors"

type Item struct {
	ChrtId      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmId        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func (i *Item) Validate() error {
	if i.ChrtId <= 0 {
		return errors.New("chrt_id must be positive")
	}
	if i.TrackNumber == "" {
		return errors.New("track_number is empty")
	}
	if i.Price <= 0 {
		return errors.New("price must be positive")
	}
	if i.Rid == "" {
		return errors.New("rid is empty")
	}
	if i.Name == "" {
		return errors.New("name is empty")
	}
	if i.Sale < 0 {
		return errors.New("sale must be non-negative")
	}
	if i.TotalPrice <= 0 {
		return errors.New("total_price must be positive")
	}
	if i.NmId <= 0 {
		return errors.New("nm_id must be positive")
	}
	if i.Brand == "" {
		return errors.New("brand is empty")
	}
	if i.Status < 0 {
		return errors.New("status must be non-negative")
	}
	return nil
}
