package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "123456"

	// 测试当前数据库中的hash
	dbHash := "$2a$10$N.ZOn9G6/YLFixAOPMg/h.z7pCu6v2XyFDtC4q8NZpYYb2V7Yb0.C"

	fmt.Println("=== 测试旧hash (cost 10) ===")
	err := bcrypt.CompareHashAndPassword([]byte(dbHash), []byte(password))
	if err == nil {
		fmt.Println("✓ 旧hash验证成功")
	} else {
		fmt.Println("✗ 旧hash验证失败:", err)
	}

	// 生成新hash (cost 14)
	fmt.Println("\n=== 生成新hash (cost 14) ===")
	newHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Println("生成失败:", err)
		return
	}
	fmt.Println("新Hash:", string(newHash))

	// 验证新hash
	err = bcrypt.CompareHashAndPassword(newHash, []byte(password))
	if err == nil {
		fmt.Println("✓ 新hash验证成功")
	} else {
		fmt.Println("✗ 新hash验证失败:", err)
	}

	// 生成SQL更新语句
	fmt.Println("\n=== SQL更新语句 ===")
	fmt.Printf("UPDATE sys_user SET password = '%s' WHERE id IN (1,2,3,4,5);\n", string(newHash))
}
