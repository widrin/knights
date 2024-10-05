package validate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	nickRegex     = "^[\\p{Han}\\w]+$"
	pswRegex      = "^[\\w]{6,18}$"
	repNumRegex   = ".*[0-9]{6}.*"
	realnameRegex = "^[\\p{Han}]{2,10}$"
	onlyHanRegex  = "^[\\p{Han}]+$"
	onlyNumRegex  = "^[0-9]*$"
)

type Validate struct {
	wordPool []string
}

//敏感词库1
func (v *Validate) SetWordPool(array []string) {
	v.wordPool = array
}

func (v *Validate) CheckOnlyNum(str string) bool {
	reg := regexp.MustCompile(onlyNumRegex)
	return reg.MatchString(str)
}

func (v *Validate) CheckFamilyName(name string) bool {
	reg := regexp.MustCompile(nickRegex)
	return reg.MatchString(name)
}

//昵称长度不符合，请重新输入
//昵称含有非法字符，请重新输入
func (v *Validate) CheckNick(nick string) string {
	reg := regexp.MustCompile(nickRegex)
	if reg.MatchString(nick) == false {
		return "昵称由6个汉字或12个字母数字组成"
	}
	uLen := len(nick)
	rLen := len([]rune(nick))
	zhNum := (rLen - uLen) / 2
	if (uLen - zhNum + zhNum*2) > 12 {
		return "昵称长度不符合，请重新输入"
	}
	//检测敏感词
	res := v.CheckIllegalWord(nick)
	if res == false {
		return "昵称含有非法字符，请重新输入"
	}
	return ""
}

//检查工会名
func (v *Validate) CheckGuildName(gName string) string {
	reg := regexp.MustCompile(nickRegex)
	if reg.MatchString(gName) == false {
		return "昵称由6个汉字或12个字母数字组成"
	}
	uLen := len(gName)
	rLen := len([]rune(gName))
	zhNum := (rLen - uLen) / 2
	if (uLen - zhNum + zhNum*2) > 8 {
		return "长度不符合，请重新输入"
	}
	//检测敏感词
	res := v.CheckIllegalWord(gName)
	if res == false {
		return "昵称含有非法字符，请重新输入"
	}
	return ""
}

//检查只允许输入汉字
func (v *Validate) CheckOnlyHan(han string, hLen int32) string {
	reg := regexp.MustCompile(onlyHanRegex)
	if reg.MatchString(han) == false {
		return "只允许输入汉字"
	}
	uLen := len(han)
	rLen := len([]rune(han))
	zhNum := (rLen - uLen) / 2
	if int32(uLen-zhNum+zhNum*2) > hLen {
		fmt.Printf("长度不符合，%v，%v，%v", uLen, rLen, zhNum)
		return "长度不符合，请重新输入"
	}
	//检测敏感词
	res := v.CheckIllegalWord(han)
	if res == false {
		fmt.Printf("昵称含有非法字符，%v", han)
		return "昵称含有非法字符，请重新输入"
	}
	return ""
}

//只验证长度
func (v *Validate) CheckLen(han string, hLen int32) string {
	uLen := len(han)
	rLen := len([]rune(han))
	zhNum := (rLen - uLen) / 2
	if int32(uLen-zhNum+zhNum*2) > hLen {
		fmt.Printf("长度不符合，%v，%v，%v", uLen, rLen, zhNum)
		return "长度不符合，请重新输入"
	}
	//检测敏感词
	res := v.CheckIllegalWord(han)
	if res == false {
		fmt.Printf("昵称含有非法字符，%v", han)
		return "昵称含有非法字符，请重新输入"
	}
	return ""
}

func (v *Validate) CheckPswEmptyStr(psw string) bool {
	n := 0
	arr := strings.Split(psw, "")
	num := len(arr)
	for i := 0; i < num; i++ {
		if (i + 1) < num {
			if strings.EqualFold(arr[i], arr[(i+1)]) {
				n++
				if n >= 2 {
					return false
				}
			} else {
				n = 0
			}
		}
	}
	return true
}
func (v *Validate) CheckPswEmptyNum(psw string) bool {
	arr := strings.Split(psw, "")
	num := len(arr)
	sArr := make(map[int]int)
	mArr := make(map[int]int)
	for i := 0; i < num; i++ {
		val, _ := strconv.Atoi(arr[i])
		sArr[i] = val - i
		mArr[i] = val + i
	}
	var res1 bool
	var res2 bool
	for i := 0; i < num; i++ {
		if (i + 1) < num {
			if sArr[i] != sArr[i+1] {
				res1 = true
				break
			}
		}
	}
	if !res1 {
		return false
	}
	for i := 0; i < num; i++ {
		if (i + 1) < num {
			if mArr[i] != mArr[i+1] {
				res2 = true
				break
			}
		}
	}
	if !res2 {
		return false
	}
	return true
}

//替换敏感词
func (v *Validate) ReplaceIllegalWord(str string) string {
	for _, v := range v.wordPool {
		if strings.Contains(str, v) == true {
			str = strings.Replace(str, v, "**", -1)
		}
	}
	return str
}

//查找敏感词
func (v *Validate) CheckIllegalWord(str string) bool {
	if str == "" {
		return false
	}
	for _, w := range v.wordPool {
		if strings.Contains(str, w) == true {
			return false
		}
	}
	return true
}
