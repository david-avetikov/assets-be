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
	"strings"
	"unicode"
)

func SubstringLast(value string, str string) (string, string) {
	pos := strings.LastIndex(value, str)
	if pos == -1 {
		return value, ""
	}
	adjustedPos := pos + len(str)
	if adjustedPos >= len(value) {
		return value[:pos], ""
	}
	return value[:pos], value[adjustedPos:]
}

func GetLiterals(str string) []string {
	words := strings.Split(str, " ")
	if len(words) == 1 {
		return getLiterals(words[0])
	}

	result := make([]string, 0)
	for _, word := range words {
		literals := getLiterals(word)
		result = append(result, literals...)
	}
	return result
}

func getLiterals(str string) []string {
	str = strings.TrimSpace(str)
	chars := strings.Split(str, "")
	if len(chars) <= 3 {
		return []string{strings.Join(chars, "")}
	}

	literal := chars[:3]
	result := make([]string, 0)
	result = append(result, strings.Join(literal, ""))

	for i := 3; i < len(chars); i++ {
		literal = append(literal, chars[i])
		result = append(result, strings.Join(literal, ""))
	}
	return result
}

func HasUpperCase(str string) bool {
	for _, char := range str {
		if unicode.IsUpper(char) {
			return true
		}
	}
	return false
}
