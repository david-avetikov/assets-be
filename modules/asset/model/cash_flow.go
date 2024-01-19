package model

import (
	commonUtil "assets/common/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
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

type CashFlow struct {
	ID          uuid.UUID  `json:"id,omitempty" gorm:"type:uuid;primaryKey"`
	Type        string     `json:"type,omitempty"`
	Date        *time.Time `json:"date,omitempty"`
	Amount      float64    `json:"amount,omitempty"`
	Currency    string     `json:"currency,omitempty"`
	Description string     `json:"description,omitempty"`
}

func (cashFlow CashFlow) GetID() uuid.UUID {
	return cashFlow.ID
}

func (cashFlow *CashFlow) BeforeCreate(tx *gorm.DB) error {
	if commonUtil.IsZeroObject(cashFlow.ID) {
		cashFlow.ID = uuid.New()
	}
	return nil
}
