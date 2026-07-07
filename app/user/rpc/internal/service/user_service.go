package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"smartcommunity-microservices/app/user/rpc/internal/model"
	"smartcommunity-microservices/app/user/rpc/internal/repository"
	"smartcommunity-microservices/common/faceauth"
	"smartcommunity-microservices/common/mq"

	goredis "github.com/redis/go-redis/v9"
)

type UserService struct {
	userRepo *repository.UserRepo
	rdb      *goredis.Client
	mqClient *mq.Client
}

func NewUserService(userRepo *repository.UserRepo, rdb *goredis.Client, mqClient *mq.Client) *UserService {
	return &UserService{userRepo: userRepo, rdb: rdb, mqClient: mqClient}
}

type UpdateProfileRequest struct {
	Avatar   *string `json:"avatar"`
	Mobile   *string `json:"mobile"`
	Username *string `json:"username"`
	Gender   *int    `json:"gender"`
	Email    *string `json:"email"`
	RealName *string `json:"real_name"`
	Age      *int    `json:"age"`
}

// AUTH-005 + AUTH-006: GetProfile
func (s *UserService) GetProfile(userID int64) (*model.SysUser, error) {
	return s.userRepo.FindByID(userID)
}

// RegisterFace registers a face image URL for the user.
func (s *UserService) RegisterFace(userID int64, faceImageURL string) error {
	faceImageURL = strings.TrimSpace(faceImageURL)
	if faceImageURL == "" {
		return errors.New("face image url is required")
	}
	if err := faceauth.ValidateEnrollment(context.Background(), faceImageURL); err != nil {
		return err
	}
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	if user.FaceImageURL != "" && user.FaceImageURL != faceImageURL {
		s.publishCleanupEvent(user.FaceImageURL)
	}
	return s.userRepo.UpdateFields(user.ID, map[string]interface{}{
		"face_registered": true,
		"face_image_url":  faceImageURL,
	})
}

// AUTH-005: UpdateProfile
func (s *UserService) UpdateProfile(userID int64, req UpdateProfileRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	fields := make(map[string]interface{})
	if req.Avatar != nil {
		fields["avatar"] = *req.Avatar
		if user.Avatar != "" && user.Avatar != *req.Avatar {
			s.publishCleanupEvent(user.Avatar)
		}
	}
	if req.Mobile != nil {
		fields["mobile"] = *req.Mobile
	}
	if req.Username != nil {
		fields["username"] = *req.Username
	}
	if req.Gender != nil {
		fields["gender"] = *req.Gender
	}
	if req.Email != nil {
		fields["email"] = *req.Email
	}
	if req.RealName != nil {
		fields["real_name"] = *req.RealName
	}
	if req.Age != nil {
		fields["age"] = *req.Age
	}
	if len(fields) == 0 {
		return nil
	}
	return s.userRepo.UpdateFields(userID, fields)
}

func (s *UserService) publishCleanupEvent(url string) {
	if s.mqClient == nil || url == "" {
		return
	}
	event := map[string]string{"url": url}
	body, err := json.Marshal(event)
	if err != nil {
		return
	}
	_ = s.mqClient.PublishEvent(context.Background(), "file.cleanup", body)
}
