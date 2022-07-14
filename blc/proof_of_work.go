package blc

type ProofOfWork struct {
	Block *Block // 验证的区块
}

func NewProofOfWork(block *Block) *ProofOfWork {
	return &ProofOfWork{block}
}

func (pow ProofOfWork) Run() ([]byte, int64) {
	return nil, 0
}
