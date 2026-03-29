package controller

import (
	"smartcommunity/internal/model"
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service service.UserService
}

// 这里必须要用 json:"password"，否则读不到前端传的值
type RegisterRequest struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	RealName string `json:"real_name"`
	Age      int    `json:"age"`
	Gender   int    `json:"gender"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"` // 新增
}

// Register 处理注册请求
func (h *UserHandler) Register(c *gin.Context) {

	// 使用专门的请求结构体来接收参数
	var req RegisterRequest

	// 1. 绑定参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数解析失败: "+err.Error())
		return
	}

	// 2. 基础校验
	if req.Mobile == "" || req.Password == "" {
		response.Fail(c, "手机号和密码不能为空")
		return
	}

	// 3. 将请求结构体(DTO) 转换为 数据库模型(Model)
	// Default Avatar
	defaultAvatar := "https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png"
	if req.Avatar == "" {
		req.Avatar = defaultAvatar
	}

	user := model.SysUser{
		Mobile:   req.Mobile,
		Password: req.Password, // 这里接收到了明文密码，Service层会加密
		RealName: req.RealName,
		Age:      req.Age,
		Gender:   req.Gender,
		Username: req.Username,
		Avatar:   req.Avatar,
		Email:    req.Email,
	}

	// 4. 调用业务逻辑
	if err := h.Service.Register(&user); err != nil {
		response.Fail(c, err.Error())
		return
	}

	// 4. 返回成功 (Data 为 nil 或 返回部分用户信息)
	response.Success(c, gin.H{"uid": user.ID})
}

// SendCode 发送验证码
func (h *UserHandler) SendCode(c *gin.Context) {
	var req struct {
		Mobile string `json:"mobile"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if req.Mobile == "" {
		response.Fail(c, "手机号不能为空")
		return
	}
	if err := h.Service.SendSMSCode(req.Mobile); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// LoginCode 验证码登录
func (h *UserHandler) LoginCode(c *gin.Context) {
	var req struct {
		Mobile string `json:"mobile"`
		Code   string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	ip := c.ClientIP()
	ua := c.Request.UserAgent()

	token, user, err := h.Service.LoginByCode(req.Mobile, req.Code, ip, ua)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"token":     token,
		"user_info": user,
	})
}

// Login 处理登录请求
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Mobile   string `json:"mobile"`
		Password string `json:"password"`
	}

	// 1. 绑定参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	// 获取 IP 和 UserAgent
	ip := c.ClientIP()
	ua := c.Request.UserAgent()

	// 2. 调用业务逻辑
	token, user, err := h.Service.Login(req.Mobile, req.Password, ip, ua)
	if err != nil {
		// 登录失败通常报 400 或 401
		response.Fail(c, err.Error())
		return
	}

	// 3. 返回成功数据
	response.Success(c, gin.H{
		"token":     token,
		"user_info": user, // 包含头像、余额等信息
	})
}

// Logout 处理退出登录请求
func (h *UserHandler) Logout(c *gin.Context) {
	// 从中间件中获取 userID
	userID, exists := c.Get("userID")
	if !exists {
		response.Success(c, nil) // 如果没有UserID (理论上中间件会拦截)，直接返回成功
		return
	}

	if err := h.Service.Logout(userID.(int64)); err != nil {
		// 即使 Redis 删除失败，也应该告诉前端退出成功，或者记录日志
		// 这里选择忽略错误，直接返回成功，确保前端能清除状态
	}

	response.Success(c, nil)
}

// Update 修改资料
func (h *UserHandler) Update(c *gin.Context) {
	userID, _ := c.Get("userID")
	// 使用 map 接收，支持部分更新
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	// 过滤允许更新的字段，防止恶意更新余额等
	allowed := []string{"avatar", "real_name", "gender", "email", "mobile", "age", "username"}
	updates := make(map[string]interface{})
	for _, field := range allowed {
		if val, ok := req[field]; ok {
			updates[field] = val
		}
	}

	if err := h.Service.UpdateInfo(userID.(int64), updates); err != nil {
		response.Fail(c, "修改失败: "+err.Error())
		return
	}
	response.Success(c, nil)
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if err := h.Service.ChangePassword(userID.(int64), req.OldPassword, req.NewPassword); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// ForgetPassword 忘记密码
func (h *UserHandler) ForgetPassword(c *gin.Context) {
	var req struct {
		Mobile      string `json:"mobile"`
		Code        string `json:"code"` // 验证码
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if err := h.Service.ResetPassword(req.Mobile, req.Code, req.NewPassword); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// Info 获取个人信息
func (h *UserHandler) Info(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.Service.GetInfo(userID.(int64))
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}
	response.Success(c, user)
}

// RegisterFace 人脸录入
// 支持两种方式：
// 1) multipart/form-data 传 file 字段，服务端上传 MinIO 得到 URL
// 2) application/json 直接传 face_image_url
func (h *UserHandler) RegisterFace(c *gin.Context) {
	userID, _ := c.Get("userID")
	uid := userID.(int64)

	faceImageURL := ""

	file, err := c.FormFile("file")
	if err == nil && file != nil {
		storage := service.StorageService{}
		url, _, uploadErr := storage.UploadMultipartFile(file, "face")
		if uploadErr != nil {
			response.Fail(c, "face image upload failed: "+uploadErr.Error())
			return
		}
		faceImageURL = url
	} else {
		var req struct {
			FaceImageURL string `json:"face_image_url"`
		}
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			response.Fail(c, "invalid request parameters")
			return
		}
		faceImageURL = strings.TrimSpace(req.FaceImageURL)
	}

	if faceImageURL == "" {
		response.Fail(c, "face image is required")
		return
	}

	if err := h.Service.RegisterFace(uid, faceImageURL); err != nil {
		response.Fail(c, "face register failed: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"face_registered": true,
		"face_image_url":  faceImageURL,
	})
}
