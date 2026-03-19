package userhandler

import (
	"github.com/dehwyy/x-balance/internal/application/service/userservice"
	userspb "github.com/dehwyy/x-balance/internal/generated/pb/users/v1"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	UserService *userservice.Service
}

type Handler struct {
	userspb.UnimplementedUserServiceServer
	userservice *userservice.Service
}

func New(opts Opts) *Handler {
	return &Handler{userservice: opts.UserService}
}
