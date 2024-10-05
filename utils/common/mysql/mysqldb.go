package mysql

import (
	"database/sql"
	"tzgit.kaixinxiyou.com/utils/common/db"

	"time"

	_ "github.com/go-sql-driver/mysql"

	"tzgit.kaixinxiyou.com/utils/common/log"
)

var (
	//Context    *sql.DB
	//PubContext *sql.DB
	contexts = make(map[string]*sql.DB)
)

func Init(dbAddr string, dbPoolSize int32) {
	InitName("", dbAddr, dbPoolSize)
}
func InitName(name string, dbAddr string, dbPoolSize int32) {
	// dbAddr = "root:root@tcp(192.168.22.212:3306)/x_game_s1?charset=utf8&parseTime=True&loc=Local"
	log.Release("初始化Mysql数据库")
	_, ok := contexts[name]
	if ok {
		log.Release("已初始化Mysql数据库")
		return
	}
	Context, err := sql.Open("mysql", dbAddr)

	if err != nil {
		log.Fatal("mysqldb init is error(%v)", err)
	}
	contexts[name] = Context

	log.Release("初始化Mysql数据库完成")
	Context.SetMaxOpenConns(int(dbPoolSize))
	Context.SetMaxIdleConns(int(dbPoolSize) / 50)
	Context.SetConnMaxLifetime(time.Hour)

}
func getDbName() string {
	r, e := QueryRow("select database() as dbname;")
	if e != nil {
		log.Error("getDbName:%v", e)
		return ""
	}
	return r.GetString("dbname")
}

func Query(strsql string, args ...interface{}) ([]*db.DataRow, error) {
	return db.Query(contexts[""], strsql, args...)
}
func QueryRow(strsql string, args ...interface{}) (*db.DataRow, error) {
	return db.QueryRow(contexts[""], strsql, args...)
}
func DbQueryRow(context *sql.DB, strsql string, args ...interface{}) (*db.DataRow, error) {
	return db.QueryRow(context, strsql, args...)
}
func PubQueryRow(name string, strsql string, args ...interface{}) (*db.DataRow, error) {
	return db.QueryRow(contexts[name], strsql, args...)
}
func PubQuery(name string, strsql string, args ...interface{}) ([]*db.DataRow, error) {
	return db.Query(contexts[name], strsql, args...)
}
func PubQueryWithoutPrepare(name string, strsql string, args ...interface{}) ([]*db.DataRow, error) {
	return db.QueryWithoutPrepare(contexts[name], strsql, args...)
}
func Exec(query string, args ...interface{}) (sql.Result, error) {
	//log.Debug("exec:%v", query)
	return contexts[""].Exec(query, args...)
}
func PubExec(name string, query string, args ...interface{}) (sql.Result, error) {
	//log.Debug("exec:%v", query)
	return contexts[name].Exec(query, args...)
}
func Exist(name string) bool {
	_, ok := contexts[name]
	return ok
}
