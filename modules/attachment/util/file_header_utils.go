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
	commonUtil "assets/common/util"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strings"
)

func DetectContentType(fileName string) string {
	contentType := "application/octet-stream"
	_, extension := commonUtil.SubstringLast(fileName, ".")
	if value, ok := extensionAndContentTypeMap[strings.ToLower(extension)]; ok {
		contentType = value
	}
	return contentType
}

var extensionAndContentTypeMap = map[string]string{
	// Изображения
	"bmp":   "image/bmp",
	"gif":   "image/gif",
	"jpeg":  "image/jpeg",
	"jpg":   "image/jpeg",
	"pjpeg": "image/pjpeg",
	"png":   "image/png",
	"svg":   "image/svg+xml",
	"tiff":  "image/tiff",
	"webp":  "image/webp",
	"ico":   "image/vnd.microsoft.icon",

	// Аудио
	"m4a":  "audio/voice",
	"adp":  "audio/adpcm",
	"au":   "audio/basic",
	"mid":  "audio/midi",
	"mp4a": "audio/mp4",
	"mpga": "audio/mpeg",
	"oga":  "audio/ogg",
	"aac":  "audio/aac",
	"mp3":  "audio/mpeg",
	"opus": "audio/opus",
	"wav":  "audio/wav",
	"weba": "audio/webm",

	// Видео
	"3gp":  "video/3gpp",
	"3g2":  "video/3gpp2",
	"h261": "video/h261",
	"h263": "video/h263",
	"h264": "video/h264",
	"jpgv": "video/jpeg",
	"jpm":  "video/jpm",
	"mj2":  "video/mj2",
	"mp4":  "video/mp4",
	"mpg":  "video/mpeg",
	"mpeg": "video/mpeg",
	"ogv":  "video/ogg",
	"qt":   "video/quicktime",
	"mov":  "video/quicktime",
	"ts":   "video/mp2t",
	"webm": "video/webm",
	"avi":  "video/x-msvideo",

	// Контакты
	"vcf":   "text/vcard",
	"vcard": "text/vcard",

	// Документы
	"apk":    "application/vnd.android.package-archive",
	"epub":   "application/epub+zip",
	"gz":     "application/gzip",
	"jar":    "application/java-archive",
	"json":   "application/json",
	"jsonld": "application/ld+json",
	"doc":    "application/msword",
	"bin":    "application/octet-stream",
	"ogx":    "application/ogg",
	"pdf":    "application/pdf",
	"php":    "application/php",
	"rtf":    "application/rtf",
	"azw":    "application/vnd.amazon.ebook",
	"mpkg":   "application/vnd.apple.installer+xml",
	"xul":    "application/vnd.mozilla.xul+xml",
	"xls":    "application/vnd.ms-excel",
	"eot":    "application/vnd.ms-fontobject",
	"ppt":    "application/vnd.ms-powerpoint",
	"odp":    "application/vnd.oasis.opendocument.presentation",
	"ods":    "application/vnd.oasis.opendocument.spreadsheet",
	"odt":    "application/vnd.oasis.opendocument.text",
	"pptx":   "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"xlsx":   "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"docx":   "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"rar":    "application/vnd.rar",
	"vsd":    "application/vnd.visio",
	"7z":     "application/x-7z-compressed",
	"abw":    "application/x-abiword",
	"bz":     "application/x-bzip",
	"bz2":    "application/x-bzip2",
	"csh":    "application/x-csh",
	"arc":    "application/x-freearc",
	"sh":     "application/x-sh",
	"swf":    "application/x-shockwave-flash",
	"tar":    "application/x-tar",
	"xhtml":  "application/xhtml+xml",
	"xml":    "application/xml",
	"zip":    "application/zip",
	"otf":    "font/otf",
	"ttf":    "font/ttf",
	"woff":   "font/woff",
	"woff2":  "font/woff2",
	"ics":    "text/calendar",
	"css":    "text/css",
	"csv":    "text/csv",
	"html":   "text/html",
	"js":     "text/javascript",
	"txt":    "text/plain",
}
