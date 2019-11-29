package strutil

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/mozillazg/go-pinyin"
	"github.com/pborman/uuid"
	"github.com/PKUJohnson/solar/std/common"
	"github.com/PKUJohnson/solar/std/toolkit/intutil"
	"github.com/PKUJohnson/solar/std/toolkit/sliceutil"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"qiniupkg.com/x/url.v7"
)

const (
	UriRegexRule         = `(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`
	DisplayNameRegexRule = `^[-\w|\p{Han}]+$`
)

var (
	uriRep = regexp.MustCompile(UriRegexRule)
)

func GenContent(content string, length int) string {
	result := []rune(content)
	if len(result) > length {
		result = result[:length]
	}
	return strings.TrimSpace(string(result))
}

var charArr = []string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
	"k", "l", "m", "n", "o", "p", "q", "r", "s", "t",
	"u", "v", "w", "x", "y", "z", "A", "B", "C", "D",
	"E", "F", "G", "H", "I", "J", "K", "L", "M", "N",
	"O", "P", "Q", "R", "S", "T", "U", "V", "W", "X",
	"Y", "Z"}

var base64Images = []string{
	"data:image/jpeg;base64",
	"data:image/png;base64",
	"data:image/gif;base64",
	"data:text/plain;base64",
}

func FromInt64(i int64) string {
	return strconv.FormatInt(i, 10)
}

func FromUInt64(i uint64) string {
	return strconv.FormatUint(i, 10)
}

func FromFloat64(i float64) string {
	return strconv.FormatFloat(i, 'f', 2, 64)
}

func ToInt64(str string) int64 {
	if b, err := strconv.ParseInt(str, 10, 64); err != nil {
		return 0
	} else {
		return b
	}
}

func ToInt32(str string) int32 {
	/*
		int means at least 32 ,not an alias for int32
	*/
	if b, err := strconv.ParseInt(str, 10, 32); err == nil {
		return int32(b)
	}
	return 0
}

/*
 * 字符串中ID以,分割的ids
 */
func ToInt64Arr(str string) []int64 {

	res := make([]int64, 0)
	idSplit := strings.Split(str, common.CommaSeparator)

	for _, val := range idSplit {
		valint64 := intutil.ToInt64(val, 0)
		if valint64 > 0 {
			res = append(res, valint64)
		}
	}
	return res
}

// 将string数组转成int64 数组,转换失败的忽略

func StrArrToInt64Arr(strs []string) []int64 {
	res := make([]int64, 0)
	for _, str := range strs {
		r, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			res = append(res, r)
		}
	}
	return res
}

/*
 * 字符串中ID以,分割的ids
 */
func ToStringArr(str string) []string {
	res := make([]string, 0)
	idSplit := strings.Split(str, common.CommaSeparator)

	for _, val := range idSplit {
		if len(val) > 0 {
			res = append(res, val)
		}
	}
	return res
}

func ArrToString(arr []string) string {
	return strings.Join(arr, common.CommaSeparator)
}

func FromInt(val int) string {
	return strconv.Itoa(val)
}

func ToInt(str string) int {
	if b, err := strconv.ParseInt(str, 10, 64); err != nil {
		return 0
	} else {
		return int(b)
	}
}

func FromByte(val byte) string {
	return strconv.Itoa(int(val))
}

func ToByte(str string) byte {
	if b, err := strconv.ParseInt(str, 10, 8); err != nil {
		return 0
	} else {
		return byte(b)
	}
}

func FromBool(b bool) string {
	if b {
		return "true"
	} else {
		return "false"
	}
}

func ToBool(s string) bool {
	lower := strings.ToLower(s)
	if lower == "true" || lower == "yes" || lower == "1" {
		return true
	}
	return false
}

func ToStrList(str string) []string {
	if len(str) <= 1 {
		return nil
	} else if str == "null" {
		return nil
	}
	var v []string
	json.Unmarshal([]byte(str), &v)
	return v
}

func ToObject(obj interface{}, str string) {
	json.Unmarshal([]byte(str), obj)
}

func FromStrList(strs []string) string {
	if strs == nil {
		return "[]"
	}
	data, _ := json.Marshal(&strs)
	return string(data)
}

func FromObject(obj interface{}) string {
	data, _ := json.Marshal(obj)
	return string(data)
}

func ToSqlNullString(key string, isNull bool) sql.NullString {
	res := sql.NullString{
		String: key,
		Valid:  isNull,
	}
	return res
}

