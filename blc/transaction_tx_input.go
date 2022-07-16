package blc

type TXInput struct {
	TxHash    []byte // 1. 交易的Hash
	Vout      int    // 2. 存储 TXOutput 在 Vout里面的索引
	ScriptSig string // 3. 用户名
}

// 判断是否时所属地址的 TXInput
func (t *TXInput) UnLockWithAddress(address string) bool {
	return t.ScriptSig == address
}
