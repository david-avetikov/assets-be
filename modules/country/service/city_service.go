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

var citySrv CityService

type CityService interface {
	GetById(ctx context.Context, id uuid.UUID) model.City
	GetAll(ctx context.Context, page commonModel.Pageable) commonModel.Page[model.City]
	Create(ctx context.Context, city model.City) model.City
	Update(ctx context.Context, city model.City) model.City
	DeleteById(ctx context.Context, id uuid.UUID)
	FindByName(ctx context.Context, name string) model.City
}

type cityService struct {
	repository repository.CityRepository
	cache      *commonCache.Cache[commonModel.Page[model.City]]
}

func GetCityService() CityService {
	if citySrv != nil {
		return citySrv
	}

	citySrv = (&cityService{
		repository: repository.GetCityRepository(),
		cache:      commonCache.NewCache[commonModel.Page[model.City]]("cities", 24*time.Hour),
	}).init()

	return citySrv
}

func (service *cityService) GetById(ctx context.Context, id uuid.UUID) model.City {
	return service.repository.GetById(ctx, []uuid.UUID{id})[0]
}

func (service *cityService) GetAll(ctx context.Context, page commonModel.Pageable) commonModel.Page[model.City] {
	result, err := service.cache.Get(ctx)
	if err == nil {
		return result
	}
	return service.cache.Set(ctx, service.repository.GetAllWithPage(ctx, page))
}

func (service *cityService) Create(ctx context.Context, city model.City) model.City {
	defer service.cache.Evict(ctx)
	return service.repository.Create(ctx, []model.City{city})[0]
}

func (service *cityService) Update(ctx context.Context, city model.City) model.City {
	defer service.cache.Evict(ctx)
	return service.repository.Update(ctx, []model.City{city})[0]
}

func (service *cityService) DeleteById(ctx context.Context, id uuid.UUID) {
	defer service.cache.Evict(ctx)
	service.repository.DeleteById(ctx, id)
}

func (service *cityService) FindByName(ctx context.Context, name string) model.City {
	result := service.repository.FindByName(ctx, name)
	if len(result) == 0 {
		return model.NilCity
	}
	return result[0]
}

func (service *cityService) init() *cityService {
	ctx := context.Background()
	model.MoscowCity = service.createIfNotExists(ctx, model.MoscowCity)
	return service
}

func (service *cityService) createIfNotExists(ctx context.Context, city model.City) model.City {
	cityFromDB := service.FindByName(ctx, city.Name)
	if commonUtil.IsZeroObject(cityFromDB.ID) {
		return service.Create(ctx, city)
	}
	return cityFromDB
}
