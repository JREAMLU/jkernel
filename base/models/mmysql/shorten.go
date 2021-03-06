package mmysql

import (
	"errors"

	"github.com/JREAMLU/core/com"
	"github.com/JREAMLU/core/db/mysql"
	"github.com/JREAMLU/core/global"
	"github.com/JREAMLU/jkernel/base/models/mentity"
	"github.com/beego/i18n"
	"github.com/jinzhu/gorm"
)

// BASE mysql alias name
const BASE = "base"

// func ShortenIn(r mentity.Redirect) (uint64, error) {
// 	res := mysql.X.Create(&r)
// 	if res.Error != nil {
// 		return 0, res.Error
// 	}
// 	return r.ID, nil
// }

// ShortenIn insert shorten
func ShortenIn(r mentity.Redirect) (uint64, error) {
	x, err := mysql.GetXS(BASE)
	if err != nil {
		return 0, err
	}
	res := x.Create(&r)
	if res.Error != nil {
		return 0, res.Error
	}
	return r.ID, nil
}

// ShortenInBatch insert shorten batch
func ShortenInBatch(redirects []mentity.Redirect, tx *gorm.DB) error {
	if len(redirects) == 0 {
		return errors.New(i18n.Tr(global.Lang, "url.SHORTENINBATCHILLEGAL"))
	}
	sql := `
INSERT INTO redirect
(long_url, short_url, long_crc, short_crc, status, created_by_ip, updated_by_ip, created_at, updated_at)
VALUES
`
	var params []interface{}
	for k, redirect := range redirects {
		fsql := `(?, ?, ?, ?, ?, ?, ?, ?, ?)`
		sql = com.StringJoin(sql, fsql)
		if k+1 != len(redirects) {
			sql = com.StringJoin(sql, ",")
		}
		params = append(params, redirect.LongURL)
		params = append(params, redirect.ShortURL)
		params = append(params, redirect.LongCrc)
		params = append(params, redirect.ShortCrc)
		params = append(params, redirect.Status)
		params = append(params, redirect.CreatedByIP)
		params = append(params, redirect.UpdateByIP)
		params = append(params, redirect.CreateAT)
		params = append(params, redirect.UpdateAT)
	}

	if err := tx.Exec(sql, params...).Error; err != nil {
		return err
	}
	return nil
}

// GetShortens select shorten list
func GetShortens(longCRC []uint64) (r []mentity.Redirect, err error) {
	sql := `
SELECT  redirect_id, long_url, short_url, long_crc, short_crc, status, 
        created_by_ip, updated_by_ip, created_at, updated_at
FROM    redirect
WHERE   long_crc IN (?)
`

	x, err := mysql.GetXS(BASE)
	if err != nil {
		return nil, err
	}

	res := x.Raw(sql, longCRC).Scan(&r)
	if res.Error != nil {
		return r, res.Error
	}
	return r, nil
}

// func GetShortens(longCRC []uint64) (r []mentity.Redirect, err error) {
// 	sql := `
// SELECT  redirect_id, long_url, short_url, long_crc, short_crc, status,
//         created_by_ip, updated_by_ip, created_at, updated_at
// FROM    redirect
// WHERE   long_crc IN (?)
// `
//
// 	res := mysql.X.Raw(sql, longCRC).Scan(&r)
// 	if res.Error != nil {
// 		return r, res.Error
// 	}
// 	return r, nil
// }
