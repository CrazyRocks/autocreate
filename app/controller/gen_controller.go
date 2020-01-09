/**
 * @Author: Rocks
 * @Description:
 * @File:  gen_controoler
 * @Version: 1.0.0
 * @Date: 2019-10-13 21:58
 */

package controller

import (
	"autocreate/app/model"
	"autocreate/library/mlog"
	"autocreate/utils/base"
	"bytes"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
	"github.com/olekukonko/tablewriter"
	"strings"
)

/**
生成页面
*/
func Gen(r *ghttp.Request) {

	err := r.Response.WriteTpl("generator.html", g.Map{})
	if err != nil {
		glog.Error(err)
	}

}

func List(r *ghttp.Request) {
	tableName := r.GetString("tableName")
	var tables []model.Table
	if tableName != "" {
		g.DB("default").Table("information_schema.tables").Fields("table_name,engine,table_comment,create_time").Where("table_schema = (select database()) ").And("table_name like ?", "%"+tableName+"%").OrderBy("create_time desc").Structs(&tables)
	} else {
		g.DB("default").Table("information_schema.tables").Fields("table_name,engine,table_comment,create_time").Where("table_schema = (select database())").OrderBy("create_time desc").Structs(&tables)
	}

	base.Succ(r, g.Map{
		"list":       tables,
		"currPage":   1,
		"totalPage":  1,
		"totalCount": len(tables),
	})

}

/**
生成执行
*/
func Code(r *ghttp.Request) {

	tables := r.GetString("tables")
	module := r.GetString("module")
	project := r.GetString("project")
	moduleName := r.GetString("moduleName")
	if tables == "" || module == "" || moduleName == "" {
		return
	}
	tableNames := strings.Split(tables, ",")
	GenRouter(tableNames, module, project)
	GenModule(module, project)
	for _, table := range tableNames {
		GenCode(project, table, module, moduleName)
	}
	err := r.Response.WriteTpl("success.html", g.Map{})
	if err != nil {
		glog.Error(err)
	}

	tErr := r.Response.WriteTpl("success.html", g.Map{})
	if tErr != nil {
		glog.Error(err)
	}
}
func CamelCase(name string) string {
	return gstr.CamelCase(name)
}
func ScamelCase(name string) string {
	return gstr.CamelLowerCase(name)
}
func CameTable(name string) string {
	return strings.ToLower(gstr.Replace(name, "_", "/"))
}

func SubTable(name string) string {
	tables := gstr.Split(name, "_")
	println(tables)
	return tables[1]
}

func GenModule(module string, project string) {
	v := g.View()
	c := g.Config()
	folderPath := "result/" + module
	moduleContent, _ := v.Parse("module.go.html", map[string]interface{}{
		"module":   module,
		"project":  project,
		"DateTime": gtime.Date(),
		"Author":   c.GetString("author"),
		"Email":    c.GetString("email"),
	})
	path := folderPath + gfile.Separator + "module.go"
	if err := gfile.PutContents(path, strings.TrimSpace(moduleContent)); err != nil {
		mlog.Fatalf("writing model content to '%s' failed: %v", path, err)
	}
}
func GenRouter(tables []string, module string, project string) {
	v := g.View()
	v.BindFunc("CamelCase", CamelCase)
	v.BindFunc("ScamelCase", ScamelCase)
	v.BindFunc("CameTable", CameTable)
	v.BindFunc("SubTable", SubTable)
	c := g.Config()
	folderPath := "result/" + module + "/config"
	routerContent, _ := v.Parse("router.go.html", map[string]interface{}{
		"module":   module,
		"tables":   tables,
		"project":  project,
		"DateTime": gtime.Date(),
		"Author":   c.GetString("author"),
		"Email":    c.GetString("email"),
	})
	path := folderPath + gfile.Separator + "router.go"
	if err := gfile.PutContents(path, strings.TrimSpace(routerContent)); err != nil {
		mlog.Fatalf("writing model content to '%s' failed: %v", path, err)
	}
}

func GenCode(project string, table string, module string, moduleName string) {
	pathName := strings.ToLower(gstr.CamelCase(table))
	GenMenu(table, module, moduleName, pathName)
	GenVue(table, module)
	GenHtml(table, module, moduleName)
	GenGo(project, table, module, moduleName)
}

