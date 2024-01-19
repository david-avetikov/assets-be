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

var NilCountry = Country{}

type Country struct {
	ID        uuid.UUID `json:"id,omitempty" gorm:"type:uuid;primaryKey"`
	CodeIso   string    `json:"codeIso,omitempty"`
	ShortName string    `json:"shortName,omitempty"`
	FullName  string    `json:"fullName,omitempty"`
	CharCode2 string    `json:"charCode2,omitempty"`
	CharCode3 string    `json:"charCode3,omitempty"`
	PhoneCode string    `json:"phoneCode,omitempty"`
	PhoneMask string    `json:"phoneMask,omitempty"`
	Pictogram []byte    `json:"pictogram,omitempty"`
	Cities    []*City   `json:"cities,omitempty"`
	Currency  *Currency `json:"currency,omitempty"`
}

func (country Country) GetID() uuid.UUID {
	return country.ID
}

func (country *Country) BeforeCreate(tx *gorm.DB) error {
	if commonUtil.IsZeroObject(country.ID) {
		country.ID = uuid.New()
	}
	return nil
}

func NewCountry(codeIso string, shortName string, fullName string, opts ...CountryOption) Country {
	country := Country{
		CodeIso:   codeIso,
		ShortName: shortName,
		FullName:  fullName,
	}

	for _, opt := range opts {
		opt(&country)
	}

	return country
}

type CountryOption func(*Country)

func CountryWithCities(cities ...*City) CountryOption {
	return func(role *Country) {
		role.Cities = append(role.Cities, cities...)
	}
}

func CountryWithCurrency(currency *Currency) CountryOption {
	return func(role *Country) {
		role.Currency = currency
	}
}
