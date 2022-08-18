package app_interface

import "team3-task/internal/entity"

type AuthAccessChecker interface {
	CheckAccess(*entity.AuthRequest) (entity.AuthResponse, error)
}
