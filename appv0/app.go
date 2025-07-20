package appv0

import (
	"fmt"
	"gocache/appv0/db"
	"gocache/appv0/logic"
	"net/http"
)

func Run() {
	// 加载配置 连接数据库
	db.NewDb()
	db.NewRdb()
	// 使用go的http路由拉起服务

	http.HandleFunc("/get_name", logic.GetInfo)
	http.HandleFunc("/set_name", logic.SetInfoV1)

	fmt.Println("服务启动中...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Println("服务器启动失败！" + err.Error())
	}
}
