package model

import (
	commonUtil "assets/common/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

var NilCurrency = Currency{}

type Currency struct {
	ID             uuid.UUID `json:"id,omitempty" gorm:"type:uuid;primaryKey"`
	CountryId      uuid.UUID `json:"countryId,omitempty"`
	AlphabeticCode string    `json:"alphabeticCode,omitempty"`
	NumericCode    string    `json:"numericCode,omitempty"`
	Symbol         string    `json:"symbol,omitempty"`
}

func (currency Currency) GetID() uuid.UUID {
	return currency.ID
}

func (currency *Currency) BeforeCreate(tx *gorm.DB) error {
	if commonUtil.IsZeroObject(currency.ID) {
		currency.ID = uuid.New()
	}
	return nil
}

func NewCurrency(alphabeticCode string, numericCode string, symbol string) Currency {
	return Currency{
		AlphabeticCode: alphabeticCode,
		NumericCode:    numericCode,
		Symbol:         symbol,
	}
}
