package repo

import (
	"local/model"
	"local/util/logger"
	"log"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepo interface {
	QueryOne(reqCtx *model.RequestContext, user *model.User) model.Response[*model.User]
	QueryMany(reqCtx *model.RequestContext, us *model.User) model.Response[[]*model.User]
	Create(reqCtx *model.RequestContext, us *model.User) model.Response[*model.User]
	Update(reqCtx *model.RequestContext, us *model.User) model.Response[*model.User]
	Delete(reqCtx *model.RequestContext, us *model.User) model.Response[*model.User]
}

func (r *userRepository) QueryOne(reqCtx *model.RequestContext, user *model.User) model.Response[*model.User] {
	logger.Info(reqCtx, "UserRepo.QueryOne called", map[string]interface{}{"user_id": user.ID, "username": user.UserName})
	var result model.User
	err := r.db.WithContext(reqCtx.Context()).Where(user).First(&result).Error
	if err != nil {
		log.Printf("Error querying one user: %v", err)
		return model.NotFound[*model.User]("User not found")
	}
	return model.SuccessResponse(&result, "User retrieved successfully")
}

func (r *userRepository) QueryMany(reqCtx *model.RequestContext, user *model.User) model.Response[[]*model.User] {
	logger.Info(reqCtx, "UserRepo.QueryMany called")
	var results []*model.User
	err := r.db.WithContext(reqCtx.Context()).Where(user).Find(&results).Error
	if err != nil {
		log.Printf("Error querying many users: %v", err)
		return model.InternalError[[]*model.User]("Failed to query users")
	}
	return model.SuccessResponse(results, "Users retrieved successfully")
}

func (r *userRepository) Create(reqCtx *model.RequestContext, user *model.User) model.Response[*model.User] {
	logger.Info(reqCtx, "UserRepo.Create called", map[string]interface{}{"username": user.UserName})
	err := r.db.WithContext(reqCtx.Context()).Create(user).Error
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return model.BadRequest[*model.User]("Failed to create user")
	}
	return model.SuccessResponse(user, "User created successfully")
}

func (r *userRepository) Update(reqCtx *model.RequestContext, user *model.User) model.Response[*model.User] {
	logger.Info(reqCtx, "UserRepo.Update called", map[string]interface{}{"user_id": user.ID})
	err := r.db.WithContext(reqCtx.Context()).Save(user).Error
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return model.BadRequest[*model.User]("Failed to update user")
	}
	return model.SuccessResponse(user, "User updated successfully")
}

func (r *userRepository) Delete(reqCtx *model.RequestContext, user *model.User) model.Response[*model.User] {
	logger.Info(reqCtx, "UserRepo.Delete called", map[string]interface{}{"user_id": user.ID})
	err := r.db.WithContext(reqCtx.Context()).Delete(user).Error
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return model.BadRequest[*model.User]("Failed to delete user")
	}
	return model.SuccessResponse(user, "User deleted successfully")
}

