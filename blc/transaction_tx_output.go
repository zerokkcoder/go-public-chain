package blc

type TXOutput struct {
	Value        int64
	ScriptPubKey string // 用户名
}

// 判断是否时所属地址的 TXOutput
func (t *TXOutput) UnLockScriptPubKeyWithAddress(address string) bool {
	return t.ScriptPubKey == address
}
