/**
 * @Author: Rocks
 * @Description: 
 * @File:  column
 * @Version: 1.0.0
 * @Date: 2019-10-14 00:20
 */

package model

type Column struct {
	// columns START
	//列名
	ColumnName string;
	//列名类型
	DataType string;
	//列名备注
	Comments string;
	//属性名称(第一个字母大写)，如：user_name => UserName
	AttrName string;
	//属性名称(第一个字母小写)，如：user_name => userName
	Attrname string;
	//属性类型
	AttrType string;
	//auto_increment
	Extra string;

	ColumnKey string;
}
