package authrepo

import (
	"errors"
	"fmt"
	"time"

	authmodel "github.com/DavG20/Tikus_Event_Api/Internal/pkg/auth/Auth_Model"
	constants "github.com/DavG20/Tikus_Event_Api/pkg/Utils/Constants"
	helper "github.com/DavG20/Tikus_Event_Api/pkg/Utils/Helper"

	"gorm.io/gorm"
)

type AuthRepo struct {
	DB *gorm.DB
}

func NewAuth(db *gorm.DB) AuthRepo {
	return AuthRepo{DB: db}
}

func (authRepo *AuthRepo) CreateUser(userInput *authmodel.UserRegisterInput) (dbResponse *authmodel.DBResponse, err error) {
	hashPassword, err := helper.HashPassword(userInput.Password)
	if err != nil {
		fmt.Println("error hashing password")
		return nil, err
	}
	user := &authmodel.AuthModel{
		UserName:   userInput.UserName,
		Email:      userInput.Email,
		Password:   hashPassword,
		CreatedOn:  time.Now().Format("2006-01-02 15:04:05.12"),
		ProfileUrl: userInput.ProfileUrl,
	}
	fmt.Println(user)
	fmt.Println("userInput", userInput)
	if user == nil {
		return nil, errors.New("invalid user input to decode user")
	}

	err = authRepo.DB.Table(constants.UserTableName).Create(user).Error
	if err != nil {
		fmt.Println("error in create method repo , line 38")
		return nil, errors.New("Invalid input")
	}
	dbResponse, err = authRepo.GetDbResponse(user)
	if err != nil {
		fmt.Println("eror getting dbresponse")
		return nil, err
	}
	if dbResponse == nil {
		fmt.Println("error while decoding to dbresponse from input user")
		return nil, errors.New("empty dbresponse")
	}
	return dbResponse, nil
}
func (authRepo *AuthRepo) GetDbResponse(user *authmodel.AuthModel) (dbResponse *authmodel.DBResponse, err error) {
	dbResponse = &authmodel.DBResponse{
		UserId:     user.UserId,
		UserName:   user.UserName,
		Email:      user.Email,
		CreatedOn:  user.CreatedOn,
		ProfileUrl: user.ProfileUrl,
		Admin:      user.Admin,
	}
	if dbResponse == nil {
		return nil, errors.New("Invalid db response")
	}
	return dbResponse, nil

}

func (authRepo *AuthRepo) FindUserByUserName(userName string) (user *authmodel.AuthModel, err error) {
	// res:=authRepo.DB.Table(constants.UserTableName).Raw("select * from user_tikus_event where user_name=? ",userName).Find(&user)
	if err = authRepo.DB.Table(constants.UserTableName).Where(authmodel.AuthModel{UserName: userName}).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("inhere")
		return nil, err
	}
	return user, nil

}

func (authRepo *AuthRepo) FindUserByEmail(email string) (user *authmodel.AuthModel, err error) {
	err = authRepo.DB.Table(constants.UserTableName).Where("email=?", email).First(&user).Error
	if err != nil {

		return nil, err
	}

	return user, nil
}