func GenMenu(table string, module string, moduleName string, pathName string) {
	v := g.View()
	menuContent, err := v.Parse("menu.sql.html", map[string]interface{}{
		"pathName":   pathName,
		"module":     module,
		"moduleName": moduleName,
	})
	if err != nil {
		mlog.Debug(err.Error())
	}
	folderPath := "result/sql"
	path := folderPath + gfile.Separator + table + "_menu.sql"

	if err := gfile.PutContents(path, strings.TrimSpace(menuContent)); err != nil {
		mlog.Fatalf("writing model content to '%s' failed: %v", path, err)
	}
}
func GenHtml(table string, module string, moduleName string) {

	tablePath := strings.ToLower(gstr.Replace(table, "_", "/"))

	var tableModel model.Table
	err := g.DB("default").Table("information_schema.tables").Fields("table_name TableName, engine, table_comment TableComment, create_time CreateTime").Where("table_schema = (select database())").And("table_name =?", table).Struct(&tableModel)
	if err != nil {
		mlog.Fatalf("查询错误%s", err.Error())
	}
	var columns []model.Column
	err = g.DB("default").Table("information_schema.columns").Fields("column_name ColumnName, data_type DataType, column_comment Comments, column_key ColumnKey, extra Extra").Where("table_schema = (select database())").And("table_name =?", table).OrderBy("ordinal_position").Structs(&columns)
	if err != nil {
		mlog.Fatalf("查询错误%s", err.Error())
	}
	v := g.View()
	c := g.Config()

	var PK model.Column
	for i := 0; i < len(columns); i++ {

		columns[i].AttrName = gstr.CamelCase(columns[i].ColumnName)
		if columns[i].ColumnKey == "PRI" {
			PK = columns[i]
		}
	}

	folderPath := "result/html"

	indexContent, err := v.Parse("html/list.html", map[string]interface{}{
		"moduleName": module,
		"tableName":  tableModel.TableName,
		"comments":   tableModel.TableComment,
		"structName": gstr.CamelCase(table),
		"classname":  strings.ToLower(gstr.CamelLowerCase(table)),
		"pathName":   tablePath,
		"pk":         PK,
		"title":      tableModel.TableComment,
		"columns":    columns,
		"nowTime":    gtime.Now().Unix(),
		"DateTime":   gtime.Date(),
		"Author":     c.GetString("author"),
		"Email":      c.GetString("email"),
	})

	if err != nil {
		mlog.Fatalf("解析错误%s", err.Error())
	}

	path := folderPath + gfile.Separator + tablePath + ".html"
	if err := gfile.PutContents(path, strings.TrimSpace(indexContent)); err != nil {
		mlog.Fatalf("writing model content to '%s' failed: %v", path, err)
	}

	mlog.Print(PK)

	editContent, err := v.Parse("html/list.js.html", map[string]interface{}{
		"moduleName": module,
		"tableName":  tableModel.TableName,
		"comments":   tableModel.TableComment,
		"structName": gstr.CamelCase(table),
		"pathName":   tablePath,
		"pk":         PK,
		"columns":    columns,
		"classname":  strings.ToLower(gstr.CamelLowerCase(table)),
		"DateTime":   gtime.Date(),
		"nowTime":    gtime.Now().Unix(),
		"Author":     c.GetString("author"),
		"Email":      c.GetString("email"),
	})
	if err != nil {
		mlog.Fatalf("解析错误%s", err.Error())
	}
	jsPath := "result/js"
	editPath := jsPath + gfile.Separator + tablePath + ".js"
	if err := gfile.PutContents(editPath, strings.TrimSpace(editContent)); err != nil {
		mlog.Fatalf("writing model content to '%s' failed: %v", path, err)
	}

}
func GenVue(table string, module string) {

	var tableModel model.Table
	err := g.DB("default").Table("information_schema.tables").Fields("table_name TableName, engine, table_comment TableComment, create_time CreateTime").Where("table_schema = (select database())").And("table_name =?", table).Struct(&tableModel)
	if err != nil {
		mlog.Fatalf("查询错误%s", err.Error())
	}
	var columns []model.Column
	err = g.DB("default").Table("information_schema.columns").Fields("column_name ColumnName, data_type DataType, column_comment Comments, column_key ColumnKey, extra Extra").Where("table_schema = (select database())").And("table_name =?", table).OrderBy("ordinal_position").Structs(&columns)
	if err != nil {
		mlog.Fatalf("查询错误%s", err.Error())
	}
	v := g.View()
	c := g.Config()

	var PK model.Column
	for i := 0; i < len(columns); i++ {

		columns[i].AttrName = gstr.CamelCase(columns[i].ColumnName)
		if columns[i].ColumnKey == "PRI" {
			PK = columns[i]
		}
	}

	folderPath := "result/vue/" + module

	indexContent, err := v.Parse("vue/index.vue.html", map[string]interface{}{
		"moduleName": module,
		"tableName":  tableModel.TableName,
		"comments":   tableModel.TableComment,
		"structName": gstr.CamelCase(table),
		"pathName":   strings.ToLower(gstr.CamelLowerCase(table)),
		"pk":         PK,
		"columns":    columns,
		"DateTime":   gtime.Date(),
		"Author":     c.GetString("author"),
		"Email":      c.GetString("email"),
	})

	if err != nil {
		mlog.Fatalf("解析错误%s", err.Error())
	}
	path := folderPath + gfile.Separator + strings.ToLower(gstr.CamelLowerCase(table)) + ".vue"
	if err := gfile.PutContents(path, strings.TrimSpace(indexContent)); err != nil {
		mlog.Fatalf("writing model content to '%s' failed: %v", path, err)
	}

	mlog.Print(PK)

	editContent, err := v.Parse("vue/add-or-update.vue.html", map[string]interface{}{
		"moduleName": module,
		"tableName":  tableModel.TableName,
		"comments":   tableModel.TableComment,
		"structName": gstr.CamelCase(table),
		"pathName":   strings.ToLower(gstr.CamelLowerCase(table)),
		"pk":         PK,
		"columns":    columns,
		"classname":  strings.ToLower(gstr.CamelLowerCase(table)),
		"DateTime":   gtime.Date(),
		"Author":     c.GetString("author"),
		"Email":      c.GetString("email"),
	})
	if err != nil {
		mlog.Fatalf("解析错误%s", err.Error())
	}
	editPath := folderPath + gfile.Separator + strings.ToLower(gstr.CamelLowerCase(table)) + "-add-or-update.vue"
	if err := gfile.PutContents(editPath, strings.TrimSpace(editContent)); err != nil {
		mlog.Fatalf("writing model content to '%s' failed: %v", path, err)
	}

}

