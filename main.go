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

	// 查看数据
	err = db.View(func(tx *bolt.Tx) error {
		// 获取BlockBucket表
		b := tx.Bucket([]byte("BlockBucket"))

		// 往表里面存储数据
		if b != nil {
			data := b.Get([]byte("l"))
			fmt.Printf("%s\n", data)
			data = b.Get([]byte("lll"))
			fmt.Printf("%s\n", data)
		}

		// 返回nil，以便数据库处理相应操作
		return nil
	})
	// 更新失败
	if err != nil {
		log.Panic(err)
	}
}
