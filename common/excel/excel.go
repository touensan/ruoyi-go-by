package excel

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

const tagName = "excel"

type fieldColumn struct {
	fieldIndex int
	order      int
	index      int
	name       string
	width      float64
	replace    map[string]string
	reverse    map[string]string
}

func NormalDynamicExport(sheet, title, fields string, isGhbj, isIgnore bool, list interface{}, changeHead map[string]string) (*excelize.File, error) {
	if sheet == "" {
		sheet = "Sheet1"
	}

	value := reflect.ValueOf(list)
	if value.Kind() != reflect.Slice {
		return nil, errors.New("无效的数据类型")
	}

	elemType := value.Type().Elem()
	for elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}
	if elemType.Kind() != reflect.Struct {
		return nil, errors.New("无效的数据类型")
	}

	columns, err := exportColumns(elemType, fields, isIgnore, changeHead)
	if err != nil {
		return nil, err
	}

	file := excelize.NewFile()
	if sheet != "Sheet1" {
		if _, err = file.NewSheet(sheet); err != nil {
			return nil, err
		}
		_ = file.DeleteSheet("Sheet1")
	}
	if _, err = file.GetSheetIndex(sheet); err != nil {
		return nil, err
	}

	headerRow := 1
	dataRow := 2
	if title != "" {
		headerRow = 2
		dataRow = 3
		endCell, _ := excelize.CoordinatesToCellName(max(1, len(columns)), 1)
		_ = file.SetCellValue(sheet, "A1", title)
		_ = file.MergeCell(sheet, "A1", endCell)
	}

	for i, column := range columns {
		cell, _ := excelize.CoordinatesToCellName(i+1, headerRow)
		_ = file.SetCellValue(sheet, cell, column.name)
		if column.width > 0 {
			col, _ := excelize.ColumnNumberToName(i + 1)
			_ = file.SetColWidth(sheet, col, col, column.width)
		}
	}

	for rowIndex := 0; rowIndex < value.Len(); rowIndex++ {
		item := indirect(value.Index(rowIndex))
		if !item.IsValid() || item.Kind() != reflect.Struct {
			continue
		}
		for colIndex, column := range columns {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, dataRow+rowIndex)
			cellValue := valueString(item.Field(column.fieldIndex))
			if replacement, ok := column.replace[cellValue]; ok {
				cellValue = replacement
			}
			_ = file.SetCellValue(sheet, cell, cellValue)
		}
	}

	return file, nil
}

func ImportExcel(file *excelize.File, dst interface{}, headIndex, startRow int) error {
	sheet := file.GetSheetName(0)
	if sheet == "" {
		return errors.New("工作表不存在")
	}
	return importSheet(file, dst, sheet, headIndex, startRow)
}

func DownLoadExcel(fileName string, res http.ResponseWriter, file *excelize.File) {
	res.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	res.Header().Set("Content-Disposition", "attachment; filename*=UTF-8''"+url.PathEscape(fileName)+".xlsx")
	res.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	if err := file.Write(res); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func exportColumns(elemType reflect.Type, fields string, isIgnore bool, changeHead map[string]string) ([]fieldColumn, error) {
	columns := make([]fieldColumn, 0)
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		if field.PkgPath != "" {
			continue
		}
		tagValue := field.Tag.Get(tagName)
		if tagValue == "" {
			continue
		}
		if fields != "" {
			selected := strings.Contains(fields, field.Name+",")
			if isIgnore && selected {
				continue
			}
			if !isIgnore && !selected {
				continue
			}
		}
		column, err := parseColumn(tagValue)
		if err != nil {
			return nil, err
		}
		if column.name == "" {
			column.name = field.Name
		}
		if changeHead != nil && changeHead[field.Name] != "" {
			column.name = changeHead[field.Name]
		}
		column.fieldIndex = i
		column.order = i
		columns = append(columns, column)
	}

	sort.SliceStable(columns, func(i, j int) bool {
		if columns[i].index >= 0 && columns[j].index >= 0 {
			return columns[i].index < columns[j].index
		}
		if columns[i].index >= 0 {
			return true
		}
		if columns[j].index >= 0 {
			return false
		}
		return columns[i].order < columns[j].order
	})

	return columns, nil
}