func GenGo(project string, table string, module string, moduleName string) {
	db := g.DB("default")
	folderPath := "result/" + module
	generateModelContentFile(project, db, table, folderPath, table, "default")
	generateControllerContentFile(project, db, table, folderPath, table, module, "default")
}
func generateControllerContentFile(project string, db gdb.DB, table string, folderPath, packageName, module, groupName string) {
	camelName := gstr.CamelCase(table)
	v := g.View()
	v.BindFunc("CameTable", CameTable)
	c := g.Config()
	tablePath := strings.ToLower(gstr.Replace(table, "_", "/"))
	controllerContent, err := v.Parse("controller.go.html", map[string]interface{}{
		"TplTableName":   table,
		"TplProject":     project,
		"TplModelName":   camelName,
		"TplGroupName":   groupName,
		"TplPackageName": packageName,
		"pathName":       tablePath,
		"moduleName":     module,
		"table":          table,
		"DateTime":       gtime.Date(),
		"Author":         c.GetString("author"),
		"Email":          c.GetString("email"),
	})
	if err != nil {
		print(err.Error())
	}
	name := gstr.Trim(gstr.SnakeCase(table), "-_.")
	path := folderPath + gfile.Separator + "controller" + gfile.Separator + name + "_controller.go"
	if err := gfile.PutContents(path, strings.TrimSpace(controllerContent)); err != nil {
		mlog.Fatalf("writing controller content to '%s' failed: %v", path, err)
	}
}
func generateModelContentFile(project string, db gdb.DB, table string, folderPath, packageName, groupName string) {

	fields, err := db.TableFields(table)
	if err != nil {
		mlog.Fatalf("fetching tables fields failed for table '%s':\n%v", table, err)
	}
	camelName := gstr.CamelCase(table)
	structDefine := generateStructDefinition(table, fields)
	columnDefine := ""
	index := 1
	for _, field := range fields {
		if index == len(fields) {
			columnDefine = columnDefine + "t." + field.Name + " as " + gstr.CamelCase(field.Name)
		} else {
			columnDefine = columnDefine + "t." + field.Name + " as " + gstr.CamelCase(field.Name) + ","
		}
		index = index + 1
	}
	extraImports := ""
	if strings.Contains(structDefine, "gtime.Time") {
		extraImports = `
import (
	"github.com/gogf/gf/os/gtime"
)
`
	}
	v := g.View()
	c := g.Config()

	modelContent, err := v.Parse("model.go.html", map[string]interface{}{
		"TplTableName":    table,
		"TplProject":      project,
		"TplModelName":    camelName,
		"TplGroupName":    groupName,
		"columnDefine":    columnDefine,
		"TplPackageName":  packageName,
		"TplExtraImports": extraImports,
		"TplStructDefine": structDefine,
		"DateTime":        gtime.Date(),
		"Author":          c.GetString("author"),
		"Email":           c.GetString("email"),
	})
	name := gstr.Trim(gstr.SnakeCase(table), "-_.")
	if len(name) > 5 && name[len(name)-5:] == "_test" {
		// Add suffix to avoid the table name which contains "_test",
		// which would make the go file a testing file.
		name += "_table"
	}
	path := folderPath + gfile.Separator + "model" + gfile.Separator + name + "_model.go"
	if err := gfile.PutContents(path, strings.TrimSpace(modelContent)); err != nil {
		mlog.Fatalf("writing model content to '%s' failed: %v", path, err)
	}
}
func generateStructDefinition(table string, fields map[string]*gdb.TableField) string {
	buffer := bytes.NewBuffer(nil)
	array := make([][]string, len(fields))
	for _, field := range fields {
		array[field.Index] = generateStructField(field)
	}
	tw := tablewriter.NewWriter(buffer)
	tw.SetBorder(false)
	tw.SetRowLine(false)
	tw.SetAutoWrapText(false)
	tw.SetColumnSeparator("")
	tw.AppendBulk(array)
	tw.Render()
	stContent := buffer.String()
	// Let's do this hack for tablewriter!
	stContent = gstr.Replace(stContent, "  #", "")
	buffer.Reset()
	buffer.WriteString("type " + gstr.CamelCase(table) + " struct {\n")
	buffer.WriteString(stContent)
	buffer.WriteString("}")
	return buffer.String()
}
func generateStructField(field *gdb.TableField) []string {
	var typeName, ormTag, jsonTag string
	t, _ := gregex.ReplaceString(`\(.+\)`, "", field.Type)
	arr := gstr.Split(t, " ")
	t = strings.ToLower(arr[0])
	t = strings.ToLower(t)
	fmt.Printf("当前Type:%s", field.Type)
	fmt.Printf("当前字段:%s:%s\n", field.Name, t)
	switch t {
	case "binary", "varbinary", "blob", "tinyblob", "mediumblob", "longblob":
		typeName = "[]byte"

	case "bit", "int", "tinyint", "smallint", "mediumint":
		if gstr.ContainsI(field.Type, "unsigned") {
			typeName = "uint"
		} else {
			typeName = "int"
		}

	case "bigint":
		if gstr.ContainsI(field.Type, "unsigned") {
			typeName = "uint64"
		} else {
			typeName = "int64"
		}

	case "float", "double", "decimal":
		typeName = "float64"

	case "bool":
		typeName = "bool"

	case "datetime", "timestamp", "date", "time":
		typeName = "*gtime.Time"

	default:
		// Auto detecting type.
		switch {
		case strings.Contains(t, "int64"):
			typeName = "int64"
		case strings.Contains(t, "int"):
			typeName = "int"
		case strings.Contains(t, "text") || strings.Contains(t, "char"):
			typeName = "string"
		case strings.Contains(t, "float") || strings.Contains(t, "double"):
			typeName = "float64"
		case strings.Contains(t, "bool"):
			typeName = "bool"
		case strings.Contains(t, "binary") || strings.Contains(t, "blob"):
			typeName = "[]byte"
		case strings.Contains(t, "date") || strings.Contains(t, "time"):
			typeName = "*gtime.Time"
		default:
			typeName = "string"
		}
	}
	ormTag = field.Name
	jsonTag = gstr.CamelCase(field.Name)
	if gstr.ContainsI(field.Key, "pri") {
		ormTag += ",primary"
	}
	if gstr.ContainsI(field.Key, "uni") {
		ormTag += ",unique"
	}
	return []string{
		"    #" + gstr.CamelCase(field.Name),
		" #" + typeName,
		" #" + fmt.Sprintf("`"+`orm:"%s"`, ormTag),
		" #" + fmt.Sprintf(`json:"%s,omitempty" gconv:"%s,omitempty"`+"`", gstr.CamelCase(field.Name), jsonTag),
	}
}
