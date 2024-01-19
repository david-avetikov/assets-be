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

var currencySrv CurrencyService

type CurrencyService interface {
	GetById(ctx context.Context, id uuid.UUID) model.Currency
	GetAll(ctx context.Context, page commonModel.Pageable) commonModel.Page[model.Currency]
	Create(ctx context.Context, currency model.Currency) model.Currency
	Update(ctx context.Context, currency model.Currency) model.Currency
	DeleteById(ctx context.Context, id uuid.UUID)
	FindByAlphabeticCode(ctx context.Context, alphabeticCode string) model.Currency
}

type currencyService struct {
	repository repository.CurrencyRepository
	cache      *commonCache.Cache[commonModel.Page[model.Currency]]
}

func GetCurrencyService() CurrencyService {
	if currencySrv != nil {
		return currencySrv
	}

	currencySrv = (&currencyService{
		repository: repository.GetCurrencyRepository(),
		cache:      commonCache.NewCache[commonModel.Page[model.Currency]]("currencys", 24*time.Hour),
	}).init()

	return currencySrv
}

func (service *currencyService) GetById(ctx context.Context, id uuid.UUID) model.Currency {
	return service.repository.GetById(ctx, []uuid.UUID{id})[0]
}

func (service *currencyService) GetAll(ctx context.Context, page commonModel.Pageable) commonModel.Page[model.Currency] {
	result, err := service.cache.Get(ctx)
	if err == nil {
		return result
	}
	return service.cache.Set(ctx, service.repository.GetAllWithPage(ctx, page))
}

func (service *currencyService) Create(ctx context.Context, currency model.Currency) model.Currency {
	defer service.cache.Evict(ctx)
	return service.repository.Create(ctx, []model.Currency{currency})[0]
}

func (service *currencyService) Update(ctx context.Context, currency model.Currency) model.Currency {
	defer service.cache.Evict(ctx)
	return service.repository.Update(ctx, []model.Currency{currency})[0]
}

func (service *currencyService) DeleteById(ctx context.Context, id uuid.UUID) {
	defer service.cache.Evict(ctx)
	service.repository.DeleteById(ctx, id)
}

func (service *currencyService) FindByAlphabeticCode(ctx context.Context, alphabeticCode string) model.Currency {
	result := service.repository.FindByAlphabeticCode(ctx, alphabeticCode)
	if len(result) == 0 {
		return model.NilCurrency
	}
	return result[0]
}

func (service *currencyService) init() *currencyService {
	ctx := context.Background()
	model.RubCurrency = service.createIfNotExists(ctx, model.RubCurrency)
	return service
}

func (service *currencyService) createIfNotExists(ctx context.Context, currency model.Currency) model.Currency {
	currencyFromDB := service.FindByAlphabeticCode(ctx, currency.AlphabeticCode)
	if commonUtil.IsZeroObject(currencyFromDB.ID) {
		return service.Create(ctx, currency)
	}
	return currencyFromDB
}
