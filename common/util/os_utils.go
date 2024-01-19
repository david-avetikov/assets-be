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
	"archive/tar"
	"compress/gzip"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
	"time"
)

func PrintAllEnv() {
	log.Debugf("Environments:\n%v", strings.Join(os.Environ(), "\n"))
}

func MkDir(filePath string) error {
	dirs := strings.Split(filePath, "/")
	directoryPath := strings.Join(dirs[:len(dirs)-1], "/")
	return os.MkdirAll(directoryPath, os.ModePerm)
}

func ArchiveFile(filePath string) error {
	timeMark := time.Now().UTC().Format(time.RFC3339)
	directories, fileFullName := SubstringLast(filePath, "/")
	fileName, _ := SubstringLast(fileFullName, ".")

	archiveFileName := directories + string(os.PathSeparator) + fileName + "." + timeMark + ".tar.gz"
	targetFile, err := os.Create(archiveFileName)
	if err != nil {
		return err
	}
	defer CloseChecker(targetFile)

	gzipWriter := gzip.NewWriter(targetFile)
	defer CloseChecker(gzipWriter)
	tarWriter := tar.NewWriter(gzipWriter)
	defer CloseChecker(tarWriter)

	sourceFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer CloseChecker(sourceFile)

	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(sourceFileInfo, sourceFileInfo.Name())
	if err != nil {
		return err
	}
	header.Name = fileFullName
	if err = tarWriter.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tarWriter, sourceFile)
	if err != nil {
		return err
	}

	if err = os.Remove(filePath); err != nil {
		return err
	}
	return nil
}
