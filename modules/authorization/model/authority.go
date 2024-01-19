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

var NilAuthority = Authority{}

type Authority struct {
	ID          uuid.UUID `json:"id,omitempty" gorm:"type:uuid;primaryKey"`
	Method      string    `json:"method,omitempty"`
	Description string    `json:"description,omitempty"`
}

func (authority Authority) GetID() uuid.UUID {
	return authority.ID
}

func (authority *Authority) BeforeCreate(tx *gorm.DB) error {
	if commonUtil.IsZeroObject(authority.ID) {
		authority.ID = uuid.New()
	}
	return nil
}

func NewAuthority(method string, description string) Authority {
	return Authority{
		Method:      method,
		Description: description,
	}
}
