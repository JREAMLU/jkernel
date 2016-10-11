package mmysql

import (
	"git.corp.plu.cn/phpgo/core/mysql"
	"github.com/JREAMLU/jkernel/base/models/mentity"
)

func ShortenIn(r mentity.Redirect) (uint64, error) {
	res := mysql.X.Create(&r)
	if res.Error != nil {
		return 0, res.Error
	}
	return r.ID, nil
}
