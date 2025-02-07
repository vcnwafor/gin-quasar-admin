package private

import (
	"errors"
	"github.com/Junvary/gin-quasar-admin/GQA-BACKEND/global"
	"github.com/Junvary/gin-quasar-admin/GQA-BACKEND/model"
	"gorm.io/gorm"
)

type ServiceDept struct{}

func DeptList2DeptTree(deptList []model.SysDept, pCode string) []model.SysDept {
	var deptTree []model.SysDept
	for _, v := range deptList {
		if v.ParentCode == pCode {
			v.Children = DeptList2DeptTree(deptList, v.DeptCode)
			deptTree = append(deptTree, v)
		}
	}
	return deptTree
}

func (s *ServiceDept) GetDeptList(requestDeptList model.RequestGetDeptList) (err error, dept interface{}, total int64) {
	pageSize := requestDeptList.PageSize
	offset := requestDeptList.PageSize * (requestDeptList.Page - 1)
	db := global.GqaDb.Model(&model.SysDept{})
	var deptList []model.SysDept
	//配置搜索
	if requestDeptList.DeptCode != "" {
		db = db.Where("dept_code like ?", "%"+requestDeptList.DeptCode+"%")
	}
	if requestDeptList.DeptName != "" {
		db = db.Where("dept_name like ?", "%"+requestDeptList.DeptName+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(pageSize).Offset(offset).Order(model.OrderByColumn(requestDeptList.SortBy, requestDeptList.Desc)).
		Preload("LeaderUser").Find(&deptList).Error
	deptTree := DeptList2DeptTree(deptList, "")
	return err, deptTree, total
}

func (s *ServiceDept) EditDept(toEditDept model.SysDept) (err error) {
	var sysDept model.SysDept
	if sysDept.Stable == "yesNo_yes" {
		return errors.New("系统内置不允许编辑：" + toEditDept.DeptCode)
	}
	if err = global.GqaDb.Where("id = ?", toEditDept.Id).First(&sysDept).Error; err != nil {
		return err
	}
	//err = global.GqaDb.Updates(&toEditDept).Error
	err = global.GqaDb.Save(&toEditDept).Error
	return err
}

func (s *ServiceDept) AddDept(toAddDept model.SysDept) (err error) {
	var dept model.SysDept
	if !errors.Is(global.GqaDb.Where("dept_code = ?", toAddDept.DeptCode).First(&dept).Error, gorm.ErrRecordNotFound) {
		return errors.New("此部门已存在：" + toAddDept.DeptCode)
	}
	err = global.GqaDb.Create(&toAddDept).Error
	return err
}

func (s *ServiceDept) DeleteDeptById(id uint) (err error) {
	var sysDept model.SysDept
	if sysDept.Stable == "yesNo_yes" {
		return errors.New("系统内置不允许删除：" + sysDept.DeptCode)
	}
	if err = global.GqaDb.Where("id = ?", id).First(&sysDept).Error; err != nil {
		return err
	}
	if err = global.GqaDb.Where("id = ?", id).Unscoped().Delete(&sysDept).Error; err != nil {
		return err
	}
	var sysDeptUser model.SysDeptUser
	err = global.GqaDb.Where("sys_dept_dept_code = ?", sysDept.DeptCode).Delete(&sysDeptUser).Error
	return err
}

func (s *ServiceDept) QueryDeptById(id uint) (err error, deptInfo model.SysDept) {
	var dept model.SysDept
	err = global.GqaDb.Preload("LeaderUser").Preload("CreatedByUser").Preload("UpdatedByUser").First(&dept, "id = ?", id).Error
	return err, dept
}

func (s *ServiceDept) QueryUserByDept(deptCode *model.RequestDeptCode) (err error, user []model.SysUser) {
	dept := model.SysDept{
		DeptCode: deptCode.DeptCode,
	}
	var userList []model.SysUser
	err = global.GqaDb.Model(&dept).Association("Staff").Find(&userList)
	return err, userList
}

func (s *ServiceDept) RemoveDeptUser(toDeleteDeptUser *model.RequestDeptUser) (err error) {
	var deptUser model.SysDeptUser
	err = global.GqaDb.Where("sys_dept_dept_code = ? and sys_user_username = ?", toDeleteDeptUser.DeptCode, toDeleteDeptUser.Username).Delete(&deptUser).Error
	return err
}

func (s *ServiceDept) AddDeptUser(toAddDeptUser *model.RequestDeptUserAdd) (err error) {
	var deptUser []model.RequestDeptUser
	for _, r := range toAddDeptUser.Username {
		ud := model.RequestDeptUser{
			Username: r,
			DeptCode: toAddDeptUser.DeptCode,
		}
		deptUser = append(deptUser, ud)
	}
	if len(deptUser) != 0 {
		err = global.GqaDb.Model(&model.SysDeptUser{}).Save(&deptUser).Error
		return err
	} else {
		return errors.New("本次操作没有影响！")
	}
}
