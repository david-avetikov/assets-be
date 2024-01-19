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

var countryRepo CountryRepository

type CountryRepository interface {
	commonRepository.Repository[model.Country]
	FindByCodeIso(ctx context.Context, CodeIsos ...string) []model.Country
}

type countryRepository struct {
	commonRepository.Repository[model.Country]
	*commonDB.DataSource
}

func GetCountryRepository() CountryRepository {
	if countryRepo != nil {
		return countryRepo
	}
	countryRepo = &countryRepository{
		commonRepository.NewBaseRepository[model.Country](commonDB.GetDataSource()),
		commonDB.GetDataSource(),
	}
	return countryRepo
}

func (repo *countryRepository) FindByCodeIso(ctx context.Context, CodeIsos ...string) []model.Country {
	if len(CodeIsos) == 0 {
		panic(commonError.IllegalArgumentError)
	}

	var result []model.Country
	repo.DataSource.Where("code_iso in ?", CodeIsos).Find(&result)
	return result
}
