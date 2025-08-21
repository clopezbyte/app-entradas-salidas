package models

type RmaRequest struct {
	Rma string `json:"rma"`
}

type RmaResponse struct {
	Cliente string `json:"cliente"`
}
