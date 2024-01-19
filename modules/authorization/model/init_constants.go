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

var OwnerAuthority = NewAuthority("OWNER", "Право владельца системы")

var ReadAssetAuthority = NewAuthority("READ_ASSET", "Чтение активов")
var CreateAssetAuthority = NewAuthority("CREATE_ASSET", "Создание активов")
var UpdateAssetAuthority = NewAuthority("UPDATE_ASSET", "Редактирование активов")
var DeleteAssetAuthority = NewAuthority("DELETE_ASSET", "Удаление активов")

var ReadCashFlowAuthority = NewAuthority("READ_CASH_FLOW", "Чтение приходов/расходов")
var CreateCashFlowAuthority = NewAuthority("CREATE_CASH_FLOW", "Создание приходов/расходов")
var UpdateCashFlowAuthority = NewAuthority("UPDATE_CASH_FLOW", "Редактирование приходов/расходов")
var DeleteCashFlowAuthority = NewAuthority("DELETE_CASH_FLOW", "Удаление приходов/расходов")

var CreateAttachmentAuthority = NewAuthority("CREATE_ATTACHMENT", "Создание вложений")
var DeleteAttachmentAuthority = NewAuthority("DELETE_ATTACHMENT", "Удаление вложений")

var ReadAuthorityAuthority = NewAuthority("READ_AUTHORITY", "Чтение прав")
var CreateAuthorityAuthority = NewAuthority("CREATE_AUTHORITY", "Создание прав")
var UpdateAuthorityAuthority = NewAuthority("UPDATE_AUTHORITY", "Редактирование прав")
var DeleteAuthorityAuthority = NewAuthority("DELETE_AUTHORITY", "Удаление прав")

var ReadRoleAuthority = NewAuthority("READ_ROLE", "Чтение ролей")
var CreateRoleAuthority = NewAuthority("CREATE_ROLE", "Создание ролей")
var UpdateRoleAuthority = NewAuthority("UPDATE_ROLE", "Редактирование ролей")
var DeleteRoleAuthority = NewAuthority("DELETE_ROLE", "Удаление ролей")

var ReadUserAuthority = NewAuthority("READ_USER", "Чтение пользователей")
var CreateUserAuthority = NewAuthority("CREATE_USER", "Создание пользователей")
var UpdateUserAuthority = NewAuthority("UPDATE_USER", "Редактирование пользователей")
var DeleteUserAuthority = NewAuthority("DELETE_USER", "Удаление пользователей")
var EditRoleUserAuthority = NewAuthority("EDIT_ROLE_USER", "Редактирование ролей пользователей")
var EditAuthorityUserAuthority = NewAuthority("EDIT_AUTHORITY_USER", "Редактирование прав пользователей")

var ReadCityAuthority = NewAuthority("READ_CITY", "Чтение городов")
var CreateCityAuthority = NewAuthority("CREATE_CITY", "Создание городов")
var UpdateCityAuthority = NewAuthority("UPDATE_CITY", "Редактирование городов")
var DeleteCityAuthority = NewAuthority("DELETE_CITY", "Удаление городов")

var ReadCountryAuthority = NewAuthority("READ_COUNTRY", "Чтение стран")
var CreateCountryAuthority = NewAuthority("CREATE_COUNTRY", "Создание стран")
var UpdateCountryAuthority = NewAuthority("UPDATE_COUNTRY", "Редактирование стран")
var DeleteCountryAuthority = NewAuthority("DELETE_COUNTRY", "Удаление стран")

var ReadCurrencyAuthority = NewAuthority("READ_CURRENCY", "Чтение валют")
var CreateCurrencyAuthority = NewAuthority("CREATE_CURRENCY", "Создание валют")
var UpdateCurrencyAuthority = NewAuthority("UPDATE_CURRENCY", "Редактирование валют")
var DeleteCurrencyAuthority = NewAuthority("DELETE_CURRENCY", "Удаление валют")

var AdminRole = NewRole("ADMIN", "Роль админа", RoleWithAuthorities(
	&ReadAssetAuthority, &CreateAssetAuthority, &UpdateAssetAuthority, &DeleteAssetAuthority,
	&ReadCashFlowAuthority, &CreateCashFlowAuthority, &UpdateCashFlowAuthority, &DeleteCashFlowAuthority,
	&CreateAttachmentAuthority, &DeleteAttachmentAuthority,
	&ReadAuthorityAuthority, &CreateAuthorityAuthority, &UpdateAuthorityAuthority, &DeleteAuthorityAuthority,
	&ReadRoleAuthority, &CreateRoleAuthority, &UpdateRoleAuthority, &DeleteRoleAuthority,
	&ReadUserAuthority, &CreateUserAuthority, &UpdateUserAuthority, &DeleteUserAuthority, &EditRoleUserAuthority, &EditAuthorityUserAuthority,
	&ReadCityAuthority, &CreateCityAuthority, &UpdateCityAuthority, &DeleteCityAuthority,
	&ReadCountryAuthority, &CreateCountryAuthority, &UpdateCountryAuthority, &DeleteCountryAuthority,
	&ReadCurrencyAuthority, &CreateCurrencyAuthority, &UpdateCurrencyAuthority, &DeleteCurrencyAuthority,
))

var InternalUserRole = NewRole("USER", "Роль пользователя", RoleWithAuthorities(
	&ReadUserAuthority,
))
