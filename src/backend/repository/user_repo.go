package repository

import (
	"context"
	"local/model"
	"log"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepo interface {
	QueryOne(ctx context.Context, user *model.User) *model.User
	QueryMany(ctx context.Context, us *model.User) []*model.User
	Create(ctx context.Context, us *model.User) *model.User
	Update(ctx context.Context, us *model.User) *model.User
	Delete(ctx context.Context, us *model.User) *model.User
}

func (r *userRepository) QueryOne(ctx context.Context, user *model.User) *model.User {
	var result model.User
	err := r.db.WithContext(ctx).Where(user).First(&result).Error
	if err != nil {
		log.Printf("Error querying one user: %v", err)
		return nil
	}
	return &result
}

func (r *userRepository) QueryMany(ctx context.Context, user *model.User) []*model.User {
	var results []*model.User
	err := r.db.WithContext(ctx).Where(user).Find(&results).Error
	if err != nil {
		log.Printf("Error querying many users: %v", err)
		return nil
	}
	return results
}

func (r *userRepository) Create(ctx context.Context, user *model.User) *model.User {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil
	}
	return user
}

func (r *userRepository) Update(ctx context.Context, user *model.User) *model.User {
	err := r.db.WithContext(ctx).Save(user).Error
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return nil
	}
	return user
}

func (r *userRepository) Delete(ctx context.Context, user *model.User) *model.User {
	err := r.db.WithContext(ctx).Delete(user).Error
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return nil
	}
	return user
}
