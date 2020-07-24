package model

type Response struct {
	Data Data
}

type Data struct {
	Todo Todo `json:"todo"`
}
