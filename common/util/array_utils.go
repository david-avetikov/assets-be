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
	"errors"
	"sort"
)

func Map[T any, Z any](array []T, mapFunc func(T) Z) []Z {
	var result []Z
	for _, item := range array {
		result = append(result, mapFunc(item))
	}
	return result
}

func GetMapKeys[K comparable, V any](mp map[K]V) []K {
	var keys []K
	for key := range mp {
		keys = append(keys, key)
	}
	return keys
}

func GetMapValues[K comparable, V any](mp map[K]V) []V {
	var values []V
	for _, value := range mp {
		values = append(values, value)
	}
	return values
}

func ArrayContains[T comparable](array []T, elem T) bool {
	for _, item := range array {
		if item == elem {
			return true
		}
	}
	return false
}

func ArrayFilter[T any](array []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range array {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func ArrayIndexOf[T any](array []T, predicate func(T) bool) (int, error) {
	for i, elem := range array {
		if predicate(elem) {
			return i, nil
		}
	}
	return -1, errors.New("NoSuchElement")
}

func ArrayIndexesOf[T any](array []T, predicate func(T) bool) []int {
	var result []int
	for i, elem := range array {
		if predicate(elem) {
			result = append(result, i)
		}
	}
	return result
}

func ArrayExcludeByIndex[T any](array []T, indexes ...int) []T {
	var result []T
	sort.Ints(indexes)
	for index, value := range array {
		if !ArrayContains(indexes, index) {
			result = append(result, value)
		}
	}
	return result
}

func ArrayFindFirst[T any](array []T, predicate func(T) bool) (bool, T) {
	var result T
	for _, item := range array {
		if predicate(item) {
			return true, item
		}
	}
	return false, result
}

func Unique[T comparable](array []T) []T {
	result := make([]T, 0, len(array))
	uniqueMap := make(map[T]struct{})

	var member struct{}
	for _, value := range array {
		if _, ok := uniqueMap[value]; !ok {
			uniqueMap[value] = member
			result = append(result, value)
		}
	}

	return result
}

func ArrayNotZero[T comparable](array []T) []T {
	var zero T
	var result []T
	for _, item := range array {
		if zero != item {
			result = append(result, item)
		}
	}
	return result
}

func ArrayToBatch[T any](array []T, batchSize int) [][]T {
	var result [][]T
	if array == nil || len(array) == 0 || batchSize == 0 {
		return result
	}

	iteration := len(array) / batchSize
	if len(array)%batchSize != 0 {
		iteration++
	}
	for i := 0; i < iteration; i++ {
		startIndex := i * batchSize
		endIndex := (i + 1) * batchSize
		if endIndex > len(array) {
			endIndex = len(array)
		}
		result = append(result, array[startIndex:endIndex])
	}

	return result
}

func ArrayIntersection[T any](from []T, what []T, compareFunc func(item1 T, item2 T) bool) []T {
	var result []T
	if from == nil || len(from) == 0 {
		return result
	}
	if what == nil || len(what) == 0 {
		return from
	}

	for _, item := range from {
		if !ArrayContainsItem(what, item, compareFunc) {
			result = append(result, item)
		}
	}
	return result
}

func ArrayContainsItem[T any](array []T, item T, compareFunc func(item1 T, item2 T) bool) bool {
	for _, aItem := range array {
		if compareFunc(aItem, item) {
			return true
		}
	}
	return false
}

func EqualKeyAndValue(pairs map[any]any) bool {
	for key, value := range pairs {
		if key != value {
			return false
		}
	}
	return true
}
