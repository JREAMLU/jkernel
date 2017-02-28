package entity

import "time"

type Response struct {
	Meta       MetaList    `json:"meta"`
	StatusCode int         `json:"status_code"`
	Message    interface{} `json:"message"`
	Data       interface{} `json:"data"`
}

type ResponseList struct {
	Meta       MetaList    `json:"meta"`
	StatusCode int         `json:"status_code"`
	Message    interface{} `json:"message"`
	Data       struct {
		Total int                    `json:"total"`
		List  map[string]interface{} `json:"list"`
	} `json:"data"`
}

type MetaList struct {
	RequestId string    `json:"Request-Id"`
	UpdatedAt time.Time `json:"updated_at"`
	Timezone  string    `json:"timezone"`
}
