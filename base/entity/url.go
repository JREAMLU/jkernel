package entity

// URLShorten url shorten
type URLShorten struct {
	Meta struct {
		Auth string
	} `json:"meta" valid:"Required"`
	Data struct {
		URLs []interface{} `json:"urls" valid:"Required"`
		IP   string        `json:"ip" valid:"IP"`
	} `json:"data" valid:"Required"`
	FromIP string
}

// URLExpand shorten to expand struct
type URLExpand struct {
	Shorten []string `json:"shorten" valid:"Required"`
}
