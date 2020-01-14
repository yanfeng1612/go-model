package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type Schema struct {
	Name   string
	Tables []Table
}

type Table struct {
	Name              string
	JavaBeanName      string
	JavaBeanNameLower string
	Comment           string
	Columns           []Column
}

type Column struct {
	Name             string
	Comment          string
	ColumnType       string
	JavaFieldName    string
	JavaFieldNameUP  string
	JavaFieldType    string
	TableName        string
	InListQuery      int
	SingleQuery      int
	RangeQuery       int
	UpdateField      int
	WhereUpdateField int
}

type ColumnType struct {
	OriginType string
}

func ParseSql(txt string) Schema {
	var schema Schema
	schema.Name = "fund_cloud"
	tables := make([]Table, 1)
	lines := strings.Split(txt, "\r\n")
	for _, line := range lines {
		fmt.Println(line)
		var table Table
		if strings.Contains(line, " TABLE ") {
			table.Name = line[strings.Index(line, "TABLE")+6 : strings.Index(line, "(")]
			tables = append(tables, table)
		}
	}
	schema.Tables = tables
	return schema
}

func ParseSqlFromMySqlSchema(ip string, port int, username string, password string, databaseName string) *Schema {
	db, _ := sql.Open("mysql", username+":"+password+"@/"+databaseName+"?charset=utf8")
	rows, err := db.Query("SELECT TABLE_NAME,TABLE_COMMENT FROM information_schema.TABLES WHERE TABLE_SCHEMA = ?", databaseName)
	if err != nil {
		panic(err)
	}

	tableMap := make(map[string]*Table)
	for rows.Next() {
		var tableName, comment string
		err = rows.Scan(&tableName, &comment)
		if err != nil {
			panic(err)
		}
		table := Table{
			Name:              tableName,
			JavaBeanName:      convertJavaBeanName(tableName),
			JavaBeanNameLower: strings.ToLower(string(convertJavaBeanName(tableName)[0])) + convertJavaBeanName(tableName)[1:],
			Comment:           comment,
			Columns:           make([]Column, 0)}
		tableMap[tableName] = &table
	}

	rows, err = db.Query("SELECT TABLE_NAME,COLUMN_NAME,COLUMN_COMMENT,DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ?", databaseName)
	if err != nil {
		panic(err)
	}

	columns := make([]Column, 0)
	for rows.Next() {
		var tableName, name, comment, columnType string
		err = rows.Scan(&tableName, &name, &comment, &columnType)
		if err != nil {
			panic(err)
		}
		updateField := 1
		if name == "created_time" || name == "modified_time" {
			updateField = 0
		}
		column := Column{
			Name:            name,
			Comment:         comment,
			ColumnType:      columnType,
			JavaFieldName:   convertJavaFieldName(name),
			JavaFieldNameUP: convertJavaBeanName(name),
			JavaFieldType:   convertJavaFieldType(columnType),
			UpdateField:     updateField,
			TableName:       tableName}
		columns = append(columns, column)
	}
	for _, column := range columns {
		table := tableMap[column.TableName]
		if table == nil {
			panic("cannot find table " + column.TableName)
		}
		table.Columns = append(table.Columns, column)
	}
	tables := make([]Table, 0)
	for _, table := range tableMap {
		tables = append(tables, *table)
	}
	schema := Schema{
		Name:   databaseName,
		Tables: tables}

	return &schema
}

func convertJavaBeanName(tableName string) string {
	result := strings.ToUpper(string(tableName[0]))
	for i := 1; i < len(tableName); i++ {
		if tableName[i] == '_' {
			result += string(strings.ToUpper(string(tableName[i+1])))
			i++
		} else {
			result += string(tableName[i])
		}
	}
	return result
}

func convertJavaFieldName(tableName string) string {
	result := ""
	for i := 0; i < len(tableName); i++ {
		if tableName[i] == '_' {
			result += string(strings.ToUpper(string(tableName[i+1])))
			i++
		} else {
			result += string(tableName[i])
		}
	}
	return result
}

func convertJavaFieldType(t string) string {
	result := ""
	switch t {
	case "tinyint":
		result = "Integer"
	case "int":
		result = "Integer"
	case "bigint":
		result = "Long"
	case "varchar":
		result = "String"
	case "date":
		result = "Date"
	case "datetime":
		result = "Date"
	case "timestamp":
		result = "Date"
	case "decimal":
		result = "BigDecimal"
	}
	return result
}
