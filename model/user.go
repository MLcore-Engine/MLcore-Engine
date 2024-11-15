package model

import (
	"MLcore-Engine/common"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username         string     `json:"username" gorm:"uniqueIndex;not null" validate:"required,max=12"`
	Password         string     `json:"-" gorm:"not null" validate:"required,min=8,max=20"` // use json:"-" to hide password in JSON
	DisplayName      string     `json:"display_name" gorm:"index" validate:"max=20"`
	Role             int        `json:"role" gorm:"type:int;default:1"`   // 1000: root, 100: admin, 1: common
	Status           int        `json:"status" gorm:"type:int;default:1"` // 1: enabled, 0: disabled
	Email            string     `json:"email" validate:"max=50,email"`
	GitHubID         string     `json:"github_id" gorm:"column:github_id"`
	WeChatID         string     `json:"wechat_id" gorm:"column:wechat_id"`
	VerificationCode string     `json:"-" gorm:"-"`
	Projects         []Project  `json:"projects" gorm:"many2many:user_projects"`
	Notebooks        []Notebook `json:"notebooks" gorm:"foreignKey:UserID;constraint:OnDelete:RESTRICT"`
}

func GetMaxUserId() uint {
	var user User
	DB.Last(&user)
	return user.ID
}

func GetAllUsers(offset, limit int) (users []*User, total int64, err error) {

	err = DB.Model(&User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = DB.Order("id desc").
		Limit(limit).
		Offset(offset).
		Select([]string{"id", "username", "display_name", "role", "status", "email"}).
		Find(&users).Error
	return users, total, err
}

func SearchUsers(keyword string, offset, limit int) (users []*User, total int64, err error) {
	query := DB.Model(&User{}).Select("id", "username", "display_name", "role", "status", "email")

	keywordLike := "%" + keyword + "%"
	searchQuery := query.Where(DB.Where("id LIKE ?", keywordLike).
		Or(DB.Where("username LIKE ?", keywordLike)).
		Or(DB.Where("email LIKE ?", keywordLike)).
		Or(DB.Where("display_name LIKE ?", keywordLike)))

	err = searchQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = searchQuery.Limit(limit).
		Offset(offset).
		Order("id DESC").
		Find(&users).Error

	return users, total, err
}

func GetUserById(id uint, selectAll bool) (*User, error) {
	if id == 0 {
		return nil, errors.New("id can not be empty！")
	}

	var user User
	var err error = nil
	if selectAll {
		err = DB.Where("id = ?", id).First(&user).Error
	} else {
		err = DB.Select([]string{"id", "username", "display_name", "role", "status", "email", "wechat_id", "github_id"}).First(&user, "id = ?", id).Error
	}
	return &user, err
}

func DeleteUserById(id int) (err error) {
	if id == 0 {
		return errors.New("id can not be empty！")
	}
	var user User
	err = DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return err
	}
	return user.Delete()
}

func (user *User) Insert() error {
	var err error
	if user.Password != "" {
		user.Password, err = common.Password2Hash(user.Password)
		if err != nil {
			return err
		}
	}
	err = DB.Create(user).Error
	return err
}

func (user *User) Update(updatePassword bool) error {
	var err error
	if updatePassword {
		user.Password, err = common.Password2Hash(user.Password)
		if err != nil {
			return err
		}
	}
	err = DB.Model(user).Updates(user).Error
	return err
}

func (user *User) Delete() error {
	if user.ID == 0 {
		return errors.New("id 为空！")
	}
	err := DB.Delete(user).Error
	return err
}

// ValidateAndFill check password & user status
func (user *User) ValidateAndFill() (err error) {
	// When querying with struct, GORM will only query with non-zero fields,
	// that means if your field’s value is 0, '', false or other zero values,
	// it won’t be used to build query conditions
	password := user.Password
	if user.Username == "" || password == "" {
		return errors.New("username or password is empty")
	}
	DB.Preload("Projects").Where(User{Username: user.Username}).First(user)
	okay := common.ValidatePasswordAndHash(password, user.Password)
	if !okay || user.Status != common.UserStatusEnabled {
		return errors.New("username or password is incorrect, or user is banned")
	}
	return nil
}

func (user *User) FillUserById() error {
	if user.ID == 0 {
		return errors.New("id is empty")
	}

	result := DB.Where("id = ?", user.ID).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with ID %d not found", user.ID)
		}
		return fmt.Errorf("error querying user: %w", result.Error)
	}

	return nil
}

func (user *User) FillUserByEmail() error {
	if user.Email == "" {
		return errors.New("email is empty")
	}

	res := DB.Where("email = ?", user.Email).First(user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user's email %s not found", user.Email)
		}
		return fmt.Errorf("find user has a error %w", res.Error)
	}

	return nil
}

func (user *User) FillUserByGitHubId() error {
	if user.GitHubID == "" {
		return errors.New("GitHub ID is empty")
	}

	result := DB.Where("github_id = ?", user.GitHubID).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with GitHub ID %s not found", user.GitHubID)
		}
		return fmt.Errorf("error querying user: %w", result.Error)
	}

	return nil
}

func (user *User) FillUserByWeChatId() error {
	if user.WeChatID == "" {
		return errors.New("WeChat ID is empty")
	}

	result := DB.Where("wechat_id = ?", user.WeChatID).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with WeChat ID %s not found", user.WeChatID)
		}
		return fmt.Errorf("error querying user: %w", result.Error)
	}

	return nil
}

func (user *User) FillUserByUsername() error {
	if user.Username == "" {
		return errors.New("username is empty")
	}

	result := DB.Where("username = ?", user.Username).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with username %s not found", user.Username)
		}
		return fmt.Errorf("error querying user: %w", result.Error)
	}

	return nil
}

func ValidateUserToken(token string) (*User, error) {
	if token == "" {
		return nil, errors.New("token is empty")
	}

	// Remove "Bearer " prefix if present
	token = strings.TrimPrefix(token, "Bearer ")

	var user User
	result := DB.Where("token = ?", token).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no user found with the provided token")
		}
		return nil, fmt.Errorf("error validating token: %w", result.Error)
	}

	return &user, nil
}

func IsEmailAlreadyTaken(email string) (bool, error) {
	var count int64
	result := DB.Model(&User{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("error checking email: %w", result.Error)
	}
	return count > 0, nil
}

func IsWeChatIdAlreadyTaken(wechatId string) (bool, error) {
	var count int64
	result := DB.Model(&User{}).Where("wechat_id = ?", wechatId).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("error checking WeChat ID: %w", result.Error)
	}
	return count > 0, nil
}

func IsGitHubIdAlreadyTaken(githubId string) (bool, error) {
	var count int64
	result := DB.Model(&User{}).Where("github_id = ?", githubId).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("error checking GitHub ID: %w", result.Error)
	}
	return count > 0, nil
}

func ResetUserPasswordByEmail(email string, password string) error {
	if email == "" || password == "" {
		return errors.New("email or password is empty")
	}
	hashedPassword, err := common.Password2Hash(password)
	if err != nil {
		return err
	}
	err = DB.Model(&User{}).Where("email = ?", email).Update("password", hashedPassword).Error
	return err
}
