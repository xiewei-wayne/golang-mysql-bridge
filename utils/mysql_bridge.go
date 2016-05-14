
package utils

import (
    "database/sql"
    "log"
    "reflect"
    _ "github.com/go-sql-driver/mysql"
    "strings"
    "strconv"
    "time"
)

/****************************************************************
 *                      全局变量                                 *
 ****************************************************************/

/* 使用这个Db链接 */
var Db = getDB()

/* 数据库相关配置 */
var MySQLUsername string = "root"               // 用户名
var MySQLPassword string = "12345"              // 密码
var MySQLDatabase string = "dsj2"               // 数据库
var MySQLHostPort string = "localhost:3307"     // 主机
/* 数据库链接资源,主要用来链接数据库 */
var DatabaseSource string = MySQLUsername + ":" + MySQLPassword +
                            "@tcp(" + MySQLHostPort + ")/" +
                            MySQLDatabase + "?charset=utf8&loc=Asia%2FShanghai"

/****************************************************************
 *                          相关函数                             *
 ****************************************************************/

/* 得到数据库链接 */
func getDB() *sql.DB {
    db, err := sql.Open("mysql", DatabaseSource)

    if err != nil {
        log.Println(err)
        return nil
    }

    if db == nil {
        log.Println("数据库连接失败")
        return nil
    }

    return db
}

/****************************************************************
 *                      核心方法                                 *
 ****************************************************************/

// 根据一条SQL语句执行SELECT操作,返回所有查询到的数据,以map的形式存储
// 全部数据是一个map slice形式
// sql = "SELECT * FROM users"
func SelectSqlToStringMap(sql string, args ...interface{}) (resultRows []map[string]interface{}, err error) {
    db := getDB()

    sqlPrepare, err := db.Prepare(sql)
    if err != nil {
        return nil, err
    }
    defer sqlPrepare.Close()

    res, err := sqlPrepare.Query(args...)
    if err != nil {
        return nil, err
    }
    defer res.Close()

    fields, err := res.Columns()
    if err != nil {
        return nil, err
    }

    for res.Next() {
        result := make(map[string] interface{})
        var scanResultContainers []interface{}
        for i := 0; i < len(fields); i++ {
            var scanResultContainer interface{}
            scanResultContainers = append(scanResultContainers, &scanResultContainer)
        }

        if err := res.Scan(scanResultContainers...); err != nil {
            return nil, err
        }

        for ii, key := range fields {
            rawValue := reflect.Indirect(reflect.ValueOf(scanResultContainers[ii]))
            //if row is null then ignore
            if rawValue.Interface() == nil {
                continue
            }

            result[key] = RawValueToString(rawValue.Interface())
        }

        resultRows = append(resultRows, result)
    }

    return resultRows, nil
}

// 根据一条SQL语句执行UPDATE,INSERT,DELETE操作
func ExecSqlToStringMap(sqlString string, args ...interface{}) (resultRows []map[string]interface{}, err error) {

    resultRows = make([]map[string] interface{}, 0)
    resultMap := make(map[string] interface{})

    res, err := getDB().Exec(sqlString, args ...)
    if err != nil {
        resultMap["error"] = err.Error()
    } else {
        resultMap["RowsAffected"], err = res.RowsAffected()
        resultMap["LastInsertId"], err = res.LastInsertId()
    }

    resultRows = append(resultRows, resultMap)

    return
}

// 根据一条SQL语句执行UPDATE,INSERT,DELETE操作
// rowsAffected:影响的条数
// lastInsertId:最后插入的id
func ExecSql(sql string, args ...interface{}) (rowsAffected int64, lastInsertId int64, err error) {
    db := getDB()

    sqlPrepare, err := db.Prepare(sql)
    if err != nil {
        return 0, 0, err
    }
    defer sqlPrepare.Close()

    res, err := sqlPrepare.Exec(args...)
    if err != nil {
        return 0, 0, err
    }

    rowsAffected, _ = res.RowsAffected()
    lastInsertId, _ = res.LastInsertId()

    return
}

// 无序,json对象
// 同时发送多个SQL查询语句给后台
func ExecSqlBySqlMap(sqlObject map[string]string) (map[string]interface{}) {

    sqlResults := make(map[string]interface{})

    var resultRows []map[string]interface{}

    for key, sql := range sqlObject {

        sql = strings.TrimLeft(sql, " \n")

        sqlCmd := strings.ToLower( strings.Split(sql, " ")[0] )

        switch sqlCmd {

            case "select", "desc":
                resultRows, _ = SelectSqlToStringMap(sql)

            case "insert", "delete", "update":
                resultRows, _ = ExecSqlToStringMap(sql)
        }

        sqlResults[key] = resultRows
    }

    return sqlResults
}

// 有序,json数组
// 同时发送多个SQL查询语句给后台
func ExecSqlBySqlSlice(sqlObject []string) ([]interface{}) {

    sqlResults := make([]interface{}, 0)

    var resultRows []map[string]interface{}

    for _, sql := range sqlObject {

        sql = strings.TrimLeft(sql, " ")

        sqlCmd := strings.ToLower( strings.Split(sql, " ")[0] )

        switch sqlCmd {

        case "select":
            resultRows, _ = SelectSqlToStringMap(sql)

        case "insert", "delete", "update":
            resultRows, _ = ExecSqlToStringMap(sql)
        }

        sqlResults = append(sqlResults, resultRows)
    }

    return sqlResults
}

/****************************************************************
 *                      辅助方法                                 *
 ****************************************************************/

// 将不同类型的数据转换成string类型
func RawValueToString(rawValue interface{}) string {

    valueType := reflect.TypeOf(rawValue)
    value := reflect.ValueOf(rawValue)

    var str string
    switch reflect.TypeOf(rawValue).Kind() {
    case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        //fmt.Println("Int8")
        str = strconv.FormatInt(value.Int(), 10)
    case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        //fmt.Println("Uint8")
        str = strconv.FormatUint(value.Uint(), 10)
    case reflect.Float32, reflect.Float64:
        //fmt.Println("Float32")
        str = strconv.FormatFloat(value.Float(), 'f', -1, 64)
    case reflect.Slice:
        //fmt.Println("Slice")
        if valueType.Elem().Kind() == reflect.Uint8 {
            str = string(value.Bytes())
            break
        }
    case reflect.String:
        //fmt.Println("String")
        str = value.String()
    //时间类型
    case reflect.Struct:
        //fmt.Println("Struct")
        str = rawValue.(time.Time).Format("2006-01-02 15:04:05.000 -0700")
    case reflect.Bool:
        //fmt.Println("Bool")
        if value.Bool() {
            str = "1"
        } else {
            str = "0"
        }
    }

    return str
}
