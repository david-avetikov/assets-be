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

var NilRole = Role{}

type Role struct {
	ID          uuid.UUID    `json:"id,omitempty" gorm:"type:uuid;primaryKey"`
	Name        string       `json:"name,omitempty"`
	Description string       `json:"description,omitempty"`
	Authorities []*Authority `json:"authorities,omitempty" gorm:"many2many:role_authority;"`
}

func (role Role) GetID() uuid.UUID {
	return role.ID
}

func (role *Role) BeforeCreate(tx *gorm.DB) error {
	if commonUtil.IsZeroObject(role.ID) {
		role.ID = uuid.New()
	}
	return nil
}

func NewRole(name string, description string, opts ...RoleOption) Role {
	role := Role{
		Name:        name,
		Description: description,
		Authorities: []*Authority{},
	}

	for _, opt := range opts {
		opt(&role)
	}

	return role
}

type RoleOption func(*Role)

func RoleWithAuthorities(authorities ...*Authority) RoleOption {
	return func(role *Role) {
		role.Authorities = append(role.Authorities, authorities...)
	}
}
