package gorm

type Cago struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Length   uint64 `json:"length"`
	MaxAge   uint64 `json:"max_age"`
	CreateAt uint64 `json:"create_at"`
	UpdateAt uint64 `json:"update_at"`
}

func (c *Cago) TableName() string { return "cagos" }
