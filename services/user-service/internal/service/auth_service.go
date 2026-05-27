package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"smartcommunity-microservices/pkg/auth"
	"smartcommunity-microservices/services/user-service/internal/model"
	"smartcommunity-microservices/services/user-service/internal/repository"

	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo  *repository.UserRepo
	roleRepo  *repository.RoleRepo
	logRepo   *repository.LoginLogRepo
	resetRepo *repository.PasswordResetRepo
	rdb       *goredis.Client
	jwtSecret string
	jwtTTL    time.Duration
}

func NewAuthService(
	userRepo *repository.UserRepo,
	roleRepo *repository.RoleRepo,
	logRepo *repository.LoginLogRepo,
	resetRepo *repository.PasswordResetRepo,
	rdb *goredis.Client,
	jwtSecret string,
	jwtTTL time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		logRepo:   logRepo,
		resetRepo: resetRepo,
		rdb:       rdb,
		jwtSecret: jwtSecret,
		jwtTTL:    jwtTTL,
	}
}

type LoginResult struct {
	Token            string
	User             *model.SysUser
	IsNewUser        bool
	ProfileCompleted bool
}

type RegisterRequest struct {
	Mobile   string `json:"mobile" binding:"required"`
	Password string `json:"password" binding:"required"`
	RealName string `json:"real_name"`
	Age      int    `json:"age"`
	Gender   int    `json:"gender"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
}

// AUTH-001: Register
func (s *AuthService) Register(req RegisterRequest) (*model.SysUser, error) {
	count, err := s.userRepo.CountByMobile(req.Mobile)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("手机号已注册")
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	realName := req.RealName
	if realName == "" {
		realName = "未完善资料"
	}

	age := req.Age
	if age <= 0 {
		age = 1
	}

	gender := req.Gender
	if gender != 0 && gender != 1 && gender != 2 {
		gender = 0
	}

	avatar := req.Avatar
	if avatar == "" {
		avatar = "https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png"
	}
	username := req.Username
	if username == "" {
		username = req.Mobile
	}

	user := &model.SysUser{
		Mobile:   req.Mobile,
		Password: hash,
		RealName: realName,
		Age:      age,
		Gender:   gender,
		Username: username,
		Avatar:   avatar,
		Email:    req.Email,
		Role:     "user",
		Status:   1,
		Balance:  100.00,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// AUTH-002: Login
func (s *AuthService) Login(mobile, password, ip, userAgent string) (string, *model.SysUser, error) {
	user, err := s.userRepo.FindByMobile(mobile)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logLoginFailure(0, mobile, ip, userAgent, "用户不存在", "user")
			return "", nil, errors.New("手机号或密码错误")
		}
		return "", nil, err
	}

	if !auth.CheckPasswordHash(password, user.Password) {
		s.logLoginFailure(user.ID, mobile, ip, userAgent, "密码错误", user.Role)
		return "", nil, errors.New("手机号或密码错误")
	}

	if user.Status == 0 {
		s.logLoginFailure(user.ID, mobile, ip, userAgent, "账户已冻结", user.Role)
		return "", nil, errors.New("账户已被冻结")
	}

	token, err := auth.GenerateToken(s.jwtSecret, user.ID, user.Role, s.jwtTTL)
	if err != nil {
		return "", nil, err
	}

	ctx := context.Background()
	redisKey := fmt.Sprintf("login:token:%d", user.ID)
	if err := s.rdb.Set(ctx, redisKey, token, s.jwtTTL).Err(); err != nil {
		return "", nil, err
	}

	// LOG-001/002: write login log
	if user.Role == "admin" {
		_ = s.logRepo.CreateAdminLog(&model.AdminLoginLog{
			AdminUserID: user.ID,
			Mobile:      mobile,
			LoginTime:   time.Now(),
			IP:          ip,
			UserAgent:   userAgent,
			Success:     true,
		})
	} else {
		_ = s.logRepo.CreateUserLog(&model.UserLoginLog{
			UserID:    user.ID,
			Mobile:    mobile,
			LoginTime: time.Now(),
			IP:        ip,
			UserAgent: userAgent,
			Success:   true,
		})
	}

	return token, user, nil
}

func (s *AuthService) issueToken(user *model.SysUser) (string, error) {
	token, err := auth.GenerateToken(s.jwtSecret, user.ID, user.Role, s.jwtTTL)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	redisKey := fmt.Sprintf("login:token:%d", user.ID)
	if err := s.rdb.Set(ctx, redisKey, token, s.jwtTTL).Err(); err != nil {
		return "", err
	}
	return token, nil
}

// AUTH-008a: Send login SMS code. This is separate from password-reset codes
// because it allows unregistered mobiles to continue into one-click signup.
func (s *AuthService) SendLoginCode(mobile string) (string, error) {
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	ctx := context.Background()
	redisKey := fmt.Sprintf("sms:login:%s", mobile)
	if err := s.rdb.Set(ctx, redisKey, code, 5*time.Minute).Err(); err != nil {
		return "", err
	}

	if err := sendSMSViaSpug(mobile, code); err != nil {
		return "", err
	}

	return code, nil
}

func sendSMSViaSpug(mobile, code string) error {
	url := "https://push.spug.cc/send/nbONk8gz2Vr34gXG"
	payload := map[string]interface{}{
		"code":    code,
		"targets": mobile,
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("短信发送失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("短信服务异常, 状态码: %d", resp.StatusCode)
	}
	return nil
}

// AUTH-008b: Login by SMS code. If the mobile is not registered yet, create a
// normal resident account and bind it to the "user" role.
func (s *AuthService) LoginByCode(mobile, code, ip, userAgent string) (*LoginResult, error) {
	ctx := context.Background()
	redisKey := fmt.Sprintf("sms:login:%s", mobile)

	storedCode, err := s.rdb.Get(ctx, redisKey).Result()
	if err != nil {
		return nil, errors.New("验证码已过期或未发送")
	}
	if storedCode != code {
		return nil, errors.New("验证码错误")
	}
	_ = s.rdb.Del(ctx, redisKey)

	user, err := s.userRepo.FindByMobile(mobile)
	isNewUser := false
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		hash, err := auth.HashPassword("123456")
		if err != nil {
			return nil, err
		}

		user = &model.SysUser{
			Mobile:   mobile,
			Password: hash,
			RealName: "未完善资料",
			Age:      1,
			Gender:   0,
			Username: mobile,
			Avatar:   "https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png",
			Role:     "user",
			Status:   1,
			Balance:  0,
		}
		if err := s.userRepo.Create(user); err != nil {
			return nil, fmt.Errorf("自动注册失败: %w", err)
		}
		isNewUser = true

		if s.roleRepo != nil {
			if role, err := s.roleRepo.FindByCode("user"); err == nil {
				_ = s.roleRepo.BindUserRoles(user.ID, []int64{role.ID})
			}
		}
	} else if user.Status == 0 {
		s.logLoginFailure(user.ID, mobile, ip, userAgent, "账户已冻结", user.Role)
		return nil, errors.New("账户已被冻结")
	}

	token, err := s.issueToken(user)
	if err != nil {
		return nil, err
	}

	_ = s.logRepo.CreateUserLog(&model.UserLoginLog{
		UserID:    user.ID,
		Mobile:    mobile,
		LoginTime: time.Now(),
		IP:        ip,
		UserAgent: userAgent,
		Success:   true,
	})

	return &LoginResult{
		Token:            token,
		User:             user,
		IsNewUser:        isNewUser,
		ProfileCompleted: !isAutoRegisteredProfile(user),
	}, nil
}

func isAutoRegisteredProfile(user *model.SysUser) bool {
	return user.RealName == "" || user.RealName == "未完善资料" || user.Age <= 1
}

func (s *AuthService) logLoginFailure(userID int64, mobile, ip, userAgent, reason, role string) {
	now := time.Now()
	if role == "admin" {
		_ = s.logRepo.CreateAdminLog(&model.AdminLoginLog{
			AdminUserID:   userID,
			Mobile:        mobile,
			LoginTime:     now,
			IP:            ip,
			UserAgent:     userAgent,
			Success:       false,
			FailureReason: reason,
		})
	} else {
		_ = s.logRepo.CreateUserLog(&model.UserLoginLog{
			UserID:        userID,
			Mobile:        mobile,
			LoginTime:     now,
			IP:            ip,
			UserAgent:     userAgent,
			Success:       false,
			FailureReason: reason,
		})
	}
}

// AUTH-007: Logout
func (s *AuthService) Logout(userID int64) error {
	ctx := context.Background()
	redisKey := fmt.Sprintf("login:token:%d", userID)
	return s.rdb.Del(ctx, redisKey).Err()
}

// AUTH-004: ChangePassword
func (s *AuthService) ChangePassword(userID int64, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if !auth.CheckPasswordHash(oldPassword, user.Password) {
		return errors.New("原密码错误")
	}

	hash, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(userID, hash)
}

// AUTH-003a: SendResetCode
func (s *AuthService) SendResetCode(mobile string) error {
	user, err := s.userRepo.FindByMobile(mobile)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("该手机号未注册")
		}
		return err
	}
	_ = user // user exists

	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	ctx := context.Background()
	redisKey := fmt.Sprintf("sms:reset:%s", mobile)
	if err := s.rdb.Set(ctx, redisKey, code, 5*time.Minute).Err(); err != nil {
		return err
	}

	// Store bcrypt hash in DB for audit (not plaintext)
	codeHash, _ := auth.HashPassword(code)
	_ = s.resetRepo.Create(&model.PasswordResetCode{
		Mobile:    mobile,
		CodeHash:  codeHash,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	})

	if err := sendSMSViaSpug(mobile, code); err != nil {
		return err
	}

	return nil
}

// AUTH-003b: ResetPassword
func (s *AuthService) ResetPassword(mobile, code, newPassword string) error {
	ctx := context.Background()
	redisKey := fmt.Sprintf("sms:reset:%s", mobile)

	storedCode, err := s.rdb.Get(ctx, redisKey).Result()
	if err != nil {
		return errors.New("验证码已过期或未发送")
	}
	if storedCode != code {
		return errors.New("验证码错误")
	}

	_ = s.rdb.Del(ctx, redisKey)

	hash, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByMobile(mobile)
	if err != nil {
		return err
	}

	// Mark the most recent reset code as used
	_ = s.resetRepo.MarkUsedByMobile(mobile)

	return s.userRepo.UpdatePassword(user.ID, hash)
}
