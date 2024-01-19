package custom_error

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
	"assets/common/iface"
	"net/http"
)

type baseHttpError struct {
	err  string
	code int
}

func NewHttpError(err string, code int) iface.HttpError { return &baseHttpError{err, code} }
func (httpErr baseHttpError) Error() string             { return httpErr.err }
func (httpErr baseHttpError) Code() int                 { return httpErr.code }

var (
	NeedAuthorizationHeaderError          = NewHttpError("need authorization header", http.StatusUnauthorized)
	UnknownGrantTypeError                 = NewHttpError("unknown grant_type", http.StatusBadRequest)
	UnsupportedTokenTypeError             = NewHttpError("unsupported token type", http.StatusBadRequest)
	NotFoundError                         = NewHttpError("not found", http.StatusNotFound)
	TokenExpiredError                     = NewHttpError("token is expired", http.StatusUnauthorized)
	TokenInvalidError                     = NewHttpError("for this operation need authorization header with valid bearer token", http.StatusUnauthorized)
	UsernameAndPasswordMastNotBeNullError = NewHttpError("username and password mast not be null", http.StatusBadRequest)
	UsernameOrPasswordIsIncorrectError    = NewHttpError("username or password is incorrect", http.StatusBadRequest)
	UserBlockedError                      = NewHttpError("user is blocked", http.StatusForbidden)
	NotEnoughRightsError                  = NewHttpError("not enough rights", http.StatusForbidden)
	ParseZeroValueError                   = NewHttpError("parse zero value", http.StatusBadRequest)
	IllegalArgumentError                  = NewHttpError("illegal argument", http.StatusBadRequest)
	CacheNotFound                         = NewHttpError("cache not found", http.StatusOK)
	AlreadyExists                         = NewHttpError("already exists", http.StatusConflict)
	CouldNotConnectError                  = NewHttpError("couldn't connect", http.StatusServiceUnavailable)
	NotImplementedError                   = NewHttpError("not implemented", http.StatusNotImplemented)
	LicenseNotActivated                   = NewHttpError("license not activated", http.StatusConflict)
	LicenseAlreadyActivated               = NewHttpError("license already activated", http.StatusConflict)
)
