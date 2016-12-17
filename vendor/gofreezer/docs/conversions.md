# conversions

## mysql

### null column

如果数据库中存在`null`字段, 提供以下几种方法来保证`storage`正确的解析。

- 声明字段类型为： `sql.NullString`, `sql.NullBool`,`sql.NullInt64`,`sql.NullFloat64`
- 声明字段类型为指针类型
- 设置数据库字段的`NOT NULL`属性
