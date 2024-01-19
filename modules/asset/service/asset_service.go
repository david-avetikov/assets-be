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
	"assets/modules/asset/model"
	"assets/modules/asset/repository"
	"context"
	"github.com/google/uuid"
	"time"
)

var assetSrv AssetService

type AssetService interface {
	GetById(ctx context.Context, id uuid.UUID) model.Asset
	GetAll(ctx context.Context, page commonModel.Pageable) commonModel.Page[model.Asset]
	Create(ctx context.Context, asset model.Asset) model.Asset
	Update(ctx context.Context, asset model.Asset) model.Asset
	DeleteById(ctx context.Context, id uuid.UUID)
}

type assetService struct {
	repository repository.AssetRepository
	cache      *commonCache.Cache[commonModel.Page[model.Asset]]
}

func GetAssetService() AssetService {
	if assetSrv != nil {
		return assetSrv
	}

	assetSrv = &assetService{
		repository: repository.GetAssetRepository(),
		cache:      commonCache.NewCache[commonModel.Page[model.Asset]]("assets", 24*time.Hour),
	}

	return assetSrv
}

func (service *assetService) GetById(ctx context.Context, id uuid.UUID) model.Asset {
	return service.repository.GetById(ctx, []uuid.UUID{id})[0]
}

func (service *assetService) GetAll(ctx context.Context, page commonModel.Pageable) commonModel.Page[model.Asset] {
	result, err := service.cache.Get(ctx)
	if err == nil {
		return result
	}
	return service.cache.Set(ctx, service.repository.GetAllWithPage(ctx, page))
}

func (service *assetService) Create(ctx context.Context, asset model.Asset) model.Asset {
	defer service.cache.Evict(ctx)
	return service.repository.Create(ctx, []model.Asset{asset})[0]
}

func (service *assetService) Update(ctx context.Context, asset model.Asset) model.Asset {
	defer service.cache.Evict(ctx)
	return service.repository.Update(ctx, []model.Asset{asset})[0]
}

func (service *assetService) DeleteById(ctx context.Context, id uuid.UUID) {
	defer service.cache.Evict(ctx)
	service.repository.DeleteById(ctx, id)
}
