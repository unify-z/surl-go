package UserService

import (
	"github.com/unify-z/go-surl/internal/errors"
	"github.com/unify-z/go-surl/internal/helpers"
	"github.com/unify-z/go-surl/internal/repository/surl"
	User2 "github.com/unify-z/go-surl/internal/repository/user"
	"github.com/unify-z/go-surl/internal/service/verify_code"
	"github.com/unify-z/go-surl/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserManager struct {
	repo          *User2.UserRepo
	verifyCodeSvc *verify_code.VerifyCodeManager
	jwtHelper     *helpers.JWTHelper
}

func NewUserManager(userRepo *User2.UserRepo, jwtHelper *helpers.JWTHelper, verifyCodeSvc *verify_code.VerifyCodeManager) *UserManager {
	return &UserManager{repo: userRepo, jwtHelper: jwtHelper, verifyCodeSvc: verifyCodeSvc}
}

func (um *UserManager) RegisterUser(username, passwordMD5 string, email string, code string, isRegisterByAdmin bool) (*User2.UserModel, error) {
	var isAdmin bool
	userCount, err := um.repo.CountUsers()
	if err != nil {
		return nil, err
	}
	if userCount == 0 {
		isAdmin = true
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordMD5), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &User2.UserModel{
		Username:     username,
		PasswordHash: string(hashedPassword),
		IsAdmin:      isAdmin,
		IsBanned:     false,
		Email:        email,
	}
	// skip code verification when if registered by admin
	if !isRegisterByAdmin {
		codeVerifyResult := um.verifyCodeSvc.VerifyCode(email, code)
		if !codeVerifyResult {
			return nil, errors.ErrInvalidVerifyCode
		}
	}
	err = um.repo.CreateUser(user)
	if err != nil {
		if utils.IsUniqueConstraintError(err) {
			return nil, errors.ErrUserAlreadyExists
		}
		return nil, err
	}
	return user, nil
}

func (um *UserManager) LoginUser(username, passwordMD5 string) (*User2.UserModel, error, string) {
	user, err := um.repo.GetUserByUsername(username)
	if err != nil {
		if user == nil {
			return nil, errors.ErrUserNotFound, ""
		}
		return nil, err, ""
	}
	passwdErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(passwordMD5))
	if passwdErr != nil {
		return nil, errors.ErrInvalidCredentials, ""
	}
	if user.IsBanned {
		return nil, errors.ErrUserBanned, ""
	}
	token, err := um.jwtHelper.GenerateUserToken(user.ID, user.Username, user.IsAdmin)
	if err != nil {
		return nil, err, ""
	}
	return user, nil, token
}

func (um *UserManager) GetUserByID(id uint) (*User2.UserModel, error) {
	user, err := um.repo.GetUserByID(id)
	if err != nil {
		if user == nil {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (um *UserManager) UpdateUser(user *User2.UserModel) error {
	return um.repo.UpdateUser(user)
}

func (um *UserManager) DeleteUser(id uint) error {
	return um.repo.DeleteUser(id)
}

func (um *UserManager) ListUsersWithoutSURLs() ([]User2.UserModel, error) {
	return um.repo.ListUsersWithoutSURLs()
}

func (um *UserManager) GetUserSURLs(id uint) ([]surl.ShortURL, error) {
	return um.repo.GetUserSURLS(id)
}

func (um *UserManager) ValidateJWTToken(token string) (*User2.UserModel, error) {
	result, err := um.jwtHelper.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	userIDInt, ok := result["user_id"].(float64)
	if !ok {
		return nil, errors.ErrInvalidToken
	}
	user, err := um.GetUserByID(uint(userIDInt))
	if err != nil {
		return nil, err
	}
	if user.IsBanned {
		return nil, errors.ErrUserBanned
	}
	return user, nil
}
