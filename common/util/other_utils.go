package util

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
	log "github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"io"
	"os"
	"reflect"
)

var ApplicationName = "Please set application name"

func InitApplication(applicationName string) {
	setLocalProfileIfNeeded()
	setApplicationName(applicationName)
	InitLogger()
	PrintAllEnv()
}

func setLocalProfileIfNeeded() {
	if os.Getenv("PROFILE") == "" && os.Getenv("KUBERNETES_NAMESPACE") != "" {
		if err := os.Setenv("PROFILE", "kubernetes"); err != nil {
			log.WithError(err).Error("PROFILE is not set, use default kubernetes")
		}
	} else if os.Getenv("PROFILE") == "" {
		if err := os.Setenv("PROFILE", "local"); err != nil {
			log.WithError(err).Error("PROFILE is not set, use default local")
		}
	}
}

func setApplicationName(applicationName string) {
	ApplicationName = applicationName
}

func CloseChecker(fn io.Closer) {
	if err := fn.Close(); err != nil {
		log.WithError(err).Error("Couldn't invoke close method")
	}
}

func CheckedError(err error) {
	closer.Checked(func() error { return err }, true)
}

func CheckedRun(runnable iface.Runnable) {
	closer.Checked(func() error {
		runnable.Run()
		return nil
	}, true)
}

func IsZeroObject(object any) bool {
	return reflect.ValueOf(object).IsZero()
}
