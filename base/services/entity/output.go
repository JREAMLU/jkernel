package entity

type DataList struct {
	Total int                    `json:"total"`
	List  map[string]interface{} `json:"list"`
}
