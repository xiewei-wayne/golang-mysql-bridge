package main

import (
    "fmt"
    "mysql_bridge/utils"
    "encoding/json"
)

func main()  {
    execSqlMap()
    execSqlSlice()
    execSqlByParam()
}

// 无序sql
// 支持SELECT,DELETE,UPDATE,INSERT语句
// SELECT返回查询到的数据
// DELETE,UPDATE,INSERTF返回影响的条数
// 可以通过userList,deleteUser,updateUser,insertUser获取对应的值
func execSqlMap() {
    sqls := make(map[string]string)
    sqls["userList"] = "SELECT * FROM user"
    sqls["deleteUser"] = "DELETE FROM user WHERE id=1"
    sqls["updateUser"] = "UPDATE user set name='zhangsan' WHERE id=4"
    sqls["insertUser"] = "INSERT INTO user(name) VALUES('lisi')"

    sqlResult := utils.ExecSqlBySqlMap(sqls)

    jsonResult, _ := json.Marshal(sqlResult)

    fmt.Println(string(jsonResult))
}
/*
{
    "deleteUser": [
        {
            "LastInsertId": 0,
            "RowsAffected": 0
        }
    ],
    "insertUser": [
        {
            "LastInsertId": 11,
            "RowsAffected": 1
        }
    ],
    "updateUser": [
        {
            "LastInsertId": 0,
            "RowsAffected": 0
        }
    ],
    "userList": [
        {
            "id": "3",
            "name": "Jack"
        },
        {
            "id": "4",
            "name": "zhangsan"
        },
        {
            "id": "5",
            "name": "Tom"
        },
        {
            "id": "6",
            "name": "lisi"
        },
        {
            "id": "7",
            "name": "lisi"
        },
        {
            "id": "8",
            "name": "lisi"
        },
        {
            "id": "9",
            "name": "lisi"
        },
        {
            "id": "10",
            "name": "lisi"
        },
        {
            "id": "11",
            "name": "lisi"
        }
    ]
}
 */

// 有序sql
// 支持SELECT,DELETE,UPDATE,INSERT语句
// SELECT返回查询到的数据
// DELETE,UPDATE,INSERTF返回影响的条数
// 可以通过索引0,1,2,3得到sql语句执行的对应的值
func execSqlSlice() {
    sqls := make([]string, 0)
    sqls = append(sqls, "SELECT * FROM user")
    sqls = append(sqls, "DELETE FROM user WHERE id=1")
    sqls = append(sqls, "UPDATE user set name='zhangsan' WHERE id=4")
    sqls = append(sqls, "INSERT INTO user(name) VALUES('lisi')")

    sqlResult := utils.ExecSqlBySqlSlice(sqls)

    jsonResult, _ := json.Marshal(sqlResult)

    fmt.Println(string(jsonResult))
}
/*
[
    [
        {
            "id": "3",
            "name": "Jack"
        },
        {
            "id": "4",
            "name": "zhangsan"
        },
        {
            "id": "5",
            "name": "Tom"
        },
        {
            "id": "6",
            "name": "lisi"
        },
        {
            "id": "7",
            "name": "lisi"
        },
        {
            "id": "8",
            "name": "lisi"
        },
        {
            "id": "9",
            "name": "lisi"
        },
        {
            "id": "10",
            "name": "lisi"
        },
        {
            "id": "11",
            "name": "lisi"
        },
        {
            "id": "12",
            "name": "lisi"
        }
    ],
    [
        {
            "LastInsertId": 0,
            "RowsAffected": 0
        }
    ],
    [
        {
            "LastInsertId": 0,
            "RowsAffected": 0
        }
    ],
    [
        {
            "LastInsertId": 13,
            "RowsAffected": 1
        }
    ]
]
 */

// 执行带参数的sql
// 支持UPDATE,INSERT,DELETE语句
// 下面这个例子是删除id未6的数据,返回数据量,有id=6的,值不为0,为删除的数量
func execSqlByParam() {
    sql := "DELETE FROM user WHERE id=?"
    sqlParam := make([]interface{}, 0)
    sqlParam = append(sqlParam, 6)

    rowsAffected, lastInsertId, _ := utils.ExecSql(sql, sqlParam...)

    fmt.Println("rowsAffected:", rowsAffected, "lastInsertId:", lastInsertId)
}
/*
rowsAffected: 1 lastInsertId: 0
 */