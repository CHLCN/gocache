package logic

import (
	"context"
	"fmt"
	"gocache/appv0/db"
	"net/http"
	"strconv"
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
		_ = db.Rdb.Set(context.Background(), key, fmt.Sprint(ret), time.Second*30).Err()
	}

	// 返回值
	_, _ = fmt.Fprintf(writer, "这是查数据库的结果："+fmt.Sprint(ret))
}

func SetInfoV0(w http.ResponseWriter, r *http.Request) {
	// 为了便于测试，这里依然采用
	if r.Method != http.MethodGet {
		http.Error(w, "只支持 Get 请求", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	id := r.URL.Query().Get("id")

	//key := fmt.Sprintf("book_%s", id)
	key := fmt.Sprintf("xxcode")
	//双写策略，同时写Redis和DB
	_ = db.Rdb.Set(context.Background(), key, name, time.Second*30).Err()
	info := db.Info{}
	idInt, _ := strconv.Atoi(id)
	info.Save(idInt, name)
	// 向客户端发送响应
	fmt.Fprintf(w, "双写策略完成。")
}

func SetInfoV1(w http.ResponseWriter, r *http.Request) {
	// 为了便于测试，这里依然采用
	if r.Method != http.MethodGet {
		http.Error(w, "只支持 Get 请求", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	id := r.URL.Query().Get("id")

	//key := fmt.Sprintf("book_%s", id)
	key := fmt.Sprintf("xxcode")

	//双写策略，先写数据库，再写缓存
	info := db.Info{}
	idInt, _ := strconv.Atoi(id)
	info.Save(idInt, name)

	data := info.Get(idInt)
	err := db.Rdb.Set(context.Background(), key, data, time.Second*30).Err()
	if err != nil {
		fmt.Println(err.Error())
		_, _ = fmt.Fprintf(w, "双写策略失败！")
		return
	}
	// 向客户端发送响应
	fmt.Fprintf(w, "双写策略完成。")
}

func SetBookV2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持 Get 请求", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	id := r.URL.Query().Get("id")

	info := db.Info{}
	idInt, _ := strconv.Atoi(id)

	//启动事务
	db.DB.Begin()
	if err := info.Save(idInt, name); err != nil {
		db.DB.Rollback()
		_, _ = fmt.Fprintf(w, "数据库更新失败！")
		return
	}

	//key := fmt.Sprintf("book_%s", id)
	key := fmt.Sprintf("xxcode")
	if err := db.Rdb.Set(context.Background(), key, name, time.Second*5).Err(); err != nil {
		db.DB.Rollback()
		_, _ = fmt.Fprintf(w, "缓存更新失败！")
		return
	}
	db.DB.Commit()

	// 向客户端发送响应
	_, _ = fmt.Fprintf(w, "双写策略完成。")
}
