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

var authoritySrv AuthorityService

type AuthorityService interface {
	GetById(ctx context.Context, id uuid.UUID) model.Authority
	GetAll(ctx context.Context, page commonModel.Pageable) commonModel.Page[model.Authority]
	Create(ctx context.Context, authority model.Authority) model.Authority
	Update(ctx context.Context, authority model.Authority) model.Authority
	DeleteById(ctx context.Context, id uuid.UUID)
	FindByMethod(ctx context.Context, method string) model.Authority
}

type authorityService struct {
	repository repository.AuthorityRepository
	cache      *commonCache.Cache[commonModel.Page[model.Authority]]
}

func GetAuthorityService() AuthorityService {
	if authoritySrv != nil {
		return authoritySrv
	}

	authoritySrv = (&authorityService{
		repository: repository.GetAuthorityRepository(),
		cache:      commonCache.NewCache[commonModel.Page[model.Authority]]("authorities", 24*time.Hour),
	}).init()

	return authoritySrv
}

func (service *authorityService) GetById(ctx context.Context, id uuid.UUID) model.Authority {
	return service.repository.GetById(ctx, []uuid.UUID{id})[0]
}

func (service *authorityService) GetAll(ctx context.Context, page commonModel.Pageable) commonModel.Page[model.Authority] {
	result, err := service.cache.Get(ctx)
	if err == nil {
		return result
	}
	return service.cache.Set(ctx, service.repository.GetAllWithPage(ctx, page))
}

func (service *authorityService) Create(ctx context.Context, authority model.Authority) model.Authority {
	defer service.cache.Evict(ctx)
	return service.repository.Create(ctx, []model.Authority{authority})[0]
}

func (service *authorityService) Update(ctx context.Context, authority model.Authority) model.Authority {
	defer service.cache.Evict(ctx)
	return service.repository.Update(ctx, []model.Authority{authority})[0]
}

func (service *authorityService) DeleteById(ctx context.Context, id uuid.UUID) {
	defer service.cache.Evict(ctx)
	service.repository.DeleteById(ctx, id)
}

func (service *authorityService) FindByMethod(ctx context.Context, method string) model.Authority {
	result := service.repository.FindByMethod(ctx, method)
	if len(result) == 0 {
		return model.NilAuthority
	}
	return result[0]
}

func (service *authorityService) init() *authorityService {
	ctx := context.Background()
	model.OwnerAuthority = service.createIfNotExists(ctx, model.OwnerAuthority)

	model.ReadAssetAuthority = service.createIfNotExists(ctx, model.ReadAssetAuthority)
	model.CreateAssetAuthority = service.createIfNotExists(ctx, model.CreateAssetAuthority)
	model.UpdateAssetAuthority = service.createIfNotExists(ctx, model.UpdateAssetAuthority)
	model.DeleteAssetAuthority = service.createIfNotExists(ctx, model.DeleteAssetAuthority)

	model.ReadCashFlowAuthority = service.createIfNotExists(ctx, model.ReadCashFlowAuthority)
	model.CreateCashFlowAuthority = service.createIfNotExists(ctx, model.CreateCashFlowAuthority)
	model.UpdateCashFlowAuthority = service.createIfNotExists(ctx, model.UpdateCashFlowAuthority)
	model.DeleteCashFlowAuthority = service.createIfNotExists(ctx, model.DeleteCashFlowAuthority)

	model.CreateAttachmentAuthority = service.createIfNotExists(ctx, model.CreateAttachmentAuthority)
	model.DeleteAttachmentAuthority = service.createIfNotExists(ctx, model.DeleteAttachmentAuthority)

	model.ReadAuthorityAuthority = service.createIfNotExists(ctx, model.ReadAuthorityAuthority)
	model.CreateAuthorityAuthority = service.createIfNotExists(ctx, model.CreateAuthorityAuthority)
	model.UpdateAuthorityAuthority = service.createIfNotExists(ctx, model.UpdateAuthorityAuthority)
	model.DeleteAuthorityAuthority = service.createIfNotExists(ctx, model.DeleteAuthorityAuthority)

	model.ReadRoleAuthority = service.createIfNotExists(ctx, model.ReadRoleAuthority)
	model.CreateRoleAuthority = service.createIfNotExists(ctx, model.CreateRoleAuthority)
	model.UpdateRoleAuthority = service.createIfNotExists(ctx, model.UpdateRoleAuthority)
	model.DeleteRoleAuthority = service.createIfNotExists(ctx, model.DeleteRoleAuthority)

	model.ReadUserAuthority = service.createIfNotExists(ctx, model.ReadUserAuthority)
	model.CreateUserAuthority = service.createIfNotExists(ctx, model.CreateUserAuthority)
	model.UpdateUserAuthority = service.createIfNotExists(ctx, model.UpdateUserAuthority)
	model.DeleteUserAuthority = service.createIfNotExists(ctx, model.DeleteUserAuthority)
	model.EditRoleUserAuthority = service.createIfNotExists(ctx, model.EditRoleUserAuthority)
	model.EditAuthorityUserAuthority = service.createIfNotExists(ctx, model.EditAuthorityUserAuthority)

	model.ReadCityAuthority = service.createIfNotExists(ctx, model.ReadCityAuthority)
	model.CreateCityAuthority = service.createIfNotExists(ctx, model.CreateCityAuthority)
	model.UpdateCityAuthority = service.createIfNotExists(ctx, model.UpdateCityAuthority)
	model.DeleteCityAuthority = service.createIfNotExists(ctx, model.DeleteCityAuthority)

	model.ReadCountryAuthority = service.createIfNotExists(ctx, model.ReadCountryAuthority)
	model.CreateCountryAuthority = service.createIfNotExists(ctx, model.CreateCountryAuthority)
	model.UpdateCountryAuthority = service.createIfNotExists(ctx, model.UpdateCountryAuthority)
	model.DeleteCountryAuthority = service.createIfNotExists(ctx, model.DeleteCountryAuthority)

	model.ReadCurrencyAuthority = service.createIfNotExists(ctx, model.ReadCurrencyAuthority)
	model.CreateCurrencyAuthority = service.createIfNotExists(ctx, model.CreateCurrencyAuthority)
	model.UpdateCurrencyAuthority = service.createIfNotExists(ctx, model.UpdateCurrencyAuthority)
	model.DeleteCurrencyAuthority = service.createIfNotExists(ctx, model.DeleteCurrencyAuthority)

	return service
}

func (service *authorityService) createIfNotExists(ctx context.Context, authority model.Authority) model.Authority {
	authorityFromDB := service.FindByMethod(ctx, authority.Method)
	if commonUtil.IsZeroObject(authorityFromDB.ID) {
		return service.Create(ctx, authority)
	}
	return authorityFromDB
}
