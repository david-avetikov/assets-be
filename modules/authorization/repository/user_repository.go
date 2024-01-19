package repository

/*
 * Copyright © 2024, "DEADLINE TEAM" LLC
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are not permitted.
 *
 * THIS SOFTWARE IS PROVIDED BY "DEADLINE TEAM" LLC "AS IS" AND ANY
 * EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL "DEADLINE TEAM" LLC BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * No reproductions or distributions of this code is permitted without
 * written permission from "DEADLINE TEAM" LLC.
 * Do not reverse engineer or modify this code.
 *
 * © "DEADLINE TEAM" LLC, All rights reserved.
 */

import (
	commonError "assets/common/custom_error"
	commonDB "assets/common/db"
	commonRepository "assets/common/repository"
	commonUtil "assets/common/util"
	"assets/modules/authorization/model"
	"context"
	"github.com/google/uuid"
)

var userRepo UserRepository

type UserRepository interface {
	commonRepository.Repository[model.User]
	FindByUsername(ctx context.Context, usernames []string) []model.User
	FindByEmail(ctx context.Context, emails []string) []model.User
	RemoveRoles(ctx context.Context, id uuid.UUID, roles []model.Role) model.User
	RemoveAuthorities(ctx context.Context, id uuid.UUID, authorities []model.Authority) model.User
}

type userRepository struct {
	commonRepository.Repository[model.User]
	*commonDB.DataSource
}

func GetUserRepository() UserRepository {
	if userRepo != nil {
		return userRepo
	}
	userRepo = &userRepository{
		commonRepository.NewBaseRepository[model.User](commonDB.GetDataSource()),
		commonDB.GetDataSource(),
	}
	return userRepo
}

func (repo *userRepository) FindByUsername(ctx context.Context, usernames []string) []model.User {
	if len(usernames) == 0 {
		panic(commonError.IllegalArgumentError)
	}

	var result []model.User
	repo.DataSource.Preload("Roles.Authorities").Preload("AdditionalAuthorities").Where("username in ?", usernames).Find(&result)
	return result
}

func (repo *userRepository) FindByEmail(ctx context.Context, emails []string) []model.User {
	if len(emails) == 0 {
		panic(commonError.IllegalArgumentError)
	}

	var result []model.User
	repo.DataSource.Where("email in ?", emails).Find(&result)
	return result
}

func (repo *userRepository) RemoveRoles(ctx context.Context, id uuid.UUID, roles []model.Role) model.User {
	commonUtil.Must(repo.DataSource.Model(&model.User{ID: id}).Association("Roles").Delete(&roles))
	return repo.GetById(ctx, []uuid.UUID{id})[0]
}

func (repo *userRepository) RemoveAuthorities(ctx context.Context, id uuid.UUID, authorities []model.Authority) model.User {
	commonUtil.Must(repo.DataSource.Model(&model.User{ID: id}).Association("Authorities").Delete(&authorities))
	return repo.GetById(ctx, []uuid.UUID{id})[0]
}
