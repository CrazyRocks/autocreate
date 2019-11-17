/**
 * @Author: Rocks
 * @Description: 
 * @File:  table
 * @Version: 1.0.0
 * @Date: 2019-10-13 22:34
 */

package model

type Table struct {
	// columns START
	TableName    string `json:"tableName"`
	Engine       string `json:"engine"`
	TableComment string `json:"tableComment"`
	CreateTime   string `json:"createTime"`
	Pk string
}


