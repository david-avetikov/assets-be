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

import "strings"

var NilPageable = Pageable{}

type Pageable struct {
	Size int   `json:"size"`
	Page int   `json:"page"`
	Sort *Sort `json:"sort,omitempty"`
}

func NewPageable(size int, page int) Pageable {
	if size == 0 {
		size = 20
	}
	return Pageable{Size: size, Page: page}
}

func (pageable Pageable) WithSortString(sortStr string) Pageable {
	if sortStr != "" {
		if fieldAndOrder := strings.Split(sortStr, ","); len(fieldAndOrder) == 2 {
			field, orderStr := fieldAndOrder[0], fieldAndOrder[1]

			var sortOrder SortOrder
			switch orderStr {
			case "desc":
				sortOrder = Desc
			default:
				sortOrder = Asc
			}

			pageable.Sort = NewSort(field, sortOrder)
		} else {
			field := fieldAndOrder[0]
			pageable.Sort = NewSort(field, Asc)
		}
	}
	return pageable
}

func (pageable Pageable) WithSort(sort Sort) Pageable {
	pageable.Sort = &sort
	return pageable
}
