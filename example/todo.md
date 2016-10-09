SOA:
	1.long to short
		完整的url 经过算法 变成 几个字母
		返回完整的短url, 前缀放在数据库里
		判断
			a.是否存在
	2.short to long

	3.建表语句
		CREATE TABLE `redirect` (
			`redirect_id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '短网址唯一id,自增长',
			`long_url` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '原始url',
			`short_url` CHAR(25) NOT NULL DEFAULT '' COMMENT '短url',
			`long_crc` INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '原始url crc',
			`short_crc` INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '短url crc',
			`status` TINYINT(4) UNSIGNED NOT NULL DEFAULT 0 COMMENT '状态 0:删除 1:正常',
			`created_by_ip` INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建者ip',
			`updated_by_ip` INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新者ip',
			`created_at` INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间timestamp',
			`updated_at` INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间timestamp',
			PRIMARY KEY (`redirect_id`),
			KEY `long_crc` (`long_crc`),
			KEY `short_url` (`short_url`)
		) ENGINE = INNODB DEFAULT CHARSET=utf8 COMMENT='短网址表';

	4.redis
		全部走redis

技术疑问:
	0.路由					  			 √
		namespace		  				√
		疑问								√
		NSBefore 
		NSAfter
		路由域名						√
	1.参数以及验证
		验证json的合法性 						√
		多维的meta data 					  √
		json相关知识 				 		 √
		dataParams为数组: 定义struct  		 √
		接受的参数							√
		参数的验证							√
		映射到struct 							√
		struct里的验证 							√
		继承数组验证 								√
		错误机制的封装 						√
		Content-Type验证					√
2.baseConrtoller						√
3.filter
5.mysql 								√
6.redis
logs 									√
test
返回201 								√
签名验证
7.日志
	基于队列 Kakfa
	MongoDB
8.错误页面
9.代码提示 							√
10.封装返回 						√
11.ffjson							√

go:
	todo
		logs
			Kakfa->MongoDB
			封装kafak
		
		签名验证 
		思考request_id

	已完成
		sublime 配置保存 	//
		i18n 	//
		参数的验证 	//
		meta验证
			路由获取version 	//框架受限制
			把meta参数放入新的struct 返回 	//
		错误err 	//
		返回值封装 	//
			meta信息中的requestid 	错误时也返回requestid 	//
				map检查key是否存在 							//
				input 返回值改成 2个返回 					//

			路由匹配配置化 		//
			core单独移出 	 	//	
				目录结构变化 	//
		content-type判断 		//
		nginx环境下header 下划线_格式bug 					//
		mysql 							//
			gorm 							//
		redis 							//
		test 						//
		返回201 											//
		调用api方法封装 								//
		部署 												//
		