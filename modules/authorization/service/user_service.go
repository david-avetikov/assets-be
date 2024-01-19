package service

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
	commonCache "assets/common/cache"
	commonModel "assets/common/model"
	commonUtil "assets/common/util"
	"assets/modules/authorization/model"
	"assets/modules/authorization/repository"
	"context"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

var userSrv UserService

type UserService interface {
	GetById(ctx context.Context, id uuid.UUID) model.User
	GetAll(ctx context.Context, page commonModel.Pageable) commonModel.Page[model.User]
	Create(ctx context.Context, user model.User) model.User
	Update(ctx context.Context, user model.User) model.User
	DeleteById(ctx context.Context, id uuid.UUID)
	FindByUsername(ctx context.Context, username string) model.User
	FindByEmail(ctx context.Context, email string) model.User

	GetCurrent(ctx context.Context) model.User
	AddRoles(ctx context.Context, id uuid.UUID, rolesIds []uuid.UUID) model.User
	RemoveRoles(ctx context.Context, id uuid.UUID, rolesIds []uuid.UUID) model.User
	AddAuthorities(ctx context.Context, id uuid.UUID, authoritiesIds []uuid.UUID) model.User
	RemoveAuthorities(ctx context.Context, id uuid.UUID, authoritiesIds []uuid.UUID) model.User

	Authorize(ctx context.Context, username string, password string) bool
}

type userService struct {
	repository          repository.UserRepository
	roleRepository      repository.RoleRepository
	authorityRepository repository.AuthorityRepository
	cache               *commonCache.Cache[commonModel.Page[model.User]]
}

func GetUserService() UserService {
	if userSrv != nil {
		return userSrv
	}

	userSrv = (&userService{
		repository:          repository.GetUserRepository(),
		roleRepository:      repository.GetRoleRepository(),
		authorityRepository: repository.GetAuthorityRepository(),
		cache:               commonCache.NewCache[commonModel.Page[model.User]]("users", 24*time.Hour),
	}).init()

	return userSrv
}

func (service *userService) GetById(ctx context.Context, id uuid.UUID) model.User {
	return service.repository.GetById(ctx, []uuid.UUID{id})[0]
}

func (service *userService) GetAll(ctx context.Context, page commonModel.Pageable) commonModel.Page[model.User] {
	result, err := service.cache.Get(ctx)
	if err == nil {
		return result
	}
	return service.cache.Set(ctx, service.repository.GetAllWithPage(ctx, page))
}

func (service *userService) Create(ctx context.Context, user model.User) model.User {
	defer service.cache.Evict(ctx)
	return service.repository.Create(ctx, []model.User{user})[0]
}

func (service *userService) Update(ctx context.Context, user model.User) model.User {
	defer service.cache.Evict(ctx)
	return service.repository.Update(ctx, []model.User{user})[0]
}

func (service *userService) DeleteById(ctx context.Context, id uuid.UUID) {
	defer service.cache.Evict(ctx)
	service.repository.DeleteById(ctx, id)
}

func (service *userService) FindByUsername(ctx context.Context, username string) model.User {
	result := service.repository.FindByUsername(ctx, []string{strings.ToLower(username)})
	if len(result) == 0 {
		return model.NilUser
	}
	return result[0]
}

func (service *userService) FindByEmail(ctx context.Context, email string) model.User {
	result := service.repository.FindByEmail(ctx, []string{strings.ToLower(email)})
	if len(result) == 0 {
		return model.NilUser
	}
	return result[0]
}

func (service *userService) GetCurrent(ctx context.Context) model.User {
	tokenInfo := commonUtil.MustGetCurrentTokenInfo(ctx)
	return service.repository.GetById(ctx, []uuid.UUID{tokenInfo.UserId})[0]
}

func (service *userService) AddRoles(ctx context.Context, id uuid.UUID, rolesIds []uuid.UUID) model.User {
	user := service.GetById(ctx, id)
	roles := service.roleRepository.GetById(ctx, rolesIds)
	user.Roles = append(user.Roles, commonUtil.Map(roles, func(it model.Role) *model.Role { return &it })...)
	return service.Update(ctx, user)
}

func (service *userService) RemoveRoles(ctx context.Context, id uuid.UUID, rolesIds []uuid.UUID) model.User {
	return service.repository.RemoveRoles(ctx, id, service.roleRepository.GetById(ctx, rolesIds))
}

func (service *userService) AddAuthorities(ctx context.Context, id uuid.UUID, authoritiesIds []uuid.UUID) model.User {
	user := service.GetById(ctx, id)
	authorities := service.authorityRepository.GetById(ctx, authoritiesIds)
	user.AdditionalAuthorities = append(user.AdditionalAuthorities, commonUtil.Map(authorities, func(it model.Authority) *model.Authority { return &it })...)
	return service.Update(ctx, user)
}

func (service *userService) RemoveAuthorities(ctx context.Context, id uuid.UUID, authoritiesIds []uuid.UUID) model.User {
	return service.repository.RemoveAuthorities(ctx, id, service.authorityRepository.GetById(ctx, authoritiesIds))
}

func (service *userService) Authorize(ctx context.Context, username string, password string) bool {
	users := service.repository.FindByUsername(ctx, []string{username})
	if len(users) == 0 {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(password)) == nil
}

func (service *userService) init() *userService {
	ctx := context.Background()
	assetss := service.repository.FindByUsername(ctx, []string{"assets"})
	if len(assetss) == 0 {
		assets := service.Create(ctx, model.NewUser(
			"assets",
			"assets",
			"",
			"admin@deadline.team",
			model.UserWithPassword("Deadline@777"),
		))
		service.AddRoles(ctx, assets.ID, []uuid.UUID{model.AdminRole.ID})
	}
	return service
}

func findUserInArray(users []model.User, user model.User) *model.User {
	for _, usr := range users {
		if usr.Username == user.Username {
			return &usr
		}
	}
	return nil
}
