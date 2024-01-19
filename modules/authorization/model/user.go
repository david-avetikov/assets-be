package model

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
	commonUtil "assets/common/util"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
)

var NilUser = User{}

type User struct {
	ID                    uuid.UUID    `json:"id,omitempty" gorm:"type:uuid;primaryKey"`
	Username              string       `json:"username,omitempty"`
	Password              string       `json:"-"`
	FirstName             string       `json:"firstName,omitempty"`
	LastName              string       `json:"lastName,omitempty"`
	Email                 string       `json:"email,omitempty"`
	PhoneNumber           int64        `json:"phoneNumber,omitempty"`
	Roles                 []*Role      `json:"roles,omitempty" gorm:"many2many:user_role;"`
	AdditionalAuthorities []*Authority `json:"authorities,omitempty" gorm:"many2many:user_additional_authority;"`
	IsBlocked             bool         `json:"isBlocked,omitempty"`
}

func (user User) GetID() uuid.UUID {
	return user.ID
}

func (user User) GetFullName() string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	if commonUtil.IsZeroObject(user.ID) {
		user.ID = uuid.New()
	}
	return nil
}

func NewUser(username string, firstName string, lastName string, email string, opts ...UserOption) User {
	user := User{
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Email:     strings.ToLower(email),
		Roles:     []*Role{},
	}

	for _, opt := range opts {
		opt(&user)
	}

	return user
}

type UserOption func(*User)

func UserWithPassword(password string) UserOption {
	pwdHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return func(user *User) {
		user.Password = string(pwdHash)
	}
}

func UserWithPhoneNumber(phoneNumber int64) UserOption {
	return func(user *User) {
		if phoneNumber != 0 {
			user.PhoneNumber = phoneNumber
		}
	}
}
