package util

import (
	"regexp"
	"strings"
)

const (
	ChineseWordRegx         = "/[\\u4e00-\\u9fa5]/"
	ChineseAppendixWordRegx = "/[\\u9FA6-\\u9fcb]/"
	ChineseTokenRegx        = "/[^\\x00-\\xff]/"
	NumberWordRegx          = "/[0-9]/"
)

var (
	ChineseWordPattern         *regexp.Regexp
	ChineseAppendixWordPattern *regexp.Regexp
	ChineseTokenPattern        *regexp.Regexp
	NumberWordPattern          *regexp.Regexp
)

func init() {
	ChineseWordPattern, _ = regexp.Compile(ChineseWordRegx)
	ChineseAppendixWordPattern, _ = regexp.Compile(ChineseAppendixWordRegx)
	ChineseTokenPattern, _ = regexp.Compile(ChineseTokenRegx)
	NumberWordPattern, _ = regexp.Compile(NumberWordRegx)
}

// WordCount 统计给定的字符串字符数据
func WordCount(words string) int64 {
	iTotal := 0
	inum := 0
	eTotal := 0
	sTotal := 0

	for _, word := range strings.Split("", words) {
		// 基本汉字
		if ChineseWordPattern.MatchString(word) {
			iTotal++
		}
		// 基本汉字补充
		if ChineseAppendixWordPattern.MatchString(word) {
			iTotal++
		}
		// 中文标点加中文字
		if ChineseTokenPattern.MatchString(word) {
			sTotal++
		} else {
			// 英文
			eTotal++
		}
		// 数字
		if NumberWordPattern.MatchString(word) {
			inum++
		}
	}
	return int64(iTotal + eTotal + sTotal + inum)
}
