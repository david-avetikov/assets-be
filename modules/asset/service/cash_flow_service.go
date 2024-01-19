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

var cashFlowSrv CashFlowService

type CashFlowService interface {
	GetById(ctx context.Context, id uuid.UUID) model.CashFlow
	GetAll(ctx context.Context, page commonModel.Pageable) commonModel.Page[model.CashFlow]
	Create(ctx context.Context, cashFlow model.CashFlow) model.CashFlow
	Update(ctx context.Context, cashFlow model.CashFlow) model.CashFlow
	DeleteById(ctx context.Context, id uuid.UUID)
}

type cashFlowService struct {
	repository repository.CashFlowRepository
	cache      *commonCache.Cache[commonModel.Page[model.CashFlow]]
}

func GetCashFlowService() CashFlowService {
	if cashFlowSrv != nil {
		return cashFlowSrv
	}

	cashFlowSrv = &cashFlowService{
		repository: repository.GetCashFlowRepository(),
		cache:      commonCache.NewCache[commonModel.Page[model.CashFlow]]("cashFlows", 24*time.Hour),
	}

	return cashFlowSrv
}

func (service *cashFlowService) GetById(ctx context.Context, id uuid.UUID) model.CashFlow {
	return service.repository.GetById(ctx, []uuid.UUID{id})[0]
}

func (service *cashFlowService) GetAll(ctx context.Context, page commonModel.Pageable) commonModel.Page[model.CashFlow] {
	result, err := service.cache.Get(ctx)
	if err == nil {
		return result
	}
	return service.cache.Set(ctx, service.repository.GetAllWithPage(ctx, page))
}

func (service *cashFlowService) Create(ctx context.Context, cashFlow model.CashFlow) model.CashFlow {
	defer service.cache.Evict(ctx)
	return service.repository.Create(ctx, []model.CashFlow{cashFlow})[0]
}

func (service *cashFlowService) Update(ctx context.Context, cashFlow model.CashFlow) model.CashFlow {
	defer service.cache.Evict(ctx)
	return service.repository.Update(ctx, []model.CashFlow{cashFlow})[0]
}

func (service *cashFlowService) DeleteById(ctx context.Context, id uuid.UUID) {
	defer service.cache.Evict(ctx)
	service.repository.DeleteById(ctx, id)
}
