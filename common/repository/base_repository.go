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
	"assets/common/custom_error"
	"assets/common/db"
	"assets/common/iface"
	"assets/common/model"
	"assets/common/util"
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type Repository[T iface.Entity] interface {
	GetById(ctx context.Context, ids []uuid.UUID) []T
	GetAll(ctx context.Context) []T
	GetAllWithPage(ctx context.Context, page model.Pageable) model.Page[T]
	Create(ctx context.Context, entities []T) []T
	Update(ctx context.Context, entities []T) []T
	DeleteById(ctx context.Context, ids ...uuid.UUID)
}

type baseRepository[T iface.Entity] struct {
	ds *db.DataSource
}

func NewBaseRepository[T iface.Entity](dataSource *db.DataSource) Repository[T] {
	var entity T
	util.Must(dataSource.AutoMigrate(&entity))
	return &baseRepository[T]{ds: dataSource}
}

func (baseRepo *baseRepository[T]) GetById(ctx context.Context, ids []uuid.UUID) []T {
	if len(ids) == 0 {
		panic(custom_error.IllegalArgumentError)
	}

	var result []T
	util.Must(baseRepo.ds.Preload(clause.Associations).Where("id in ?", util.Map(ids, func(it uuid.UUID) string { return it.String() })).Find(&result).Error)
	return result
}

func (baseRepo *baseRepository[T]) GetAll(ctx context.Context) []T {
	var result []T
	util.Must(baseRepo.ds.Preload(clause.Associations).Find(&result).Error)
	return result
}

func (baseRepo *baseRepository[T]) GetAllWithPage(ctx context.Context, page model.Pageable) model.Page[T] {
	var result []T
	util.Must(baseRepo.ds.Preload(clause.Associations).Limit(page.Size).Offset(page.Page * page.Size).Order(fmt.Sprintf("%s %s", page.Sort.Field, page.Sort.Order)).Find(&result).Error)
	return model.NewPage[T](result, page)
}

func (baseRepo *baseRepository[T]) Create(ctx context.Context, entities []T) []T {
	util.Must(baseRepo.ds.Preload(clause.Associations).Create(&entities).Error)
	return entities
}

func (baseRepo *baseRepository[T]) Update(ctx context.Context, entities []T) []T {
	util.Must(baseRepo.ds.Preload(clause.Associations).Save(&entities).Error)
	return entities
}

func (baseRepo *baseRepository[T]) DeleteById(ctx context.Context, ids ...uuid.UUID) {
	if len(ids) == 0 {
		panic(custom_error.IllegalArgumentError)
	}

	var entity []T
	util.Must(baseRepo.ds.Where("id in ?", util.Map(ids, func(it uuid.UUID) string { return it.String() })).Delete(&entity).Error)
}
