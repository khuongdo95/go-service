package services

import (
	"context"

	"github.com/khuongdo95/go-pkg/common/response"
	"github.com/khuongdo95/go-service/internal/adapter/dtos"
	"github.com/khuongdo95/go-service/internal/generated/ent"
	entuser "github.com/khuongdo95/go-service/internal/generated/ent/user"
	entipwhitelist "github.com/khuongdo95/go-service/internal/generated/ent/useripwhitelist"
	"github.com/khuongdo95/go-service/internal/usecase/services/signin"
	"github.com/khuongdo95/go-service/internal/usecase/services/signup"
	token "github.com/khuongdo95/go-service/internal/usecase/services/token"
	"github.com/khuongdo95/go-service/internal/usecase/services/transformer"
)

type UserService struct {
	entUser *ent.Client
	signUp  signup.SignUp
	signIn  signin.SignIn
}

type UserServiceIf interface {
	Get(ctx context.Context, id string) (*ent.User, *response.AppError)
	SignUp(ctx context.Context, username, password string) (*token.Tokens, *response.AppError)
	SignIn(ctx context.Context, username, password string) (*token.Tokens, *response.AppError)
	Create(ctx context.Context, req *dtos.CreateUserReq) (*dtos.UserRes, *response.AppError)
	Update(ctx context.Context, req *dtos.UpdateUserReq) (*dtos.UserRes, *response.AppError)
	Delete(ctx context.Context, id string) *response.AppError
	List(ctx context.Context, req *dtos.ListUserReq) (*dtos.ListUserRes, *response.AppError)
}

func NewUserService(entUser *ent.Client, signUp signup.SignUp, signIn signin.SignIn) *UserService {
	return &UserService{
		entUser: entUser,
		signUp:  signUp,
		signIn:  signIn,
	}
}

func (u *UserService) Get(ctx context.Context, id string) (*ent.User, *response.AppError) {
	user, err := u.entUser.User.Query().Where(entuser.ID(id)).WithIdentity().WithIPWhiteList().Only(ctx)
	if err != nil {
		return nil, response.ConvertDatabaseError(err)
	}
	return user, nil
}

func (u *UserService) SignUp(ctx context.Context, username, password string) (*token.Tokens, *response.AppError) {
	return nil, nil
}

func (u *UserService) SignIn(ctx context.Context, username, password string) (*token.Tokens, *response.AppError) {
	return u.signIn.SignIn(ctx, username, password)
}

func (u *UserService) Create(ctx context.Context, req *dtos.CreateUserReq) (*dtos.UserRes, *response.AppError) {
	request := &dtos.SignUpReq{
		Username: req.Username,
		Password: req.Password,
		Name:     req.Name,
		Email:    req.Email,
	}
	data, err := u.signUp.SignUp(ctx, "CLIENT", *req.IpAddress, request)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (u *UserService) Update(ctx context.Context, req *dtos.UpdateUserReq) (*dtos.UserRes, *response.AppError) {
	userDB, err := u.entUser.User.UpdateOneID(req.Id).SetNillableName(req.Name).SetNillableEmail(req.Email).Save(ctx)
	if err != nil {
		return nil, response.ConvertDatabaseError(err)
	}
	err = u.entUser.UserIpWhiteList.Update().Where(entipwhitelist.UserIDEQ(userDB.ID)).SetNillableIPAddress(req.IpAddress).Exec(ctx)
	if err != nil {
		return nil, response.ConvertDatabaseError(err)
	}
	return nil, nil
}

func (u *UserService) Delete(ctx context.Context, id string) *response.AppError {
	err := u.entUser.User.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return response.ConvertDatabaseError(err)
	}
	return nil
}

func (u *UserService) List(ctx context.Context, req *dtos.ListUserReq) (*dtos.ListUserRes, *response.AppError) {
	var userResult []*dtos.UserRes

	query := u.entUser.User.Query()
	count, err := query.Count(ctx)
	if err != nil {
		return nil, response.ConvertDatabaseError(err)
	}

	enitites, err := query.
		Limit(int(req.Pagination.PageSize)).
		Offset(int(req.Pagination.PageIndex * req.Pagination.PageSize)).
		WithIPWhiteList().
		WithIdentity().
		All(ctx)
	if err != nil {
		return nil, response.ConvertDatabaseError(err)
	}

	for _, v := range enitites {
		userResult = append(userResult, transformer.TransformerUser(v))
	}

	return &dtos.ListUserRes{
		Users: userResult,
		Pagination: &dtos.PaginationRes{
			PageSize:  req.Pagination.PageSize,
			PageIndex: req.Pagination.PageIndex,
			Total:     int32(count),
		},
	}, nil
}
