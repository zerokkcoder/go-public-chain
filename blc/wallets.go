package blc

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "wallets_%s.dat"

type Wallets struct {
	Wallets map[string]*Wallet
}

// 创建钱包集合
func NewWallets(nodeID string) (*Wallets, error) {
	walletFile := fmt.Sprintf(walletFile, nodeID)

	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &Wallets{}
		wallets.Wallets = make(map[string]*Wallet)
		return wallets, err
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	return &wallets, nil
}

// 创建一个新钱包
func (w Wallets) CreateNewWallet(nodeID string) string {
	wallet := NewWallet()
	w.Wallets[string(wallet.GetAddress())] = wallet
	w.SaveToFile(nodeID)

	return string(wallet.GetAddress())
}

// 获取所有的钱包信息
func (w *Wallets) GetAddresses() []string {
	var addresses []string

	for address := range w.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

// 根据地址获取钱包对象
func (w *Wallets) GetWallet(address string) Wallet {
	return *w.Wallets[address]
}

// 加载钱包文件
func (w *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	w.Wallets = wallets.Wallets

	return nil
}

// 把 wallets 信息保存到一个文件中
func (w Wallets) SaveToFile(nodeID string) {

	walletFile := fmt.Sprintf(walletFile, nodeID)

	var content bytes.Buffer

	// 注册的目的，是为了，可以序列化任何类型
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&w)
	if err != nil {
		log.Panic(err)
	}

	// 将序列化以后的数据写入到文件，原来文件的数据会被覆盖
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
