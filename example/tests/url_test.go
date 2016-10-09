package test

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/JREAMLU/core/inout"
	_ "github.com/JREAMLU/jkernel/base/routers"
	"github.com/JREAMLU/jkernel/base/services"
	"github.com/JREAMLU/jkernel/base/services/entity"
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

	Convey("func /v1/url/goshorten.json", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}

func TestUrlServiceGoshorten(t *testing.T) {
	httpStatus, shorten := urlServiceGoshorten()

	Convey("func Goshorten()", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(httpStatus, ShouldEqual, 200)
		})
		Convey("result", func() {
			datalist := shorten.Data.(entity.DataList)
			So(datalist.Total, ShouldEqual, 2)
			So(len(datalist.List["http://huiyimei.com"].(string)), ShouldBeGreaterThan, 0)
			So(len(datalist.List["http://o9d.cn"].(string)), ShouldBeGreaterThan, 0)
		})
	})
}

func BenchmarkUrlServiceGoshorten(b *testing.B) {
	Convey("bench UrlServiceGoshorten \n", b, func() {
		for i := 0; i < b.N; i++ {
			urlServiceGoshorten()
		}
	})
}

func urlServiceGoshorten() (int, inout.Output) {
	data := make(map[string]interface{})
	h := `{"Accept":["application/json"],"Content-Type":["application/json;charset=UTF-8;"],"Ip":["9.9.9.9"],"Request-Id":["base-57c930de30e8bd1aac000001"],"Source":["gotest"]}`
	hm := make(http.Header)
	hm["Accept"] = []string{"application/json"}
	hm["Content-Type"] = []string{"application/json;charset=UTF-8;"}
	hm["Ip"] = []string{"9.9.9.9"}
	hm["Request-Id"] = []string{"base-57c930de30e8bd1aac000001"}
	hm["Source"] = []string{"gotest"}
	b := `{"data":{"urls":[{"long_url":"http://o9d.cn","IP":"127.0.0.1"},{"long_url":"http://huiyimei.com","IP":"192.168.1.1"}],"timestamp":1466668134,"sign":"0B490F84305C7CF4D9CDD293B936BE0D"}}`
	c := `[]`
	q := `"a":["1"],"b":["2"]}`
	qm := make(map[string][]string)
	qm["a"] = []string{"1"}
	qm["b"] = []string{"2"}

	data["header"] = []byte(h)
	data["body"] = []byte(b)
	data["cookies"] = []byte(c)
	data["querystr"] = []byte(q)
	data["headermap"] = hm
	data["cookiesslice"] = []string{""}
	data["querystrmap"] = qm

	var service services.Url

	return service.GoShorten(data)
}
