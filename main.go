package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

func main() {
	// 创建或者打开数据库
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// 创建表
	err = db.Update(func(tx *bolt.Tx) error {
		// 创建BlockBucket表
		b, err := tx.CreateBucket([]byte("BlockBucket"))
		if err != nil {
			return fmt.Errorf("create bucker: %s", err)
		}

		// 往表里面存储数据
		if b != nil {
			err := b.Put([]byte("l"), []byte("send 100 btc to shanshan...."))
			if err != nil {
				log.Panic("数据存储失败......")
			}
		}

		// 返回nil，以便数据库处理相应操作
		return nil
	})
	// 更新失败
	if err != nil {
		log.Panic(err)
	}
}
