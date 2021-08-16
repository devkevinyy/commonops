package dms_service

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/chujieyang/commonops/ops/opslog"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chujieyang/commonops/ops/conf"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type connectionManage struct {
	connectString          string
	databaseName           string
	dbClient               *sql.DB
	mutex                  *sync.Mutex
	concurrentOperationNum int // 允许的并发操作数
}

func (manage connectionManage) acquire() {
	manage.mutex.Lock()
	manage.concurrentOperationNum = manage.concurrentOperationNum + 1
	manage.mutex.Unlock()
}

func (manage connectionManage) release() {
	manage.mutex.Lock()
	manage.concurrentOperationNum = manage.concurrentOperationNum - 1
	manage.mutex.Unlock()
}

var globalMutex = sync.Mutex{}

var connectionMap map[string]connectionManage

func init() {
	connectionMap = map[string]connectionManage{}
}

func initDBClient(instanceType string, dbConnectionString string) (db *sql.DB, err error) {
	db, err = sql.Open(instanceType, dbConnectionString)
	if err != nil {
		return
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(1 * time.Hour)
	return
}

func getDBConnection(instanceType string, connection string, dbName string) (err error) {
	itemKey := fmt.Sprintf("%s-%s", connection, dbName)
	globalMutex.Lock()
	opslog.Info().Println("get db connection: ", itemKey)
	_, ok := connectionMap[itemKey]
	if !ok {
		opslog.Info().Println("当前无此库对应的连接，现在进行初始化...")
		driverName := "-"
		switch instanceType {
		case "mysql", "polardb":
			driverName = "mysql"
		case "sqlserver":
			driverName = "sqlserver"
		}
		dbClient, err1 := initDBClient(driverName, connection)
		if err1 != nil {
			err = err1
			return
		}
		connectionMap[itemKey] = connectionManage{
			connectString:          connection,
			databaseName:           dbName,
			concurrentOperationNum: 0,
			dbClient:               dbClient,
			mutex:                  &sync.Mutex{},
		}
		opslog.Info().Println("初始化完成")
	}
	globalMutex.Unlock()
	return
}

// 获取SQL语句的执行计划
func DmsExplain(instanceType string, connection string, dbName string, SQL string) (explainResult string, needDba bool, err error) {
	err = getDBConnection(instanceType, connection, dbName)
	if err != nil {
		return
	}
	itemKey := fmt.Sprintf("%s-%s", connection, dbName)
	connManage, ok := connectionMap[itemKey]
	if !ok {
		err = errors.New("未找到对应的连接")
		return
	}
	defer func() {
		connManage.release()
	}()

	connManage.acquire()

	if connManage.concurrentOperationNum > conf.DmsConcurrentOperationNum { // 同一个库允许的并发操作数
		err = errors.New(fmt.Sprintf("当前库允许的并发操作为 %d 个，请等待其他操作完成后再尝试", conf.DmsConcurrentOperationNum))
		return
	}

	rows, err := connManage.dbClient.Query(fmt.Sprintf("EXPLAIN %s", SQL))
	if err != nil {
		return
	}
	var result []map[string]string
	columns, err := rows.Columns()
	values := make([]sql.RawBytes, len(columns))
	scans := make([]interface{}, len(columns))
	for i := range values {
		scans[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			return
		}
		each := make(map[string]string)
		for i, col := range values {
			each[columns[i]] = string(col)
		}
		result = append(result, each)
	}
	for _, item := range result {
		scanRows, err := strconv.Atoi(item["rows"])
		if err != nil {
			opslog.Error().Println(err)
			needDba = true
			scanRows = 1000000
		}
		if scanRows > 20000 {
			needDba = true
		}
		explainRow := fmt.Sprintf("id: %s \r\n table: %s \r\n select_type: %s \r\n type: %s \r\n key: %s \r\n possible_keys: %s \r\n "+
			"ref: %s \r\n rows: %s \r\n filtered: %s \r\n extra: %s \r\n ----------- \r\n",
			item["id"], item["table"], item["select_type"], item["type"], item["key"], item["possible_keys"],
			item["ref"], item["rows"], item["filtered"], item["extra"])
		explainResult = fmt.Sprintf("%s %s", explainResult, explainRow)
	}
	return
}

func DmsQuery(instanceType string, connection string, dbName string, SQL string) (columns []string, result []map[string]string, duration int64, err error) {
	err = getDBConnection(instanceType, connection, dbName)
	if err != nil {
		return
	}
	itemKey := fmt.Sprintf("%s-%s", connection, dbName)
	connManage, ok := connectionMap[itemKey]
	if !ok {
		err = errors.New("未找到对应的连接")
		return
	}
	defer func() {
		connManage.release()
	}()

	connManage.acquire()

	if connManage.concurrentOperationNum > conf.DmsConcurrentOperationNum { // 同一个库允许的并发操作数
		err = errors.New(fmt.Sprintf("当前库允许的并发操作为 %d 个，请等待其他操作完成后再尝试", conf.DmsConcurrentOperationNum))
		return
	}

	startExecTime := time.Now()
	rows, err := connManage.dbClient.Query(SQL)
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	elapsed := time.Since(startExecTime)
	duration = int64(elapsed / time.Millisecond)
	columns, err = rows.Columns()
	values := make([]sql.RawBytes, len(columns))
	scans := make([]interface{}, len(columns))
	for i := range values {
		scans[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			return
		}
		each := make(map[string]string)
		for i, col := range values {
			each[columns[i]] = string(col)
			if string(col) == "\x01" {
				each[columns[i]] = "1"
			}
			if string(col) == "\x00" {
				each[columns[i]] = "0"
			}
		}
		result = append(result, each)
	}
	return
}

func DmsExec(sqlType string, instanceType string, connection string, dbName string, SQL string) (result sql.Result, duration int64, err error) {
	err = getDBConnection(instanceType, connection, dbName)
	if err != nil {
		return
	}
	itemKey := fmt.Sprintf("%s-%s", connection, dbName)
	connManage, ok := connectionMap[itemKey]
	if !ok {
		err = errors.New("未找到对应的连接")
		return
	}
	defer func() {
		connManage.release()
	}()

	connManage.acquire()

	if connManage.concurrentOperationNum > conf.DmsConcurrentOperationNum { // 同一个库允许的并发操作数
		err = errors.New(fmt.Sprintf("当前库允许的并发操作为 %d 个，请等待其他操作完成后再尝试", conf.DmsConcurrentOperationNum))
		return
	}

	startExecTime := time.Now()
	result, err = connManage.dbClient.Exec(SQL)
	elapsed := time.Since(startExecTime)
	duration = int64(elapsed / time.Millisecond)
	return
}

func DmsRollback(instanceType string, connection string, dbName string, SQL string, rollbackTableName string) (result sql.Result, duration int64, err error) {
	err = getDBConnection(instanceType, connection, dbName)
	if err != nil {
		return
	}
	itemKey := fmt.Sprintf("%s-%s", connection, dbName)
	connManage, ok := connectionMap[itemKey]
	if !ok {
		err = errors.New("未找到对应的连接")
		return
	}
	defer func() {
		connManage.release()
	}()

	connManage.acquire()

	if connManage.concurrentOperationNum > conf.DmsConcurrentOperationNum { // 同一个库允许的并发操作数
		err = errors.New(fmt.Sprintf("当前库允许的并发操作为 %d 个，请等待其他操作完成后再尝试", conf.DmsConcurrentOperationNum))
		return
	}

	startExecTime := time.Now()
	procedureSQL := strings.Replace(SQL, "'", "''", -1)
	if instanceType == "mysql" || instanceType == "polardb" {
		procedureSQL = fmt.Sprintf("call p_dms_rollback('%s', '%s')", procedureSQL, rollbackTableName)
	} else { // sqlserver
		procedureSQL = fmt.Sprintf("exec p_dms_rollback '%s', '%s'", procedureSQL, rollbackTableName)
	}
	opslog.Info().Println(procedureSQL)
	result, err = connManage.dbClient.Exec(procedureSQL)
	elapsed := time.Since(startExecTime)
	duration = int64(elapsed / time.Millisecond)
	return
}

func DmsSqlParseType(sqlInput string) (sqlType string, err error) {
	tmpSql := strings.ToLower(strings.TrimSpace(sqlInput))
	if len(tmpSql) == 0 {
		err = errors.New("输入的SQL内容为空!")
		return
	}
	sqlType = "not allow"
	if pos := strings.Index(tmpSql, "select"); pos == 0 {
		sqlType = "query"
	}
	if pos := strings.Index(tmpSql, "update"); pos == 0 {
		sqlType = "update"
	}
	if pos := strings.Index(tmpSql, "insert"); pos == 0 {
		sqlType = "insert"
	}
	if pos := strings.Index(tmpSql, "delete"); pos == 0 {
		sqlType = "delete"
	}
	if sqlInput == "not allow" {
		err = errors.New("不支持该SQL语句进行的操作")
	}
	return
}

func DmsSQLParser(dbType string, sqlInput string) (sqlType string, tableName string, err error) {
	formatSql := strings.Replace(sqlInput, "'", "\"", -1)
	if dbType == "polardb" {
		dbType = "mysql"
	}
	parseCmd := fmt.Sprintf("java -jar jars/sqlparser.jar %s '%s'", dbType, formatSql)
	cmd := exec.Command("/bin/sh", "-c", parseCmd)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if err != nil {
		opslog.Info().Println("cmd.Run() failed", errStr, outStr)
		err = errors.New(errStr + " " + outStr)
		return
	}
	parserResult := strings.Split(outStr, " ")
	sqlType = parserResult[1]
	tableName = parserResult[2]
	return
}
