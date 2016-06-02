package controllers

import (
	//"encoding/json"
	"fmt"
	"hash/crc32"
	"hash/crc64"

	"core/global"
)

type LongController struct {
	global.BaseController
}

/**
 *	@auther	jream.lu
 *	@url		https://r.jream.lu/v1/golong.json
 */
func (gl *LongController) GoLong() {
	//var user models.User
	//json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	//uid := models.AddUser(user)
	//u.Data["json"] = map[string]string{"uid": uid}
	//u.ServeJSON()

	//var domain models.Domain
	//json.Unmarshal(st.Ctx.Input.RequestBody, &domain)
	//fmt.Println("==========", &domain)

	var url string = "http://o9d.cn"
	a := Jcrc32(url)
	fmt.Println(a)

	gl.Data["json"] = map[string]string{"name": "long"}
	gl.ServeJSON()
}

func Jcrc32(str string) uint32 {
	h := crc32.NewIEEE()
	h.Write([]byte(str))
	return h.Sum32()
}

func Jcrc64(str string) uint64 {
	h := crc64.New(crc64.MakeTable(crc64.ISO))
	h.Write([]byte(str))
	return h.Sum64()
}
