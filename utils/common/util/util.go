package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net/http"
	"tzgit.kaixinxiyou.com/utils/common/convert"
	"tzgit.kaixinxiyou.com/utils/common/log"
	"tzgit.kaixinxiyou.com/utils/common/rrand"

	"math"
	"math/rand"
	"net"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}
	return reflect.ValueOf(v).IsNil()
}

// 返回分割后转换为int32
func StringSplitInt(s, sep string) []int32 {
	ret := make([]int32, 0)
	str := strings.Split(s, sep)
	for _, v := range str {
		if v != "" {
			value, err := strconv.Atoi(v)
			if err != nil {
				return ret
			}
			ret = append(ret, int32(value))
		}
	}
	return ret
}
func StringSplitUint(s, sep string) []uint32 {
	ret := make([]uint32, 0)
	str := strings.Split(s, sep)
	for _, v := range str {
		if v != "" {
			value, err := strconv.Atoi(v)
			if err != nil {
				return ret
			}
			ret = append(ret, uint32(value))
		}
	}
	return ret
}
func StringSplitInt64(s, sep string) []int64 {
	ret := make([]int64, 0)
	str := strings.Split(s, sep)
	for _, v := range str {
		if v != "" {
			value := convert.ToInt64(v)
			ret = append(ret, value)
		}
	}
	return ret
}

func TrimBlank(str string) string {
	return strings.Trim(strings.Trim(str, " "), "　")
}

// 指定元素是否在数组中存在
func ArrayExists[T comparable](arrayOrSlice []T, target T) bool {
	for _, v := range arrayOrSlice {
		if v == target {
			return true
		}
	}
	return false
}

func ArrayDeleteUint32(arr []uint32, tar uint32) ([]uint32, bool) {
	del := false
	for i := 0; i <= len(arr)-1; i++ {
		if arr[i] == tar {
			del = true
			arr = append(arr[:i], arr[i+1:]...)
			break
		}
	}
	return arr, del
}

func ArrayDeleteInt32(arr []int32, tar int32) ([]int32, bool) {
	del := false
	for i := 0; i <= len(arr)-1; i++ {
		if arr[i] == tar {
			del = true
			arr = append(arr[:i], arr[i+1:]...)
			break
		}
	}
	return arr, del
}

func ArrayDeleteInt64(arr []int64, tar int64) ([]int64, bool) {
	del := false
	for i := 0; i <= len(arr)-1; i++ {
		if arr[i] == tar {
			del = true
			arr = append(arr[:i], arr[i+1:]...)
			break
		}
	}
	return arr, del
}

func ArrayDeleteString(arr []string, tar string) ([]string, bool) {
	del := false
	for i := 0; i <= len(arr)-1; i++ {
		if arr[i] == tar {
			del = true
			arr = append(arr[:i], arr[i+1:]...)
			break
		}
	}
	return arr, del
}

// 指定元素在数组中存在个数
func ArrayExistsCount(arrayOrSlice interface{}, target interface{}) int32 {
	count := int32(0)
	targetValue := reflect.ValueOf(arrayOrSlice)
	switch reflect.TypeOf(arrayOrSlice).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			i2 := targetValue.Index(i).Interface()
			if i2 == target {
				count++
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(target)).IsValid() {
			return 1
		}
	}

	return count
}

// 指定数组是否有重复项
func ArrayExistsRepeat(arrayOrSlice interface{}) (bool, interface{}) {
	targetValue := reflect.ValueOf(arrayOrSlice)
	switch reflect.TypeOf(arrayOrSlice).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			for j := i + 1; j < targetValue.Len(); j++ {
				if targetValue.Index(i).Interface() == targetValue.Index(j).Interface() {
					return true, targetValue.Index(j).Interface()
				}
			}
		}
	}
	return false, nil
}

// 根据分隔符判断 字符串中是否包含所需字符
func ContainString(fullStr string, splitStr string, target string) bool {
	arrStr := strings.Split(fullStr, splitStr)
	return ArrayExists(arrStr, target)
}

func RemoveElementInt32(arr []int32, tar int32) []int32 {
	i := 0
	for _, val := range arr {
		if val != tar {
			arr[i] = val
			i++
		}
	}
	return arr[:i]
}

func RemoveElementInt32Mul(arr, RemoveArr []int32) []int32 {
	for _, val := range RemoveArr {
		arr = RemoveElementInt32(arr, val)
	}
	return arr
}

