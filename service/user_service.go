package service

import (
	"errors"

	"github.com/0xirvan/url-shortener/model"
	"github.com/0xirvan/url-shortener/utils"
	"github.com/0xirvan/url-shortener/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserService interface {
	GetUsers(c *fiber.Ctx, params *validation.QueryUser) ([]model.User, int64, error)
	GetUserByID(c *fiber.Ctx, id string) (*model.User, error)
	GetUserByEmail(c *fiber.Ctx, email string) (*model.User, error)
	CreateUser(c *fiber.Ctx, req *validation.CreateUser) (*model.User, error)
	CreateGoogleUser(c *fiber.Ctx, req *validation.GoogleLogin) (*model.User, error)
	UpdateUser(c *fiber.Ctx, req *validation.UpdateUser, id string) (*model.User, error)
	UpdatePassOrVerifyUser(c *fiber.Ctx, req *validation.UpdatePassOrVerifyUser, id string) error
	DeleteUser(c *fiber.Ctx, id string) error
}

type userService struct {
	log      *logrus.Logger
	db       *gorm.DB
	validate *validator.Validate
}

func NewUserService(db *gorm.DB, validate *validator.Validate) UserService {
	return &userService{
		log:      utils.Log,
		db:       db,
		validate: validate,
	}
}

func (s *userService) GetUsers(c *fiber.Ctx, params *validation.QueryUser) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	if err := s.validate.Struct(params); err != nil {
		s.log.Error("GetUsers validation error: ", err)
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	query := s.db.WithContext(c.Context()).Order("created_at desc")

	if params.Search != "" {
		query = query.Where("name LIKE ? OR email LIKE ? OR role LIKE ?",
			"%"+params.Search+"%", "%"+params.Search+"%", "%"+params.Search+"%").Order("created_at desc")
	}

	res := query.Find(&users).Count(&total)
	if res.Error != nil {
		s.log.Error("GetUsers query error: ", res.Error)
		return nil, 0, res.Error
	}

	res = query.Limit(params.Limit).Offset(offset).Find(&users)
	if res.Error != nil {
		s.log.Error("GetUsers query error: ", res.Error)
		return nil, 0, res.Error
	}

	return users, total, nil

}

func (s *userService) GetUserByID(c *fiber.Ctx, id string) (*model.User, error) {
	var user model.User

	res := s.db.WithContext(c.Context()).Where("id = ?", id).First(&user)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		s.log.Error("GetUserByID query error: ", res.Error)
		return nil, res.Error
	}

	return &user, nil
}

func (s *userService) GetUserByEmail(c *fiber.Ctx, email string) (*model.User, error) {
	var user model.User

	res := s.db.WithContext(c.Context()).Where("email = ?", email).First(&user)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		s.log.Error("GetUserByEmail query error: ", res.Error)
		return nil, res.Error
	}

	return &user, nil
}

func (s *userService) CreateUser(c *fiber.Ctx, req *validation.CreateUser) (*model.User, error) {
	if err := s.validate.Struct(req); err != nil {
		s.log.Error("CreateUser validation error: ", err)
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.log.Error("CreateUser hash password error: ", err)
		return nil, err
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
	}

	res := s.db.WithContext(c.Context()).Create(user)

	if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
		s.log.Error("CreateUser duplicate key error: ", res.Error)
		return nil, fiber.NewError(fiber.StatusConflict, "Email already exists")
	}

	if res.Error != nil {
		s.log.Error("CreateUser query error: ", res.Error)
		return nil, res.Error
	}

	return user, nil

}

func (s *userService) CreateGoogleUser(c *fiber.Ctx, req *validation.GoogleLogin) (*model.User, error) {
	if err := s.validate.Struct(req); err != nil {
		s.log.Error("CreateGoogleUser validation error: ", err)
		return nil, err
	}

	userFromDB, err := s.GetUserByEmail(c, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user := &model.User{
				Name:          req.Name,
				Email:         req.Email,
				VerifiedEmail: req.VerifiedEmail,
			}

			if err := s.db.WithContext(c.Context()).Create(user).Error; err != nil {
				s.log.Error("CreateGoogleUser query error: ", err)
				return nil, err
			}
			return user, nil
		}
		s.log.Error("CreateGoogleUser query error: ", err)
		return nil, err
	}

	userFromDB.VerifiedEmail = req.VerifiedEmail
	if err := s.db.WithContext(c.Context()).Save(userFromDB).Error; err != nil {
		s.log.Error("CreateGoogleUser query error: ", err)
		return nil, err
	}

	return userFromDB, nil
}

func (s *userService) UpdateUser(c *fiber.Ctx, req *validation.UpdateUser, id string) (*model.User, error) {
	if err := s.validate.Struct(req); err != nil {
		s.log.Error("UpdateUser validation error: ", err)
		return nil, err
	}

	if req.Name == "" && req.Email == "" && req.Password == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "No fields to update")
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			s.log.Error("UpdateUser hash password error: ", err)
			return nil, err
		}
		req.Password = hashedPassword
	}

	updateBody := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	res := s.db.WithContext(c.Context()).Where("id = ?", id).Updates(updateBody)
	if res.Error != nil {
		s.log.Error("UpdateUser query error: ", res.Error)
		return nil, res.Error
	}

	if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
		s.log.Error("UpdateUser duplicate key error: ", res.Error)
		return nil, fiber.NewError(fiber.StatusConflict, "Email already exists")
	}

	if res.RowsAffected == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	user, err := s.GetUserByID(c, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdatePassOrVerifyUser(c *fiber.Ctx, req *validation.UpdatePassOrVerifyUser, id string) error {
	if err := s.validate.Struct(req); err != nil {
		s.log.Error("UpdatePassOrVerifyUser validation error: ", err)
		return err
	}

	if req.Password == "" && !req.VerifiedEmail {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			s.log.Error("UpdatePassOrVerifyUser hash password error: ", err)
			return err
		}

		req.Password = hashedPassword
	}

	updateBody := &model.User{
		Password:      req.Password,
		VerifiedEmail: req.VerifiedEmail,
	}

	res := s.db.WithContext(c.Context()).Where("id = ?", id).Updates(updateBody)
	if res.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if res.Error != nil {
		s.log.Error("UpdatePassOrVerifyUser query error: ", res.Error)
		return res.Error
	}

	return nil

}

func (s *userService) DeleteUser(c *fiber.Ctx, id string) error {
	res := s.db.WithContext(c.Context()).Where("id = ?", id).Delete(&model.User{})
	if res.Error != nil {
		s.log.Error("DeleteUser query error: ", res.Error)
		return res.Error
	}

	if res.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return nil
}
