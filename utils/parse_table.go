package utils

import (
	"fmt"
	"html/template"
	"strings"
)

// ParseTable 解析表字段
func (config *GlobalConfig) ParseTable(fields *[]TableField) {
	if len(*fields) != 0 {
		isParse := true
		Struct := &Struct{
			Fields: func() *[]StructField {
				list := make([]StructField, len(*fields))
				return &list
			}(),
		}
		for key, field := range *fields {
			if isParse {
				Struct.Table = config.Database.GetTable(field.TableName)

				Struct.StructName = snake(Struct.Table.Name)
				if Struct.Table.Comment == "" {
					Struct.Table.Comment = "."
				}
				isParse = false
			}
			(*Struct.Fields)[key] = config.parseField(&field)
		}
		config.Structs = append(config.Structs, Struct)
	}
}

// snake aaa_bbb to AaaBbb
func snake(source string) string {
	builder := strings.Builder{}
	isHandle := false
	for _, word := range strings.Split(source, "_") {
		isHandle = true
		if word != "id" {
			for _, letter := range []rune(word) {
				if isHandle {
					builder.WriteString(strings.ToUpper(string(letter)))
					isHandle = false
				} else {
					builder.WriteString(string(letter))
				}
			}
		} else {
			builder.WriteString(strings.ToUpper(word))
		}
	}
	return builder.String()
}

// ensnake AaaBbb to aaa_bbb
func ensnake(source string) string {
	builder := strings.Builder{}
	for key, word := range []rune(source) {
		if word <= 90 {
			if key != 0 {
				builder.WriteString("_")
			}
			builder.WriteString(strings.ToLower(string(word)))
		} else {
			builder.WriteString(string(word))
		}
	}
	return builder.String()
}

// camel AaaBbb to aaaBbb
func camel(source string) string {
	builder := strings.Builder{}
	if source != "ID" {
		isHandle := true
		for _, letter := range []rune(source) {
			if isHandle {
				builder.WriteString(strings.ToLower(string(letter)))
				isHandle = false
			} else {
				builder.WriteString(string(letter))
			}
		}
	} else {
		builder.WriteString(strings.ToLower(source))
	}
	return builder.String()
}

// parseField 解析字段
func (config *GlobalConfig) parseField(tf *TableField) StructField {
	return StructField{
		Name:    snake(tf.Name),
		Type:    config.parseType(tf),
		Tag:     config.parseTag(tf),
		Comment: tf.Comment,
	}
}

// parseType 解析字段类型
func (config *GlobalConfig) parseType(tf *TableField) string {
	builder := strings.Builder{}
	// if strings.HasSuffix(tf.Column, "unsigned") {
	// 	builder.WriteString("u")
	// }
	gotype := mysqltype[tf.Type]
	if strings.Contains(gotype, "time.Time") {
		config.Imports[`"time"`] = ""
	}
	builder.WriteString(gotype)
	return builder.String()
}

// parseTag 解析字段 tag. 包含 json tag 和 orm tag
func (config *GlobalConfig) parseTag(tf *TableField) template.HTML {
	builder := strings.Builder{}
	builder.WriteString("`")

	// json tag
	{
		// aaBbCc
		// builder.WriteString(fmt.Sprintf(`json:"%s"`, camel(snake(tf.Name))))
		// aa_bb_cc
		builder.WriteString(fmt.Sprintf(`json:"%s"`, tf.Name))
	}

	builder.WriteString(" ")

	// orm tag
	{
		builder.WriteString(fmt.Sprintf(`gorm:"column:%s;type:%s`, tf.Name, tf.Column))
		if tf.IsNull == "NO" {
			builder.WriteString(";not null")
		}
		if tf.Primary == "PRI" {
			builder.WriteString(";primary_key")
		}
		if tf.Extra == "auto_increment" {
			builder.WriteString(";auto_increment")
		}
		builder.WriteString(`"`)
	}

	builder.WriteString("`")
	return template.HTML(builder.String())
}
