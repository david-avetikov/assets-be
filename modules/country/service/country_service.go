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
	"assets/modules/country/model"
	"assets/modules/country/repository"
	"context"
	"github.com/google/uuid"
	"time"
)

var countrySrv CountryService

type CountryService interface {
	GetById(ctx context.Context, id uuid.UUID) model.Country
	GetAll(ctx context.Context, page commonModel.Pageable) commonModel.Page[model.Country]
	Create(ctx context.Context, country model.Country) model.Country
	Update(ctx context.Context, country model.Country) model.Country
	DeleteById(ctx context.Context, id uuid.UUID)
	FindByCodeIso(ctx context.Context, name string) model.Country
}

type countryService struct {
	repository repository.CountryRepository
	cache      *commonCache.Cache[commonModel.Page[model.Country]]
}

func GetCountryService() CountryService {
	if countrySrv != nil {
		return countrySrv
	}

	countrySrv = (&countryService{
		repository: repository.GetCountryRepository(),
		cache:      commonCache.NewCache[commonModel.Page[model.Country]]("countrys", 24*time.Hour),
	}).init()

	return countrySrv
}

func (service *countryService) GetById(ctx context.Context, id uuid.UUID) model.Country {
	return service.repository.GetById(ctx, []uuid.UUID{id})[0]
}

func (service *countryService) GetAll(ctx context.Context, page commonModel.Pageable) commonModel.Page[model.Country] {
	result, err := service.cache.Get(ctx)
	if err == nil {
		return result
	}
	return service.cache.Set(ctx, service.repository.GetAllWithPage(ctx, page))
}

func (service *countryService) Create(ctx context.Context, country model.Country) model.Country {
	defer service.cache.Evict(ctx)
	return service.repository.Create(ctx, []model.Country{country})[0]
}

func (service *countryService) Update(ctx context.Context, country model.Country) model.Country {
	defer service.cache.Evict(ctx)
	return service.repository.Update(ctx, []model.Country{country})[0]
}

func (service *countryService) DeleteById(ctx context.Context, id uuid.UUID) {
	defer service.cache.Evict(ctx)
	service.repository.DeleteById(ctx, id)
}

func (service *countryService) FindByCodeIso(ctx context.Context, name string) model.Country {
	result := service.repository.FindByCodeIso(ctx, name)
	if len(result) == 0 {
		return model.NilCountry
	}
	return result[0]
}

func (service *countryService) init() *countryService {
	ctx := context.Background()
	model.RussiaCountry = service.createIfNotExists(ctx, model.RussiaCountry)
	return service
}

func (service *countryService) createIfNotExists(ctx context.Context, city model.Country) model.Country {
	countryFromDB := service.FindByCodeIso(ctx, city.CodeIso)
	if commonUtil.IsZeroObject(countryFromDB.ID) {
		return service.Create(ctx, city)
	}
	return countryFromDB
}