// args特指值类型数据，不支持指针或者结构体 //20211103修改 支持切片、数组
func GetKey(args ...interface{}) string {
	str := make([]string, 0)
	keyArr := make([]interface{}, 0)
	for _, key := range args {
		switch reflect.TypeOf(key).Kind() {
		case reflect.Slice, reflect.Array:
			targetValue := reflect.ValueOf(key)
			for i := 0; i < targetValue.Len(); i++ {
				str = append(str, "%v")
				keyArr = append(keyArr, targetValue.Index(i).Interface())
			}
		default:
			str = append(str, "%v")
			keyArr = append(keyArr, key)
		}
	}
	//%v_%v_%v
	result := fmt.Sprintf(strings.Join(str, "_"), keyArr...)
	return result
}

// 用|拼接字符串
func AddString(originStr, addStr interface{}) string {
	res := fmt.Sprintf("%v|%v", originStr, addStr)
	return strings.Trim(res, "|")
}

func GetKeyInt64(v1 int32, v2 int32) int64 {
	return int64(v1)<<32 + int64(v2)
}
func RevertKeyInt64(key int64) (int32, int32) { //对应GetKeyInt64
	return int32(key >> 32), int32(key)
}

// 连接2个int64
func ConnectInt64(a int64, b int64) []byte {
	byteBuf := bytes.NewBuffer([]byte{})
	binary.Write(byteBuf, binary.BigEndian, a)
	binary.Write(byteBuf, binary.BigEndian, b)
	return byteBuf.Bytes()
}
func ConnectInt64s(a []int64) []byte {
	byteBuf := bytes.NewBuffer([]byte{})
	for _, v := range a {
		binary.Write(byteBuf, binary.BigEndian, v)
	}
	return byteBuf.Bytes()
}
func GoID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
func Min(a int32, b int32) int32 {
	if a < b {
		return a
	} else {
		return b
	}
}

func GetAllIP() []net.IP {
	iFaces, err := net.Interfaces()
	ret := make([]net.IP, 0)
	if err != nil {
		return ret
	}

	for _, iFace := range iFaces {
		if iFace.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iFace.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iFace.Addrs()
		if err != nil {
			log.Error("%v", err)
			continue
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			ret = append(ret, ip)
		}
	}
	return ret
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}
	return ip
}

// 判断本机ip 是否是有C类ip
func GetHaveIpC() bool {
	allIp := GetAllIP()
	ip1 := net.ParseIP("192.168.0.0").To4()
	ip2 := net.ParseIP("192.168.255.255").To4()
	for _, v := range allIp {
		if bytes.Compare(v, ip1) >= 0 && bytes.Compare(v, ip2) <= 0 {
			return true
		}
	}
	return false
}

// 判断本机ip 是否是有B类ip
func GetHaveIpB() bool {
	allIp := GetAllIP()
	ip1 := net.ParseIP("172.16.0.0").To4()
	ip2 := net.ParseIP("172.31.255.255").To4()
	for _, v := range allIp {
		if bytes.Compare(v, ip1) >= 0 && bytes.Compare(v, ip2) <= 0 {
			return true
		}
	}
	return false
}

// 判断本机ip 是否是有A类ip
func GetHaveIpA() bool {
	allIp := GetAllIP()
	ip1 := net.ParseIP("10.0.0.0").To4()
	ip2 := net.ParseIP("10.255.255.255").To4()
	for _, v := range allIp {
		if bytes.Compare(v, ip1) >= 0 && bytes.Compare(v, ip2) <= 0 {
			return true
		}
	}
	return false
}

func GetRandomArrayInt32(a []int32) int32 {
	l := int32(len(a))
	if a == nil || l == 0 {
		return 0
	}
	i := rrand.Default.RandInterval(0, l-1)
	return a[i]
}

func GetRandomArrayInt64(a []int64) int64 {
	l := int32(len(a))
	if a == nil || l == 0 {
		return 0
	}
	i := rrand.Default.RandInterval(0, l-1)
	return a[i]
}

func GetRandomArrayString(a []string) string {
	l := int32(len(a))
	if a == nil || l == 0 {
		return ""
	}
	i := rrand.Default.RandInterval(0, l-1)
	return a[i]
}

func ZhengTaiFloat64(x float64, miu float64, sigma float64) float64 {
	randomNormal := 1 / (math.Sqrt(2*math.Pi) * sigma) * math.Pow(math.E, -math.Pow(x-miu, 2)/(2*math.Pow(sigma, 2)))
	//注意下是x-miu，我看网上好多写的是miu-miu是不对的
	return randomNormal
}

func ZhengTaiRandomInt64(r *rand.Rand, miu float64, sigma float64) int64 {
	var x int64
	var y, dScope float64
	for {
		x = r.Int63n(int64(miu)*2) + 1
		y = ZhengTaiFloat64(float64(x), miu, sigma) * 100000
		dScope = float64(r.Int63n(int64(ZhengTaiFloat64(miu, miu, sigma) * 100000)))
		//注意下传的是两个miu
		if dScope <= y {
			break
		}
	}
	return x
}

func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}
