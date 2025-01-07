package main

import (
	"SoulWorker/logger"
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

type Response struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func main() {
	logger.InitLogger()
	logger.Info("服务正在启动")

	// 连接 SQL Server 数据库
	server := "IP"
	user := "账号"
	password := "密码"
	database := "AccountDB"
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", server, user, password, database)
	dbSQLServer, err := sql.Open("mssql", connString)
	if err != nil {
		logger.Error("发生错误：", err)
	}
	defer dbSQLServer.Close()

	// 设置路由处理器
	http.HandleFunc("/api/reg", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		// 获取客户端IP
		clientIP := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			clientIP = forwarded
		}
		// 如果IP包含端口号，只取IP部分
		if strings.Contains(clientIP, ":") {
			clientIP = strings.Split(clientIP, ":")[0]
		}

		if r.Method != http.MethodPost {
			json.NewEncoder(w).Encode(Response{Msg: "请求方法错误"})
			return
		}

		err := r.ParseForm()
		if err != nil {
			logger.Error("发生错误：", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		username := r.Form.Get("username")
		password := r.Form.Get("password")

		// 检查用户名是否已存在
		var exists int
		err = dbSQLServer.QueryRow("SELECT COUNT(*) FROM AccountDB.dbo.TB_ACCOUNT WHERE ACCOUNT_ID = ?", username).Scan(&exists)
		if err != nil {
			logger.Error("发生错误：", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if exists > 0 {
			json.NewEncoder(w).Encode(Response{Code: "401", Msg: "用户已存在"})
			return
		}

		// 用户名不存在，进行注册
		sha1Password := sha1Encode(password)
		createDate := time.Now()

		_, err = dbSQLServer.Exec("INSERT INTO AccountDB.dbo.TB_ACCOUNT (ACCOUNT_ID, PASSWORD, LASTSERVER_ID, STATE, CREATE_DATE, CASH, CLEAR_TUTORIAL, ECHELON_LEVEL, ECHELON_EXP, IP, SECOND_PW_CHECK, TRADE_PW_CHECK, MAC_ADDR, POWER, IP_STRING) VALUES (?, ?, 0, 0, ?, 99999999, 1, 0, 0, ?, 0, 0, '0', 0, ?)",
			username, sha1Password, createDate, getRandomIP(), "127.0.0.1")
		if err != nil {
			logger.Error("发生错误：", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 注册成功时记录日志
		logger.Info(fmt.Sprintf("IP:%s 注册成功 用户名：%s 密码：%s", clientIP, username, password))

		json.NewEncoder(w).Encode(Response{Code: "200", Msg: "注册成功"})
	})

	logger.Info("服务器正在监听 11451 端口")
	if err := http.ListenAndServe(":11451", nil); err != nil {
		logger.Error("发生错误：", err)
	}

	logger.CloseLogger()
}

func getRandomIP() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(90000000) + 10000000 
}

func sha1Encode(input string) []byte {
	sha1Bytes := sha1.Sum([]byte(input))
	return sha1Bytes[:]
}

func enableCors(w *http.ResponseWriter) { //跨域
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
