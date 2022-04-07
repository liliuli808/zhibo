package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type Mysql struct {
	MysqlDb    *sql.DB
	MysqlDbErr error
	Config     *Config
}

// Init 初始化链接
func (mysql *Mysql) Init() {
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		mysql.Config.UserName,
		mysql.Config.PassWord,
		mysql.Config.Host,
		mysql.Config.Port,
		mysql.Config.Database,
		mysql.Config.Charset)

	// 打开连接失败
	mysql.MysqlDb, mysql.MysqlDbErr = sql.Open("mysql", dbDSN)
	//defer MysqlDb.Close();
	if mysql.MysqlDbErr != nil {
		log.Println("dbDSN: " + dbDSN)
		panic("数据源配置不正确: " + mysql.MysqlDbErr.Error())
	}

	// 最大连接数
	mysql.MysqlDb.SetMaxOpenConns(100)
	// 闲置连接数
	mysql.MysqlDb.SetMaxIdleConns(20)
	// 最大连接周期
	mysql.MysqlDb.SetConnMaxLifetime(100 * time.Second)

	if mysql.MysqlDbErr = mysql.MysqlDb.Ping(); nil != mysql.MysqlDbErr {
		panic("数据库链接失败: " + mysql.MysqlDbErr.Error())
	}

}
