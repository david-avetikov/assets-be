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
	"assets/modules/country/model"
	"context"
)

var cityRepo CityRepository

type CityRepository interface {
	commonRepository.Repository[model.City]
	FindByName(ctx context.Context, names ...string) []model.City
}

type cityRepository struct {
	commonRepository.Repository[model.City]
	*commonDB.DataSource
}

func GetCityRepository() CityRepository {
	if cityRepo != nil {
		return cityRepo
	}
	cityRepo = &cityRepository{
		commonRepository.NewBaseRepository[model.City](commonDB.GetDataSource()),
		commonDB.GetDataSource(),
	}
	return cityRepo
}

func (repo *cityRepository) FindByName(ctx context.Context, names ...string) []model.City {
	if len(names) == 0 {
		panic(commonError.IllegalArgumentError)
	}

	var result []model.City
	repo.DataSource.Where("name in ?", names).Find(&result)
	return result
}
