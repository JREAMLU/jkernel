package test

import (
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	_ "github.com/JREAMLU/jkernel/base/routers"
	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

func TestUrlGoshorten(t *testing.T) {
	r := TRollingCurl(Requests{
		Method: "POST",
		UrlStr: "/v1/url/goshorten.json?a=1&b=2",
		Header: map[string]string{
			"Content-Type": "application/json;charset=UTF-8;",
			"Accept":       "application/json",
			"Source":       "gotest",
			"ip":           "9.9.9.9",
		},
		Raw: `{"data":{"urls":[{"long_url":"http://o9d.cn","IP":"127.0.0.1"},{"long_url":"http://huiyimei.com","IP":"192.168.1.1"}],"timestamp":1466668134,"sign":"0B490F84305C7CF4D9CDD293B936BE0D"}}`,
	})
	w := httptest.NewRecorder()

	beego.BeeApp.Handlers.ServeHTTP(w, r)
	beego.Trace("testing", "TestUrlGoshorten", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("func Goshorten", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}
