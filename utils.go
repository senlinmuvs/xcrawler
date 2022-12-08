package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func CurMills() int64 {
	t := time.Now()
	return int64(time.Nanosecond) * t.UnixNano() / int64(time.Millisecond)
}

func CurTimeStr() string {
	t := time.Now()
	ts := t.Format("2006-01-02 15:04:05")
	return ts
}
func Split(s, sep string) []string {
	arr := []string{}
	arr_ := strings.Split(s, sep)
	for _, a := range arr_ {
		a = strings.Trim(a, " ")
		if a != "" {
			arr = append(arr, a)
		}
	}
	return arr
}
func ByteArrEq(arr1, arr2 []byte) bool {
	if len(arr1) != len(arr2) {
		return false
	}
	for i, e := range arr1 {
		if e != arr2[i] {
			return false
		}
	}
	return true
}
func ArrEqIg(arr1, arr2 []string, igEle ...string) bool {
	for i, ar := range arr1 {
		isIg := false
		for _, ie := range igEle {
			if ar == ie {
				isIg = true
				break
			}
		}
		if isIg {
			continue
		}
		if ar != arr2[i] {
			return false
		}
	}
	return true
}
func ToInt(s string) int {
	if s == "" {
		return 0
	}
	x, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Println("err: toInt() err", err)
		return 0
	}
	return int(x)
}
func ReadFile(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}
	return string(content)
}
func ReadLines(file string) []string {
	cont := ReadFile(file)
	lines := strings.Split(cont, "\n")
	return lines
}
func SubUnicode(s string, i, j int) string {
	if i < 0 || j < 0 {
		return ""
	}
	srune := []rune(s)
	return string(srune[i:j])
}
func Compact(s string, cs ...string) string {
	for _, c := range cs {
		for {
			cc := c + c
			i := strings.Index(s, cc)
			if i < 0 {
				break
			}
			s = strings.ReplaceAll(s, cc, c)
		}
	}
	return s
}
func Rand(n int) int {
	cur := time.Now().UnixNano()
	rand.Seed(cur)
	return rand.Intn(n)
}
