package service

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
	"assets/common/config"
	commonModel "assets/common/model"
	commonUtil "assets/common/util"
	"assets/modules/attachment/model"
	"assets/modules/attachment/repository"
	"assets/modules/attachment/util"
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
)

var attachmentSrv AttachmentService

type AttachmentService interface {
	GetById(ctx context.Context, id uuid.UUID) (date []byte, contentType string, fileName string, filePath string)
	GetMetaById(ctx context.Context, id uuid.UUID, fields string) model.Attachment
	GetAllMeta(ctx context.Context, page commonModel.Pageable, filter model.AttachmentFilter, fields string) commonModel.Page[model.Attachment]
	Create(ctx context.Context, fileHeader *multipart.FileHeader) model.Attachment
	DeleteById(ctx context.Context, id uuid.UUID)
}

type attachmentService struct {
	repository        repository.AttachmentRepository
	s3UploadListeners []func(context context.Context, attachment model.Attachment)
}

func GetAttachmentService() AttachmentService {
	if attachmentSrv != nil {
		return attachmentSrv
	}

	attachmentSrv = &attachmentService{repository: repository.GetAttachmentRepository()}
	return attachmentSrv
}

func (service *attachmentService) GetById(ctx context.Context, id uuid.UUID) (date []byte, contentType string, fileName string, filePath string) {
	attachment := service.GetMetaById(ctx, id, "")
	file := commonUtil.MustOne(os.Open(attachment.Path))
	defer commonUtil.CloseChecker(file)

	return commonUtil.MustOne(io.ReadAll(file)), attachment.MimeType, attachment.FileName, attachment.Path
}

func (service *attachmentService) GetMetaById(ctx context.Context, id uuid.UUID, fields string) model.Attachment {
	return service.repository.GetById(ctx, []uuid.UUID{id})[0]
}

func (service *attachmentService) GetAllMeta(ctx context.Context, page commonModel.Pageable, filter model.AttachmentFilter, fields string) commonModel.Page[model.Attachment] {
	return service.repository.GetAllWithPage(ctx, page)
}

func (service *attachmentService) Create(ctx context.Context, fileHeader *multipart.FileHeader) model.Attachment {
	file := commonUtil.MustOne(fileHeader.Open())
	defer commonUtil.CloseChecker(file)
	contentType := util.DetectContentType(fileHeader.Filename)
	attachment := model.NewAttachment(fileHeader.Filename, contentType, fileHeader.Size)

	ginCtx := commonUtil.ConvertContext(ctx)
	dir := fmt.Sprintf("%s/%s", config.CoreConfig.Attachment.UploadPath, attachment.CreateDate.Format("2006-01-02"))
	commonUtil.Must(os.MkdirAll(dir, os.ModePerm))
	attachment.Path = fmt.Sprintf("%s/%s", dir, attachment.FileName)
	commonUtil.Must(ginCtx.SaveUploadedFile(fileHeader, attachment.Path))

	return service.repository.Create(ctx, []model.Attachment{attachment})[0]
}

func (service *attachmentService) DeleteById(ctx context.Context, id uuid.UUID) {
	service.repository.DeleteById(ctx, id)
}
