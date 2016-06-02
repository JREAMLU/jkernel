package controllers

import (
	"core/global"
	"fmt"

	//"github.com/astaxie/beego/validation"
)

type TesController struct {
	global.BaseController
}

func (t *TesController) Post() {
	fmt.Println(2)
}

//struct

//type SuperMan struct {
//	Skill string `valid:"Required;MinSize(5)"`
//}

//type Humem struct {
//	Ip string `valid:"IP"`
//}

//type User struct {
//	Name      string `valid:"Required"`
//	Age       int    `valid:"Range(1,18)"`
//	Humems    Humem
//	SuperMans []SuperMan
//}

//func (t *TesController) Post() {
//	valid := validation.Validation{}
//	u := User{Name: "man", Age: 18, Humems: Humem{Ip: "123.1.1.1"}, SuperMans: []SuperMan{SuperMan{Skill: "abc"}, SuperMan{Skill: "def"}}}

//	//u := User{SuperMans: []SuperMan{SuperMan{Skill: "abc"}, SuperMan{Skill: "def"}}}

//	//struct []
//	/*
//		u := User{
//			SuperMans: []SuperMan{
//				SuperMan{Skill: "abc"},
//				SuperMan{Skill: "def"},
//			},
//		}
//	*/

//	fmt.Println(u)

//	b, err := valid.Valid(&u.SuperMans)

//	fmt.Println(b, "--", err)

//	if err != nil {
//		fmt.Println("-----err了-----")
//	}

//	if !b {
//		for _, err := range valid.Errors {
//			fmt.Println(err.Key, ": ", err.Message)
//		}
//	}

//	t.Data["json"] = "abc"
//	t.ServeJSON()
//}
