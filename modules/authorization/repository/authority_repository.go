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
	"assets/modules/authorization/model"
	"context"
)

var authorityRepo AuthorityRepository

type AuthorityRepository interface {
	commonRepository.Repository[model.Authority]
	FindByMethod(ctx context.Context, methods ...string) []model.Authority
}

type authorityRepository struct {
	commonRepository.Repository[model.Authority]
	*commonDB.DataSource
}

func GetAuthorityRepository() AuthorityRepository {
	if authorityRepo != nil {
		return authorityRepo
	}
	authorityRepo = &authorityRepository{
		commonRepository.NewBaseRepository[model.Authority](commonDB.GetDataSource()),
		commonDB.GetDataSource(),
	}
	return authorityRepo
}

func (repo *authorityRepository) FindByMethod(ctx context.Context, methods ...string) []model.Authority {
	if len(methods) == 0 {
		panic(commonError.IllegalArgumentError)
	}

	var result []model.Authority
	repo.DataSource.Where("method in ?", methods).Find(&result)
	return result
}
