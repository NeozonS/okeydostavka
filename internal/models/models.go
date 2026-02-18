package models

type City struct {
	StoreID string `json:"storeId"`
	FfcID   string `json:"ffcId"`
	City    string `json:"city"`

	Stores []Store `json:"storeID,omitempty"`
}
type Store struct {
	FfcID string `json:"ffcId"`
}

type Category struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Subcategories []Category `json:"subcategories,omitempty"`
}

type Product struct {
	Name       string `json:"name"`
	OfferPrice string `json:"offerPrice"`
	ProductURL string `json:"productUrl"`
}

type Stores struct {
	Stores []Store `json:"stores"`
}
type Regions struct {
	Regions []City `json:"regions"`
}
type Categories struct {
	Categories []Category `json:"categories"`
}

type Products struct {
	Products   []Product `json:"products"`
	TotalCount int       `json:"totalCount"`
	Offset     int       `json:"offset"`
	Name       string    `json:"name"`
}
type AddressResponse struct {
	Addresses []struct {
		FfcID int64 `json:"ffcId"`
	} `json:"addresses"`
}

type CategoryProducts struct {
	CategoryID   string    `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Products     []Product `json:"products"`
}
