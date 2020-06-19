package util

type Table struct {
	Header Row   `json:"header,omitempty"`
	Rows   []Row `json:"row,omitempty"`
}

type Row []string
