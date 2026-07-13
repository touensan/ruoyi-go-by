package toolcontroller

import (
	"archive/zip"
	"bytes"
	"net/http"
	"os"
	"ruoyi-go/config"
	"ruoyi-go/framework/dal"
	"ruoyi-go/framework/response"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type GenController struct{}

type genTable struct {
	TableId        int         `json:"tableId"`
	TableName      string      `json:"tableName"`
	TableComment   string      `json:"tableComment"`
	ClassName      string      `json:"className"`
	TplCategory    string      `json:"tplCategory"`
	TplWebType     string      `json:"tplWebType"`
	PackageName    string      `json:"packageName"`
	ModuleName     string      `json:"moduleName"`
	BusinessName   string      `json:"businessName"`
	FunctionName   string      `json:"functionName"`
	FunctionAuthor string      `json:"functionAuthor"`
	GenType        string      `json:"genType"`
	GenPath        string      `json:"genPath"`
	CreateTime     string      `json:"createTime"`
	UpdateTime     string      `json:"updateTime"`
	Columns        []genColumn `json:"columns,omitempty" gorm:"-"`
}

type genColumn struct {
	ColumnId      int    `json:"columnId"`
	TableId       int    `json:"tableId"`
	ColumnName    string `json:"columnName"`
	ColumnComment string `json:"columnComment"`
	ColumnType    string `json:"columnType"`
	JavaType      string `json:"javaType"`
	JavaField     string `json:"javaField"`
	IsPk          string `json:"isPk"`
	IsIncrement   string `json:"isIncrement"`
	IsRequired    string `json:"isRequired"`
	IsInsert      string `json:"isInsert"`
	IsEdit        string `json:"isEdit"`
	IsList        string `json:"isList"`
	IsQuery       string `json:"isQuery"`
	QueryType     string `json:"queryType"`
	HtmlType      string `json:"htmlType"`
	DictType      string `json:"dictType"`
	Sort          int    `json:"sort"`
}

func (*GenController) List(ctx *gin.Context) {
	tables, total := listDatabaseTables(ctx)
	response.NewSuccess().SetPageData(tables, total).Json(ctx)
}

func (*GenController) DbList(ctx *gin.Context) {
	tables, total := listDatabaseTables(ctx)
	response.NewSuccess().SetPageData(tables, total).Json(ctx)
}

func (*GenController) Detail(ctx *gin.Context) {
	tableId, _ := strconv.Atoi(ctx.Param("tableId"))
	tables, _ := listDatabaseTables(ctx)
	if len(tables) == 0 {
		response.NewSuccess().SetData("data", gin.H{
			"info":   genTable{TableId: tableId, TplCategory: "crud", TplWebType: "element-plus", GenType: "0"},
			"rows":   []genColumn{},
			"tables": []genTable{},
		}).Json(ctx)
		return
	}
	index := tableId - 1
	if index < 0 || index >= len(tables) {
		index = 0
	}
	info := tables[index]
	info.Columns = tableColumns(info.TableId, info.TableName)
	response.NewSuccess().SetData("data", gin.H{
		"info":   info,
		"rows":   info.Columns,
		"tables": tables,
	}).Json(ctx)
}

func (*GenController) Update(ctx *gin.Context) {
	response.NewSuccess().Json(ctx)
}

func (*GenController) ImportTable(ctx *gin.Context) {
	response.NewSuccess().Json(ctx)
}

func (*GenController) CreateTable(ctx *gin.Context) {
	response.NewSuccess().Json(ctx)
}

func (*GenController) Preview(ctx *gin.Context) {
	response.NewSuccess().SetData("data", gin.H{
		"domain.go.vm":     "package model\n\n// 当前接口为 RuoYi-Go BY 兼容预览，占位展示数据库结构。\n",
		"controller.go.vm": "package controller\n\n// 当前 ruoyi-go 基线未启用完整代码生成器。\n",
		"vue/index.vue.vm": "<template><div>RuoYi-Go BY 代码生成预览</div></template>\n",
	}).Json(ctx)
}

func (*GenController) Remove(ctx *gin.Context) {
	response.NewSuccess().Json(ctx)
}

func (*GenController) GenCode(ctx *gin.Context) {
	response.NewSuccess().SetMsg("当前 Go 基线未启用完整代码生成器，已接受生成请求").Json(ctx)
}

func (*GenController) SynchDb(ctx *gin.Context) {
	response.NewSuccess().Json(ctx)
}

func (*GenController) BatchGenCode(ctx *gin.Context) {
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	file, _ := writer.Create("README.txt")
	_, _ = file.Write([]byte("RuoYi-Go BY 当前后端已提供代码生成下载占位接口。\n"))
	_ = writer.Close()
	ctx.Header("Content-Type", "application/zip")
	ctx.Header("Content-Disposition", "attachment; filename=ruoyi.zip")
	ctx.Data(http.StatusOK, "application/zip", buffer.Bytes())
}

func listDatabaseTables(ctx *gin.Context) ([]genTable, int) {
	type tableRow struct {
		TableName    string
		TableComment string
		CreateTime   *time.Time
		UpdateTime   *time.Time
	}
	rows := make([]tableRow, 0)
	query := dal.Gorm.Table("information_schema.tables").
		Select("table_name, table_comment, create_time, update_time").
		Where("table_schema = ?", config.Data.Mysql.Database).
		Order("table_name")
	if tableName := ctx.Query("tableName"); tableName != "" {
		query = query.Where("table_name LIKE ?", "%"+tableName+"%")
	}
	if tableComment := ctx.Query("tableComment"); tableComment != "" {
		query = query.Where("table_comment LIKE ?", "%"+tableComment+"%")
	}
	query.Scan(&rows)

	tables := make([]genTable, 0, len(rows))
	for index, row := range rows {
		createTime, updateTime := "", ""
		if row.CreateTime != nil {
			createTime = row.CreateTime.Format("2006-01-02 15:04:05")
		}
		if row.UpdateTime != nil {
			updateTime = row.UpdateTime.Format("2006-01-02 15:04:05")
		}
		tables = append(tables, genTable{
			TableId:        index + 1,
			TableName:      row.TableName,
			TableComment:   row.TableComment,
			ClassName:      toCamel(row.TableName),
			TplCategory:    "crud",
			TplWebType:     "element-plus",
			PackageName:    "ruoyi-go/app",
			ModuleName:     moduleName(row.TableName),
			BusinessName:   row.TableName,
			FunctionName:   row.TableComment,
			FunctionAuthor: "RuoYi-Go BY",
			GenType:        "0",
			GenPath:        defaultGenPath(),
			CreateTime:     createTime,
			UpdateTime:     updateTime,
		})
	}
	return tables, len(tables)
}

func defaultGenPath() string {
	wd, err := os.Getwd()
	if err != nil || strings.TrimSpace(wd) == "" {
		return "."
	}
	return wd
}

func tableColumns(tableId int, tableName string) []genColumn {
	type columnRow struct {
		ColumnName      string
		ColumnComment   string
		ColumnType      string
		ColumnKey       string
		Extra           string
		IsNullable      string
		OrdinalPosition int
	}
	rows := make([]columnRow, 0)
	dal.Gorm.Table("information_schema.columns").
		Select("column_name, column_comment, column_type, column_key, extra, is_nullable, ordinal_position").
		Where("table_schema = ? AND table_name = ?", config.Data.Mysql.Database, tableName).
		Order("ordinal_position").
		Scan(&rows)

	columns := make([]genColumn, 0, len(rows))
	for index, row := range rows {
		columns = append(columns, genColumn{
			ColumnId:      index + 1,
			TableId:       tableId,
			ColumnName:    row.ColumnName,
			ColumnComment: row.ColumnComment,
			ColumnType:    row.ColumnType,
			JavaType:      goType(row.ColumnType),
			JavaField:     lowerFirst(toCamel(row.ColumnName)),
			IsPk:          boolFlag(row.ColumnKey == "PRI"),
			IsIncrement:   boolFlag(strings.Contains(row.Extra, "auto_increment")),
			IsRequired:    boolFlag(row.IsNullable == "NO"),
			IsInsert:      "1",
			IsEdit:        boolFlag(row.ColumnKey != "PRI"),
			IsList:        "1",
			IsQuery:       "0",
			QueryType:     "EQ",
			HtmlType:      htmlType(row.ColumnType),
			Sort:          row.OrdinalPosition,
		})
	}
	return columns
}

func toCamel(value string) string {
	parts := strings.Split(value, "_")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, "")
}

func lowerFirst(value string) string {
	if value == "" {
		return value
	}
	return strings.ToLower(value[:1]) + value[1:]
}

func moduleName(tableName string) string {
	parts := strings.Split(tableName, "_")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	return "system"
}

func goType(columnType string) string {
	t := strings.ToLower(columnType)
	switch {
	case strings.Contains(t, "bigint"), strings.Contains(t, "int"):
		return "int"
	case strings.Contains(t, "decimal"), strings.Contains(t, "double"), strings.Contains(t, "float"):
		return "float64"
	case strings.Contains(t, "datetime"), strings.Contains(t, "timestamp"), strings.Contains(t, "date"):
		return "time.Time"
	default:
		return "string"
	}
}

func htmlType(columnType string) string {
	t := strings.ToLower(columnType)
	if strings.Contains(t, "text") {
		return "textarea"
	}
	if strings.Contains(t, "datetime") || strings.Contains(t, "date") {
		return "datetime"
	}
	return "input"
}

func boolFlag(ok bool) string {
	if ok {
		return "1"
	}
	return "0"
}
