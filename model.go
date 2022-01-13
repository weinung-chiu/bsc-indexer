package portto

type Block struct {
	Number       uint64   `json:"block_num"`
	Hash         string   `json:"block_hash"`
	Time         uint64   `json:"block_time"`
	ParentHash   string   `json:"parent_hash"`
	Transactions []string `json:"transactions" gorm:"-"`
}

type Transaction struct {
	Hash      string `json:"tx_hash" gorm:"primaryKey"`
	From      string `json:"from" gorm:"column:from_addr"`
	To        string `json:"to" gorm:"column:to_addr"`
	Nonce     uint64 `json:"nonce"`
	Data      string `json:"data"`
	Value     uint64 `json:"value"`
	Logs      string `json:"logs"`
	BlockHash string `json:"-"`
}
