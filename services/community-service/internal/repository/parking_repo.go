package repository

import (
	"errors"
	"strings"
	"time"

	"smartcommunity-microservices/services/community-service/internal/model"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ParkingRepo struct {
	db *gorm.DB
}

func NewParkingRepo(db *gorm.DB) *ParkingRepo {
	return &ParkingRepo{db: db}
}

func (r *ParkingRepo) List(page, size int) ([]model.ParkingSpaceAdminView, int64, error) {
	var items []model.ParkingSpaceAdminView
	var total int64
	q := r.db.Model(&model.ParkingSpace{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Select(`
			parking_spaces.*,
			COALESCE(bindings.id, 0) AS binding_id,
			COALESCE(bindings.user_id, 0) AS user_id,
			COALESCE(bindings.car_plate, '') AS car_plate,
			COALESCE(NULLIF(sys_user.real_name, ''), NULLIF(sys_user.username, ''), '') AS user_name,
			COALESCE(sys_user.mobile, '') AS user_mobile
		`).
		Joins("LEFT JOIN user_parking_bindings AS bindings ON bindings.parking_space_id = parking_spaces.id AND bindings.status = 1").
		Joins("LEFT JOIN sys_user ON sys_user.id = bindings.user_id").
		Order("parking_spaces.id ASC").
		Offset((page - 1) * size).
		Limit(size).
		Scan(&items).Error
	return items, total, err
}

func (r *ParkingRepo) Create(item *model.ParkingSpace) error {
	err := r.db.Create(item).Error
	if isDuplicateEntry(err) {
		return ErrDuplicateParkingNo
	}
	return err
}

func (r *ParkingRepo) ListBindingsByUser(userID int64) ([]model.UserParkingBinding, error) {
	var items []model.UserParkingBinding
	err := r.db.Preload("ParkingSpace").Where("user_id = ? AND status = ?", userID, 1).Find(&items).Error
	return items, err
}

func (r *ParkingRepo) Assign(spaceID int64, mobile, carPlate string) (*model.UserParkingBinding, error) {
	var binding model.UserParkingBinding
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var space model.ParkingSpace
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&space, spaceID).Error; err != nil {
			return err
		}

		var activeBinding model.UserParkingBinding
		activeErr := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("parking_space_id = ? AND status = ?", spaceID, 1).
			First(&activeBinding).Error
		if activeErr != nil && !errors.Is(activeErr, gorm.ErrRecordNotFound) {
			return activeErr
		}

		mobile = strings.TrimSpace(mobile)
		if mobile == "" {
			now := time.Now()
			if activeBinding.ID != 0 {
				activeBinding.Status = 0
				activeBinding.UnboundAt = &now
				if err := tx.Save(&activeBinding).Error; err != nil {
					return err
				}
				binding = activeBinding
			}
			space.Status = 0
			return tx.Save(&space).Error
		}

		var user struct {
			ID int64
		}
		if err := tx.Table("sys_user").Select("id").Where("mobile = ?", mobile).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrParkingUserNotFound
			}
			return err
		}

		space.Status = 1
		if err := tx.Save(&space).Error; err != nil {
			return err
		}
		if activeBinding.ID != 0 {
			activeBinding.UserID = user.ID
			activeBinding.CarPlate = strings.TrimSpace(carPlate)
			activeBinding.UnboundAt = nil
			if err := tx.Save(&activeBinding).Error; err != nil {
				return err
			}
			binding = activeBinding
			return nil
		}

		binding = model.UserParkingBinding{
			UserID:         user.ID,
			ParkingSpaceID: spaceID,
			CarPlate:       strings.TrimSpace(carPlate),
			Status:         1,
			BoundAt:        time.Now(),
		}
		return tx.Create(&binding).Error
	})
	return &binding, err
}

func (r *ParkingRepo) BindPlate(bindingID, userID int64, carPlate string) (*model.UserParkingBinding, error) {
	var binding model.UserParkingBinding
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ? AND status = ?", bindingID, userID, 1).
			First(&binding).Error; err != nil {
			return err
		}
		binding.CarPlate = carPlate
		return tx.Save(&binding).Error
	})
	return &binding, err
}

func (r *ParkingRepo) Stats() (map[string]int64, error) {
	var total, free, bound int64
	if err := r.db.Model(&model.ParkingSpace{}).Count(&total).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.ParkingSpace{}).Where("status = ?", 0).Count(&free).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.ParkingSpace{}).Where("status = ?", 1).Count(&bound).Error; err != nil {
		return nil, err
	}
	return map[string]int64{"total": total, "free": free, "bound": bound, "used": bound}, nil
}

func isDuplicateEntry(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
