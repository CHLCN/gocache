package logic

import (
	"context"
	"fmt"
	"gocache/appv0/db"
	"net/http"
	"time"
)

func GetInfo(writer http.ResponseWriter, request *http.Request) {
	// 必须是GET请求
	if request.Method != http.MethodGet {
		http.Error(writer, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 查看缓存是否存在
	key := fmt.Sprintf("xxcode")
	cache, _ := db.Rdb.Get(context.Background(), key).Result()
	if cache != "" {
		_, _ = fmt.Fprintf(writer, "这是查缓存的结果："+fmt.Sprint(cache))
		return
	}

	// 查询数据库
	t := db.Info{}
	ret := t.Get(1)

	// 设置缓存ret
	if ret.ID > 0 {
		_ = db.Rdb.Set(context.Background(), key, fmt.Sprint(ret), time.Second*5).Err()
	}

	// 返回值
	_, _ = fmt.Fprintf(writer, "这是查数据库的结果："+fmt.Sprint(ret))
}
