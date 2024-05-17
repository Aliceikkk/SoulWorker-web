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

type IPInfo struct {
	Count      int
	LastTime   time.Time
	BlockUntil time.Time
}

var ipMap = make(map[string]*IPInfo)
var mutex = &sync.Mutex{}

func main() {
	logger.InitLogger() // 初始化日志

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

	// 连接 MySQL 数据库
	mysqlServer := "127.0.0.1"
	mysqlUser := "root"
	mysqlPassword := "d1ed1f74cb1ef5e0"
	mysqlDatabase := "IP_Info"
	mysqlConnString := fmt.Sprintf("%s:%s@tcp(%s)/%s", mysqlUser, mysqlPassword, mysqlServer, mysqlDatabase)
	dbMySQL, err := sql.Open("mysql", mysqlConnString)
	if err != nil {
		logger.Error("发生错误：", err)
	}
	defer dbMySQL.Close()

	// 设置路由处理器
	http.HandleFunc("/api/reg", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		if r.Method != http.MethodPost {
			json.NewEncoder(w).Encode(Response{Msg: "曹尼玛傻逼"}) // 如果不是post请求就骂他
			return
		}

		ip := strings.Split(r.RemoteAddr, ":")[0]

		mutex.Lock()
		info, ok := ipMap[ip]
		if !ok {
			info = &IPInfo{
				Count:    0,
				LastTime: time.Now(),
			}
			ipMap[ip] = info
		}

		// 检查这个 IP 地址是否被封禁
		if time.Now().Before(info.BlockUntil) {
			json.NewEncoder(w).Encode(Response{Code: "429", Msg: "Too many requests, please try again later."})
			mutex.Unlock()
			return
		}

		if time.Since(info.LastTime) > 1*time.Minute {
			info.Count = 0
		}

		info.Count++

		// 如果在 1 分钟内请求了 10 次，封禁这个 IP 地址 60 分钟
		if info.Count >= 10 {
			info.BlockUntil = time.Now().Add(60 * time.Minute)
			_, err = dbMySQL.Exec("INSERT INTO Ban_ip (IP, TIME) VALUES (?, ?)",
				ip, info.BlockUntil.Unix())
			if err != nil {
				logger.Error("发生错误：", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				mutex.Unlock()
				return
			}
			json.NewEncoder(w).Encode(Response{Code: "429", Msg: "Too many requests, please try again later."})
			mutex.Unlock()
			return
		}

		info.LastTime = time.Now()
		mutex.Unlock()

		// 记录日志
		logger.Info("Received a request from IP: ", ip)

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

		// 如果用户名已存在，返回错误信息
		if exists > 0 {
			json.NewEncoder(w).Encode(Response{Code: "401", Msg: "用户已存在"})
			return
		}

		// 检查 IP 是否已经注册过三个账号
		var num int
		err = dbMySQL.QueryRow("SELECT NUM FROM Info WHERE IP = ?", ip).Scan(&num)
		if err != nil && err != sql.ErrNoRows {
			logger.Error("发生错误：", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 如果 IP 已经注册过三个账号，返回错误信息
		if num >= 3 {
			json.NewEncoder(w).Encode(Response{Code: "429", Msg: "每个 IP 只能注册三个账号"})
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

		// 更新 IP 注册数量
		if num > 0 {
			_, err = dbMySQL.Exec("UPDATE Info SET NUM = NUM + 1 WHERE IP = ?", ip)
		} else {
			_, err = dbMySQL.Exec("INSERT INTO Info (IP, NUM) VALUES (?, 1)", ip)
		}
		if err != nil {
			logger.Error("发生错误：", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(Response{Code: "200", Msg: "注册成功"})
	})

	logger.Info("服务器正在监听 8080 端口")
	if err := http.ListenAndServe(":8080", nil); err != nil {
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
