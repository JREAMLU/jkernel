package entity

// IPInfo ip info struct
type IPInfo struct {
	IPs []string `json:"ips" valid:"Required"`
}
