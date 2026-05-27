package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "123456"

	// 使用 cost 14 生成 hash（与后端一致）
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Println("生成hash失败:", err)
		return
	}

	fmt.Println("密码:", password)
	fmt.Println("Hash (cost 14):", string(hash))

	// 验证
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == nil {
		fmt.Println("✓ 验证成功")
	} else {
		fmt.Println("✗ 验证失败")
	}
}
