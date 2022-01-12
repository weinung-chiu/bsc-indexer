package portto

type Block struct {
	Number     uint64 `json:"block_num"`
	Hash       string `json:"block_hash"`
	Time       uint64 `json:"block_time"`
	ParentHash string `json:"parent_hash"`
}
