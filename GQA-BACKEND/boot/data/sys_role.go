package data

import (
	"fmt"
	"gin-quasar-admin/global"
	"gin-quasar-admin/model/system"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

var SysRole = new(sysRole)

type sysRole struct{}

var sysRoleData = []system.SysRole{
	{
		GqaModel: global.GqaModel{
			Id:       1,
			Status:   "on",
			Sort:     1,
			Desc:     "超级管理员角色",
			CreateAt: time.Now(),
			CreateBy: "admin",
			UpdateAt: time.Now(),
		},
		RoleCode: "super-admin",
		RoleName: "超级管理组",
	},
}

func (s *sysRole) Init() error {
	return global.GqaDb.Transaction(func(tx *gorm.DB) error {
		var count int64
		tx.Model(&system.SysRole{}).Count(&count)
		if count != 0 {
			fmt.Println("[Gin-Quasar-Admin] --> sys_role 表的初始数据已存在！数据量：", count)
			global.GqaLog.Error("sys_role 表的初始数据已存在！", zap.Any("数据量", count))
			return nil
		}
		if err := tx.Create(&sysRoleData).Error; err != nil { // 遇到错误时回滚事务
			return err
		}
		fmt.Println("[Gin-Quasar-Admin] --> sys_role 表初始数据成功！")
		global.GqaLog.Error("sys_role 表初始数据成功！")
		return nil
	})
}