func RandomString(length int) string {
	if length < 1 {
		return ""
	}
	result := make([]string, length)
	for i := 0; i < length; i++ {
		result = append(result, charArr[rand.Intn(61)])
	}
	return strings.Join(result, "")
}

func Sub(str string, begin, length int) string {
	rs := []rune(str)
	lth := len(rs)

	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}

func ToLike(str string) string {
	return fmt.Sprintf("%%%s%%", str)
}

func UUID() string {
	return uuid.New()
}

func Diff(arr1 []string, arr2 []string) []string {
	res := make([]string, 0)
	for _, val1 := range arr1 {
		exist := false
		for _, val2 := range arr2 {
			if val1 == val2 {
				exist = true
				break
			}
		}
		if !exist {
			res = append(res, val1)
		}
	}
	return res
}

func FromChinese(str string) string {
	res := strings.Join(pinyin.LazyPinyin(str, pinyin.NewArgs()), "")
	if IsEmpty(res) {
		res = str
	}
	return res
}

func IsEmpty(str string) bool {
	return str == ""
}

func AsciiPinyin(s string) string {
	options := pinyin.NewArgs()
	options.Fallback = func(r rune, options pinyin.Args) []string {
		if r < 128 {
			return []string{string(r)}
		} else {
			return []string{}
		}
	}
	result := strings.Join(pinyin.LazyPinyin(s, options), "")
	return strings.ToLower(result)
}

func AsciiPinyinLower(s string) string {
	s = strings.ToLower(AsciiPinyin(s))
	r := regexp.MustCompile(`[^a-z0-9]`)
	s = r.ReplaceAllString(s, "")
	return s
}

func IsMobile(str string) bool {
	reg := regexp.MustCompile(common.MobileRegex)
	return reg.MatchString(str)
}

func ProcessMobile(mobile string) string {
	if mobile[0:1] != "+" {
		return mobile
	}

	mobileHead := mobile[0:3]
	if mobileHead == "+86" {
		return mobile[3:]
	}
	return mobile
}

func IsEmail(str string) bool {
	reg := regexp.MustCompile(common.EmailRegex)
	return reg.MatchString(str)
}

func TrimSpace(str string) string {
	return strings.Trim(str, " ")
}

func StrInArray(str string, strs []string) bool {
	ina := false
	for _, st := range strs {
		if st == str {
			ina = true
			break
		}
	}
	return ina
}

func Md5(str string) string {
	h := md5.New()
	return hex.EncodeToString(h.Sum([]byte(str)))
}

func Md5Object(obj interface{}) string {
	data, _ := json.Marshal(obj)
	return Md5(string(data))
}

func DecodePercent(s string) string {
	if decoded, err := url.Unescape(s); err != nil {
		return ""
	} else {
		return decoded
	}
}

func DecodeCsvArg(s string) []string {
	vs := strings.Split(DecodePercent(s), string(44))
	vs = sliceutil.RemoveEmptyStrings(vs)
	vs = sliceutil.RemoveDuplicateStrings(vs)
	sort.Strings(vs)
	return vs
}

func HMacMd5Sign(key string, str string) string {
	h := hmac.New(md5.New, []byte(key))
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func ContentHasBase64Image(con string) bool {
	for _, val := range base64Images {
		if strings.Index(con, val) > -1 {
			return true
		}
	}
	return false
}

func GbkLen(s string) int {
	data, err := ioutil.ReadAll(
		transform.NewReader(bytes.NewReader([]byte(s)),
			simplifiedchinese.GB18030.NewEncoder()))
	if err != nil {
		return len(s)
	} else {
		return len(data)
	}
}

func CompileUris(str string) []string {
	res := uriRep.FindAllString(str, -1)

	return res
}

func ToFloat64(str string) float64 {
	i, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return i
}

func IsInvalidDisplayName(displayName string) bool {
	reg := regexp.MustCompile(DisplayNameRegexRule)
	return !reg.MatchString(displayName)
}

// 正则过滤sql注入的方法
// 参数 : 要匹配的语句
func FilteredSQLInject(to_match_str string) bool {
	//过滤 ‘
	//ORACLE 注解 --  /**/
	//关键字过滤 update ,delete
	// 正则的字符串, 不能用 " " 因为" "里面的内容会转义
	str := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|delete|insert|trancate|declare|exec|count|into|drop|execute)\b)`
	re, err := regexp.Compile(str)
	if err != nil {
		return false
	}
	return re.MatchString(to_match_str)
}