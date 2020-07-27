package model

// Response represents the full gql response object
type Response struct {
	Data Data
}

// Data represents the data block of the gql response object
type Data struct {
	Todo Todo `json:"todo"`
}