func importSheet(file *excelize.File, dst interface{}, sheet string, headIndex, startRow int) error {
	rows, err := file.GetRows(sheet)
	if err != nil {
		return errors.New(sheet + "工作表不存在")
	}

	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr || dstValue.Elem().Kind() != reflect.Slice {
		return errors.New("无效的数据类型")
	}

	sliceValue := dstValue.Elem()
	elemType := sliceValue.Type().Elem()
	for elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}
	if elemType.Kind() != reflect.Struct {
		return errors.New("无效的数据类型")
	}

	columns, err := importColumns(elemType)
	if err != nil {
		return err
	}

	headers := []string{}
	if headIndex >= 0 && headIndex < len(rows) {
		headers = rows[headIndex]
	}

	for rowIndex, row := range rows {
		if rowIndex < startRow {
			continue
		}
		item := reflect.New(elemType).Elem()
		for _, column := range columns {
			cellIndex := column.index
			if cellIndex < 0 {
				cellIndex = headerIndex(headers, column.name)
			}
			if cellIndex < 0 || cellIndex >= len(row) {
				continue
			}
			cellValue := row[cellIndex]
			if replacement, ok := column.reverse[cellValue]; ok {
				cellValue = replacement
			}
			setFieldValue(item.Field(column.fieldIndex), cellValue)
		}
		sliceValue.Set(reflect.Append(sliceValue, item))
	}

	return nil
}

func importColumns(elemType reflect.Type) ([]fieldColumn, error) {
	columns := make([]fieldColumn, 0)
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		if field.PkgPath != "" {
			continue
		}
		tagValue := field.Tag.Get(tagName)
		if tagValue == "" {
			continue
		}
		column, err := parseColumn(tagValue)
		if err != nil {
			return nil, err
		}
		if column.name == "" {
			column.name = field.Name
		}
		column.fieldIndex = i
		column.order = i
		columns = append(columns, column)
	}
	return columns, nil
}

func parseColumn(tagValue string) (fieldColumn, error) {
	column := fieldColumn{
		index:   -1,
		replace: make(map[string]string),
		reverse: make(map[string]string),
	}
	for _, part := range strings.Split(tagValue, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		key, value, ok := strings.Cut(part, ":")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		switch key {
		case "name", "title":
			column.name = value
		case "index":
			index, err := strconv.Atoi(value)
			if err != nil {
				return column, fmt.Errorf("无效的 Excel 列索引: %w", err)
			}
			column.index = index
		case "width":
			width, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return column, fmt.Errorf("无效的 Excel 列宽: %w", err)
			}
			column.width = width
		case "replace":
			for _, item := range strings.Split(value, ",") {
				raw, display, ok := strings.Cut(item, "_")
				if !ok {
					continue
				}
				column.replace[raw] = display
				column.reverse[display] = raw
			}
		}
	}
	return column, nil
}

func indirect(value reflect.Value) reflect.Value {
	for value.IsValid() && value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return reflect.Value{}
		}
		value = value.Elem()
	}
	return value
}

func valueString(value reflect.Value) string {
	value = indirect(value)
	if !value.IsValid() {
		return ""
	}
	if value.CanInterface() {
		if t, ok := value.Interface().(time.Time); ok {
			if t.IsZero() {
				return ""
			}
			return t.Format("2006-01-02 15:04:05")
		}
		if s, ok := value.Interface().(fmt.Stringer); ok {
			return s.String()
		}
	}
	return fmt.Sprint(value.Interface())
}

func setFieldValue(field reflect.Value, value string) {
	if !field.CanSet() {
		return
	}
	field = indirectAssignable(field)
	if !field.IsValid() || !field.CanSet() {
		return
	}
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		parsed, _ := strconv.ParseBool(value)
		field.SetBool(parsed)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, _ := strconv.ParseInt(value, 10, field.Type().Bits())
		field.SetInt(parsed)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsed, _ := strconv.ParseUint(value, 10, field.Type().Bits())
		field.SetUint(parsed)
	case reflect.Float32, reflect.Float64:
		parsed, _ := strconv.ParseFloat(value, field.Type().Bits())
		field.SetFloat(parsed)
	}
}

func indirectAssignable(value reflect.Value) reflect.Value {
	for value.IsValid() && value.Kind() == reflect.Ptr {
		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}
		value = value.Elem()
	}
	return value
}

func headerIndex(headers []string, name string) int {
	for index, header := range headers {
		if strings.TrimSpace(header) == name {
			return index
		}
	}
	return -1
}
