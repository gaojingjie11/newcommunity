package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"smartcommunity/internal/global"
	"smartcommunity/internal/model"
	"smartcommunity/pkg/utils"
	"strings"
	"time"
)

type UserService struct{}

// Register 用户注册逻辑 [cite: 39]
func (s *UserService) Register(user *model.SysUser) error {
	// 1. 检查手机号是否已存在
	var count int64
	global.DB.Model(&model.SysUser{}).Where("mobile = ?", user.Mobile).Count(&count)
	if count > 0 {
		return errors.New("该手机号已注册")
	}

	// 2. 密码加密
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hash

	// 3. 设置默认值
	if user.Role == "" {
		user.Role = "user"
	}
	user.Balance = 100.00
	user.Status = 1
	// 确保 email 唯一性检查可以加在这里，暂时略过

	// 4. 写入数据库
	return global.DB.Create(user).Error
}

// Login 用户登录逻辑
func (s *UserService) Login(mobile, password, ip, userAgent string) (string, *model.SysUser, error) {
	var user model.SysUser
	// 1. 根据手机号查询用户
	if err := global.DB.Where("mobile = ?", mobile).First(&user).Error; err != nil {
		return "", nil, errors.New("账号不存在")
	}

	// 2. 验证密码
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", nil, errors.New("密码错误")
	}

	if user.Status == 0 {
		return "", nil, errors.New("账号已冻结")
	}

	//// 3. 生成 Token
	//token, err := utils.GenerateToken(user.ID, user.Role)
	//if err != nil {
	//	return "", nil, errors.New("Token生成失败")
	//}
	//
	//return token, &user, nil
	// [原有代码] 3. 生成 Token
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", nil, errors.New("Token生成失败")
	}

	// --- [新增 Redis 逻辑] ---
	// Key 格式: "login:token:{userId}"
	// Value: token 字符串
	// 过期时间: 与 JWT 过期时间保持一致 (24小时)
	redisKey := fmt.Sprintf("login:token:%d", user.ID)
	ctx := context.Background()

	// Set 第三个参数是过期时间，这里设为 24 小时
	err = global.RDB.Set(ctx, redisKey, token, 24*time.Hour).Err()
	if err != nil {
		// 如果 Redis 存入失败，建议报错，保证安全性
		return "", nil, errors.New("登录服务异常")
	}
	// -----------------------

	return token, &user, nil
}

// Logout 用户退出登录
func (s *UserService) Logout(userID int64) error {
	ctx := context.Background()
	redisKey := fmt.Sprintf("login:token:%d", userID)
	return global.RDB.Del(ctx, redisKey).Err()
}

// UpdateInfo 修改用户信息
// UpdateInfo 修改用户信息
func (s *UserService) UpdateInfo(userID int64, updates map[string]interface{}) error {
	// 只更新非空字段，由Controller层组装 updates map
	return global.DB.Model(&model.SysUser{}).Where("id = ?", userID).Updates(updates).Error
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID int64, oldPwd, newPwd string) error {
	var user model.SysUser
	if err := global.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}
	if !utils.CheckPasswordHash(oldPwd, user.Password) {
		return errors.New("旧密码错误")
	}
	hash, _ := utils.HashPassword(newPwd)
	return global.DB.Model(&user).Update("password", hash).Error
}

// ResetPassword 重置密码 (忘记密码)
func (s *UserService) ResetPassword(mobile, code, newPwd string) error {
	// 1. 校验验证码 (这里Mock一下，假设验证码是 "123456")
	if code != "123456" {
		return errors.New("验证码错误")
	}
	var user model.SysUser
	if err := global.DB.Where("mobile = ?", mobile).First(&user).Error; err != nil {
		return errors.New("该手机号未注册")
	}
	hash, _ := utils.HashPassword(newPwd)
	return global.DB.Model(&user).Update("password", hash).Error
}

// GetInfo 获取最新用户信息 (刷新页面用)
func (s *UserService) GetInfo(userID int64) (*model.SysUser, error) {
	var user model.SysUser
	err := global.DB.First(&user, userID).Error
	return &user, err
}

func (s *UserService) RegisterFace(userID int64, faceImageURL string) error {
	faceImageURL = strings.TrimSpace(faceImageURL)
	if faceImageURL == "" {
		return errors.New("face image url is required")
	}

	var user model.SysUser
	if err := global.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	return global.DB.Model(&user).Updates(map[string]interface{}{
		"face_registered": true,
		"face_image_url":  faceImageURL,
	}).Error
}

// SendSMSCode 发送验证码
func (s *UserService) SendSMSCode(mobile string) error {
	// 1. 生成6位随机验证码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06d", rnd.Intn(1000000))

	// 2. 存入 Redis, 有效期 5 分钟
	ctx := context.Background()
	key := fmt.Sprintf("sms:code:%s", mobile)
	if err := global.RDB.Set(ctx, key, code, 5*time.Minute).Err(); err != nil {
		return errors.New("系统繁忙，请稍后再试")
	}

	// 3. 调用 Spug Push API 发送
	// 模板 URL: https://push.spug.cc/send/nbONk8gz2Vr34gXG
	url := "https://push.spug.cc/send/nbONk8gz2Vr34gXG"

	payload := map[string]interface{}{
		"code":    code,
		"targets": mobile,
	}
	jsonBody, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return errors.New("短信发送失败: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("短信服务异常")
	}

	return nil
}

// LoginByCode 验证码登录
func (s *UserService) LoginByCode(mobile, code, ip, userAgent string) (string, *model.SysUser, error) {
	// 0. 处理空格
	code = strings.TrimSpace(code) // Assuming utils has Trim or just use strings.TrimSpace?
	// Let's use strings.TrimSpace, need to import strings

	// 1. 校验验证码
	ctx := context.Background()
	key := fmt.Sprintf("sms:code:%s", mobile)
	val, err := global.RDB.Get(ctx, key).Result()
	if err != nil {
		return "", nil, errors.New("验证码已失效或错误")
	}
	if val != code {
		return "", nil, errors.New("验证码错误")
	}

	// 验证成功后，立即删除验证码，防止重复使用
	global.RDB.Del(ctx, key)

	// 2. 查询用户，如果不存在则自动注册
	var user model.SysUser
	if err := global.DB.Where("mobile = ?", mobile).First(&user).Error; err != nil {
		// 自动注册逻辑
		user = model.SysUser{
			Mobile:   mobile,
			Username: mobile, // 默认用户名 = 手机号
			Role:     "user",
			Status:   1,
			Balance:  0.0,
			Avatar:   "https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png",
		}
		// 默认密码 123456
		hash, _ := utils.HashPassword("123456")
		user.Password = hash

		if err := global.DB.Create(&user).Error; err != nil {
			return "", nil, errors.New("自动注册失败: " + err.Error())
		}
	} else {
		// 如果用户存在，检查状态
		if user.Status == 0 {
			return "", nil, errors.New("账号已冻结")
		}
	}

	// 3. 生成 Token
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", nil, errors.New("Token生成失败")
	}

	// 4. 存入 Redis (同普通登录)
	tokenKey := fmt.Sprintf("login:token:%d", user.ID)
	err = global.RDB.Set(ctx, tokenKey, token, 24*time.Hour).Err()
	if err != nil {
		return "", nil, errors.New("登录服务异常")
	}

	return token, &user, nil
}
