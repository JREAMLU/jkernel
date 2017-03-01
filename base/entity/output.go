package entity

// DataList output struct
type DataList struct {
	Total int                    `json:"total"`
	List  map[string]interface{} `json:"list"`
}

// NewDataList return datalist
func NewDataList() *DataList {
	return &DataList{}
}
