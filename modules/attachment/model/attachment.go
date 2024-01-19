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
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Attachment struct {
	ID         uuid.UUID  `json:"id,omitempty" gorm:"type:uuid;primaryKey"`
	CreateDate *time.Time `json:"createDate,omitempty"`
	FileName   string     `json:"fileName,omitempty"`
	MimeType   string     `json:"mimeType,omitempty"`
	Size       int64      `json:"size,omitempty"`
	Path       string     `json:"path,omitempty"`
}

func (attachment Attachment) GetID() uuid.UUID {
	return attachment.ID
}

func (attachment *Attachment) BeforeCreate(tx *gorm.DB) error {
	if commonUtil.IsZeroObject(attachment.ID) {
		attachment.ID = uuid.New()
	}
	return nil
}

type AttachmentFilter struct {
	IDs      []string `form:"ids"`
	FileName string   `form:"fileName"`
}

func NewAttachment(fileName string, mimeType string, size int64) Attachment {
	now := time.Now()
	return Attachment{
		CreateDate: &now,
		FileName:   fileName,
		MimeType:   mimeType,
		Size:       size,
	}
}
