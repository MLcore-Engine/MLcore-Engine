package model

import (
	"MLcore-Engine/common"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func createRootAccountIfNeed() error {
	var user User
	//if user.Status != common.UserStatusEnabled {
	if err := DB.First(&user).Error; err != nil {
		common.SysLog("no user exists, create a root user for you: username is root, password is 123456")
		hashedPassword, err := common.Password2Hash("123456")
		if err != nil {
			return err
		}
		rootUser := User{
			Username:    "root",
			Password:    hashedPassword,
			Role:        common.RoleRootUser,
			Status:      common.UserStatusEnabled,
			DisplayName: "Root User",
		}
		DB.Create(&rootUser)
	}
	return nil
}

func CountTable(tableName string) (num int64) {
	DB.Table(tableName).Count(&num)
	return
}

func InitDB() (err error) {
	var db *gorm.DB
	if os.Getenv("SQL_DSN") != "" {
		// Use MySQL
		db, err = gorm.Open(mysql.Open(os.Getenv("SQL_DSN")), &gorm.Config{
			PrepareStmt: true, // precompile SQL
		})
	} else {
		// Use SQLite
		db, err = gorm.Open(sqlite.Open(common.SQLitePath), &gorm.Config{
			PrepareStmt: true, // precompile SQL
		})
		common.SysLog("SQL_DSN not set, using SQLite as database")
	}
	if err == nil {

		DB = db

		err := db.AutoMigrate(&File{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&User{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Option{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Notebook{})
		if err != nil {
			return err
		}

		if err := db.AutoMigrate(&Project{}); err != nil {
			return err
		}
		if err := db.AutoMigrate(&UserProject{}); err != nil {
			return err
		}

		if err := db.AutoMigrate(&TrainingJob{}); err != nil {
			return err
		}

		if err := db.AutoMigrate(&TritonDeploy{}); err != nil {
			return err
		}

		if err := db.AutoMigrate(&Dataset{}); err != nil {
			return err
		}

		if err := db.AutoMigrate(&DatasetVersion{}); err != nil {
			return err
		}

		if err := db.AutoMigrate(&DatasetEntry{}); err != nil {
			return err
		}

		err = createRootAccountIfNeed()
		return err
	} else {
		common.FatalLog(err)
	}
	return err
}

func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	return err
}
