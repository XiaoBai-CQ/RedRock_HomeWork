package flag

import (
	"RR6/lv1Andlv2Andlv3/global"
	"RR6/lv1Andlv2Andlv3/models"
	"fmt"
)

func DatabaseAutoMigrate() {
	var err error

	//自动建表**
	err = global.DB.Set("gorm:table_option", "Engine=InnoDB").
		AutoMigrate(
			&models.Student{},
			&models.User{},
		)

	if err != nil {
		fmt.Println("自动建表失败")
	} else {
		fmt.Println("自动建表成功")
	}
}
