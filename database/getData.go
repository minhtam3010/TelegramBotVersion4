package database

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func GetData(db *gorm.DB, myStruct []interface{}, tableName []string, option string) []string {
	var (
		query_insert = "date_created >= DATE(NOW())"
		query_update = "date_updated >= DATE(NOW()) and active = 1"
		query_delete = "date_updated >= DATE(NOW()) and active = 0"
		res          []string
	)
	if len(myStruct) != len(tableName) {
		return []string{errors.New("your input is not equal either the struct or table name").Error()}
	}
	switch len(myStruct) {
	case 0:
		return []string{errors.New("not found any struct or table name in your input").Error()}
	case 1:
		switch option {
		case "insert":
			db.Table(tableName[0]).Select("*").Where(query_insert).Find(&myStruct[0])
			res := fmt.Sprintf("%v", myStruct)
			return []string{res}
		case "update":
			db.Table(tableName[0]).Select("*").Where(query_update).Find(&myStruct[0])
			res := fmt.Sprintf("%v", myStruct)
			return []string{res}
		case "delete":
			db.Table(tableName[0]).Select("*").Where(query_delete).Find(&myStruct[0])
			res := fmt.Sprintf("%v", myStruct)
			return []string{res}
		}
	default:
		switch option {
		case "insert":
			for i, inter := range myStruct {
				db.Table(tableName[i]).Select("*").Where(query_insert).Find(&inter)
				res = append(res, fmt.Sprintf("%v", inter))
			}
		case "update":
			for i, inter := range myStruct {
				db.Table(tableName[i]).Select("*").Where(query_update).Find(&inter)
				res = append(res, fmt.Sprintf("%v", inter))
			}
		case "delete":
			for i, inter := range myStruct {
				db.Table(tableName[i]).Select("*").Where(query_delete).Find(&inter)
				res = append(res, fmt.Sprintf("%v", inter))
			}
		}
	}
	return res
}

func CallBotInserted(db *gorm.DB, myStruct []interface{}, tableName []string) []string {
	if len(myStruct) != len(tableName) {
		return []string{errors.New("your input is not equal either the struct or table name").Error()}
	}
	switch len(myStruct) {
	case 0:
		return []string{errors.New("not found any struct or table name in your input").Error()}
	case 1:
		db.Table(tableName[0]).Select("*").Where("date_created >= DATE(NOW())").Find(&myStruct[0])
		resInsert := fmt.Sprintf("%v", myStruct)
		return []string{resInsert}
	default:
		var (
			resInsert []string
		)
		for i, inter := range myStruct {
			db.Table(tableName[i]).Select("*").Where("date_created >= DATE(NOW())").Find(&inter)
			resInsert = append(resInsert, fmt.Sprintf("%v", inter))
		}
		return resInsert
	}
}

func CallBotUpdated(db *gorm.DB, myStruct []interface{}, tableName []string) []string {
	if len(myStruct) != len(tableName) {
		return []string{errors.New("your input is not equal either the struct or table name").Error()}
	}
	switch len(myStruct) {
	case 0:
		return []string{errors.New("not found any struct or table name in your input").Error()}
	case 1:
		db.Table(tableName[0]).Select("*").Where("date_updated >= DATE(NOW()) and active = 1").Find(&myStruct[0])
		resUpdate := fmt.Sprintf("%v", myStruct)
		return []string{resUpdate}
	default:
		var (
			resUpdate []string
		)
		for i, inter := range myStruct {
			db.Table(tableName[i]).Select("*").Where("date_updated >= DATE(NOW()) and active = 1").Find(&inter)
			resUpdate = append(resUpdate, fmt.Sprintf("%v", inter))
		}
		return resUpdate
	}
}

func CallBotDeleted(db *gorm.DB, myStruct []interface{}, tableName []string) []string {
	if len(myStruct) != len(tableName) {
		return []string{errors.New("your input is not equal either the struct or table name").Error()}
	}
	switch len(myStruct) {
	case 0:
		return []string{errors.New("not found any struct or table name in your input").Error()}
	case 1:
		db.Table(tableName[0]).Select("*").Where("date_updated >= DATE(NOW()) and active = 0").Find(&myStruct[0])
		resDelete := fmt.Sprintf("%v", myStruct)
		return []string{resDelete}
	default:
		var (
			resDelete []string
		)
		for i, inter := range myStruct {
			db.Table(tableName[i]).Select("*").Where("date_updated >= DATE(NOW()) and active = 0").Find(&inter)
			resDelete = append(resDelete, fmt.Sprintf("%v", inter))
		}
		return resDelete
	}
}
