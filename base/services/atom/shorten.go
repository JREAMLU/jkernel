/**
 *	@Auther		    jream.lu
 *	@Description    调用外部或内部方法或接口
 */

package atom

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

const (
	Zero = iota
	One
	Two
	Three
	Four
)

/**
 *	@auther		jream.lu
 *	@intro		获取短网址链接
 *	@logic
 *	@todo
 *	@params		url, shortenDomain string   原始url, 短网址域名
 *	@return 	string
 */
func GetShortenUrl(url, shortenDomain string) string {
	//随机
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	random := r.Intn(Four)

	shortenUrl := ShortenUrl(url)
	shortenUrlParams := shortenUrl[random]

	sUrl := shortenDomain + shortenUrlParams
	return sUrl
}

/**
 *	@auther		jream.lu
 *	@intro		短网址算法
 *	@logic
 *	@todo
 *	@params		url string  原始url
 *	@return 	[]string
 */
func ShortenUrl(url string) []string {
	//62位的密码
	base62 := [...]string{"a", "b", "c", "d", "e", "f", "g", "h",
		"i", "j", "k", "l", "m", "n", "o", "p",
		"q", "r", "s", "t", "u", "v", "w", "x",
		"y", "z", "A", "B", "C", "D", "E", "F",
		"G", "H", "I", "J", "K", "L", "M", "N",
		"O", "P", "Q", "R", "S", "T", "U", "V",
		"W", "X", "Y", "Z", "0", "1", "2", "3",
		"4", "5", "6", "7", "8", "9"}

	//传进来的url进行md5
	h := md5.New()
	h.Write([]byte(url))

	//hex.EncodeToString(h.Sum(nil)) -> md5
	hex := hex.EncodeToString(h.Sum(nil))
	hexLen := len(hex)
	subHexLen := hexLen / 8

	var outPut []string
	for i := 0; i < subHexLen; i++ {
		//截取md5的长度 按8个 截4段
		subHex := SubString(hex, i*8, 8)

		//每段与0x组成16进制
		strHex := "0x" + subHex

		//将string转换成int 16进制
		sH, _ := strconv.ParseInt(strHex, 0, 64)

		//做运算
		valInt := 0x3FFFFFFF & (1 * sH)

		var out string
		for j := 0; j < 6; j++ {
			val := 0x0000003D & valInt
			out += base62[val]
			valInt = valInt >> 5
			//			fmt.Println(valInt)
		}

		//添加到outPut 切片里
		outPut = append(outPut, out)
		//		fmt.Println(outPut)
	}

	return outPut
}

/**
 *	@auther		jream.lu
 *	@intro		截取字符串长度
 *	@logic
 *	@todo
 *	@params		str string
 *              begin, length int
 *	@return 	string
 */
func SubString(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
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
