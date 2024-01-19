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
	"time"
)

var roleSrv RoleService

type RoleService interface {
	GetById(ctx context.Context, id uuid.UUID) model.Role
	GetAll(ctx context.Context, page commonModel.Pageable) commonModel.Page[model.Role]
	Create(ctx context.Context, role model.Role) model.Role
	Update(ctx context.Context, role model.Role) model.Role
	DeleteById(ctx context.Context, id uuid.UUID)
	FindByName(ctx context.Context, name string) model.Role
}

type roleService struct {
	repository repository.RoleRepository
	cache      *commonCache.Cache[commonModel.Page[model.Role]]
}

func GetRoleService() RoleService {
	if roleSrv != nil {
		return roleSrv
	}

	roleSrv = (&roleService{
		repository: repository.GetRoleRepository(),
		cache:      commonCache.NewCache[commonModel.Page[model.Role]]("roles", 24*time.Hour),
	}).init()

	return roleSrv
}

func (service *roleService) GetById(ctx context.Context, id uuid.UUID) model.Role {
	return service.repository.GetById(ctx, []uuid.UUID{id})[0]
}

func (service *roleService) GetAll(ctx context.Context, page commonModel.Pageable) commonModel.Page[model.Role] {
	result, err := service.cache.Get(ctx)
	if err == nil {
		return result
	}
	return service.cache.Set(ctx, service.repository.GetAllWithPage(ctx, page))
}

func (service *roleService) Create(ctx context.Context, role model.Role) model.Role {
	defer service.cache.Evict(ctx)
	return service.repository.Create(ctx, []model.Role{role})[0]
}

func (service *roleService) Update(ctx context.Context, role model.Role) model.Role {
	defer service.cache.Evict(ctx)
	return service.repository.Update(ctx, []model.Role{role})[0]
}

func (service *roleService) DeleteById(ctx context.Context, id uuid.UUID) {
	defer service.cache.Evict(ctx)
	service.repository.DeleteById(ctx, id)
}

func (service *roleService) FindByName(ctx context.Context, name string) model.Role {
	result := service.repository.FindByName(ctx, name)
	if len(result) == 0 {
		return model.NilRole
	}
	return result[0]
}

func (service *roleService) init() *roleService {
	ctx := context.Background()
	model.AdminRole = service.createIfNotExists(ctx, model.AdminRole)
	model.InternalUserRole = service.createIfNotExists(ctx, model.InternalUserRole)
	return service
}

func (service *roleService) createIfNotExists(ctx context.Context, role model.Role) model.Role {
	roleFromDB := service.FindByName(ctx, role.Name)
	if commonUtil.IsZeroObject(roleFromDB.ID) {
		return service.Create(ctx, role)
	}

	return service.repository.Update(ctx, []model.Role{roleFromDB})[0]
}
