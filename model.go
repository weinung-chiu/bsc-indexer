package portto

type Block struct {
	Number       uint64   `json:"block_num"`
	Hash         string   `json:"block_hash"`
	Time         uint64   `json:"block_time"`
	ParentHash   string   `json:"parent_hash"`
	Transactions []string `json:"transactions" gorm:"-"`
}

type Transaction struct {
	Hash      string
	From      string `gorm:"column:from_addr"`
	To        string `gorm:"column:to_addr"`
	Nonce     uint64
	Data      string
	Value     uint64
	Logs      string
	BlockHash string
}
