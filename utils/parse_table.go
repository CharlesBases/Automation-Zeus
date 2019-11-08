package utils

import (
	"fmt"
	"html/template"
	"strings"
)

func (config *GlobalConfig) ParseTable(fields *[]TableField) {
	if len(*fields) != 0 {
		isParse := true
		Struct := Struct{
			Fields: func() *[]StructField {
				list := make([]StructField, len(*fields))
				return &list
			}(),
		}
		for key, field := range *fields {
			if isParse {
				Struct.StructName = snake(field.TableName)
				isParse = false
			}
			(*Struct.Fields)[key] = config.parseField(&field)
		}
		*config.Structs = append(*config.Structs, Struct)
	}
}

// aaa_bbb to AaaBbb
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

// AaaBbb to aaa_bbb
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

// AaaBbb to aaaBbb
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

func (config *GlobalConfig) parseField(tf *TableField) StructField {
	return StructField{
		Name:    snake(tf.Name),
		Type:    config.parseType(tf),
		Tag:     config.parseTag(tf),
		Comment: tf.Comment,
	}
}

func (config *GlobalConfig) parseType(tf *TableField) string {
	builder := strings.Builder{}
	// if strings.HasSuffix(tf.Column, "unsigned") {
	// 	builder.WriteString("u")
	// }
	gotype := mysqltype[tf.Type]
	if strings.Contains(gotype, "time.Time") {
		config.Import[`"time"`] = ""
	}
	builder.WriteString(gotype)
	return builder.String()
}

func (config *GlobalConfig) parseTag(tf *TableField) template.HTML {
	builder := strings.Builder{}
	if config.Json || config.Gorm {
		builder.WriteString("`")
		if config.Json {
			builder.WriteString(fmt.Sprintf(`json:"%s"`, camel(snake(tf.Name))))
		}
		if config.Gorm {
			if config.Json {
				builder.WriteString(" ")
			}
			builder.WriteString(fmt.Sprintf(`gorm:"column:%s;type:%s`, tf.Name, tf.Column))
			if tf.IsNull == "NO" {
				builder.WriteString(";not null")
			}
			if tf.Primary == "PRI" {
				builder.WriteString(";primary_key")
			}
			if len(tf.Extra) != 0 {
				builder.WriteString(";auto_increment")
			}
			builder.WriteString(`"`)
		}
		builder.WriteString("`")
	}
	return template.HTML(builder.String())
}
