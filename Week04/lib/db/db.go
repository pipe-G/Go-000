package db

import (
	"fmt"
	"Week04/lib/db/database"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	goqu "gopkg.in/doug-martin/goqu.v5"
	_ "gopkg.in/doug-martin/goqu.v5/adapters/mysql"
)

type QueryOp string

const (
	EQ      QueryOp = "eq"
	NEQ     QueryOp = "neq"
	BETWEEN QueryOp = "between"
	IN      QueryOp = "in"
	NOTIN   QueryOp = "notIn"
	GT      QueryOp = "gt"
	GTE     QueryOp = "gte"
	LT      QueryOp = "lt"
	LTE     QueryOp = "lte"
	LIKE    QueryOp = "like"
)

type DataBase struct {
	*gorm.DB
	Goqu *goqu.Database
}

var DB *DataBase
var Goqu *goqu.Database

func Open(debug bool, dialect, dburl string) error {
	err := database.Open(debug, dialect, dburl)
	if err != nil {
		return err
	}

	DB = &DataBase{
		DB:   database.DB(),
		Goqu: goqu.New(dialect, database.DB().DB()),
	}
	Goqu = DB.Goqu

	database.DB().DB().SetMaxIdleConns(10)
	database.DB().DB().SetMaxOpenConns(128)

	return nil
}

type Condition map[string]interface{}

// SqlArgsFilter ...
type SqlArgsFilter struct {
	goqu.Ex
	filterMap map[string]interface{}
	exOr      goqu.ExOr
}

// QueryFilter ...
func QueryFilter(filterMap Condition) *SqlArgsFilter {
	return &SqlArgsFilter{
		Ex:        goqu.Ex{},
		filterMap: filterMap,
		exOr:      make(goqu.ExOr),
	}
}

// Where ...
func (f *SqlArgsFilter) Where(field string, op QueryOp,
	filterField string) *SqlArgsFilter {
	if f.filterMap == nil {
		return f
	}

	value, ok := f.filterMap[filterField]
	if ok {
		if opEx := f.execOp(op, value); opEx != nil {
			f.Ex[field] = opEx
		}

	}
	return f
}

// Or ...
func (f *SqlArgsFilter) Or(field string, op QueryOp, filterField string) *SqlArgsFilter {
	if f.filterMap == nil {
		return f
	}

	value, ok := f.filterMap[filterField]
	if ok {
		if opEx := f.execOp(op, value); opEx != nil {
			f.exOr[field] = opEx
		}
	}
	return f
}

func (f *SqlArgsFilter) execOp(op QueryOp, value interface{}) goqu.Op {
	var opEx goqu.Op
	switch op {
	case EQ, NEQ, GT, LT, LTE, GTE:
		opEx = goqu.Op{string(op): value}
	case BETWEEN:
		kind := reflect.TypeOf(value).Kind()
		v := reflect.ValueOf(value)
		if kind == reflect.Slice && v.Len() == 2 {
			opEx = goqu.Op{
				"between": goqu.RangeVal{
					Start: v.Index(0).Interface(),
					End:   v.Index(1).Interface(),
				},
			}
		}
	case IN:
		opEx = goqu.Op{string(IN): value}
	case NOTIN:
		opEx = goqu.Op{string(NOTIN): value}
	case LIKE:
		value, ok := value.(string)
		if ok {
			opEx = goqu.Op{string(LIKE): "%" + value + "%"}
		}
	default:
	}

	return opEx
}

// End ...
func (f *SqlArgsFilter) End() []goqu.Expression {
	var ex []goqu.Expression
	if len(f.Ex) > 0 {
		ex = append(ex, f.Ex)
	}
	if len(f.exOr) > 0 {
		ex = append(ex, f.exOr)
	}
	return ex
}

// PageQuery ...
func (sysDB *DataBase) PageQuery(query *goqu.Dataset, scaner *gorm.DB, pageIndex int64,
	pageSize int64, outRows interface{}, selectEx ...interface{}) (int64, error) {
	var selectQuery = query
	if selectEx != nil {
		selectQuery = query.Select(selectEx...)
	}
	//count, err := sysDB.Goqu.From(selectQuery).Count()
	count, err := query.Count()
	if err != nil {
		return 0, err
	}
	selectQuery = query.
		Offset(uint((pageIndex - 1) * pageSize)).
		Limit(uint(pageSize))

	sql, args, err := selectQuery.ToSql()
	if err != nil {
		return 0, err
	}
	// use gorm to scan rows
	result := scaner.Raw(sql, args...).Find(outRows)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

// Query ...
func (sysDB *DataBase) Query(query *goqu.Dataset, scaner *gorm.DB, outRows interface{}, selectEx ...interface{}) error {
	selectQuery := query
	if selectEx != nil {
		selectQuery = query.Select(selectEx...)
	}

	sql, args, err := selectQuery.ToSql()
	if err != nil {
		return err
	}

	// use gorm to scan rows
	err = scaner.Raw(sql, args...).Find(outRows).Error
	if err != nil {
		return err
	}

	return nil
}

func (sysDB *DataBase) QueryAll(query *goqu.Dataset, outRows interface{}, selectEx ...interface{}) error {
	selectQuery := query
	if selectEx != nil {
		selectQuery = query.Select(selectEx...)
	}

	sql, args, err := selectQuery.ToSql()
	//fmt.Println("sql, args", sql, args)
	if err != nil {
		return err
	}

	// use gorm to scan rows
	err = sysDB.Raw(sql, args...).Scan(outRows).Error
	if err != nil {
		return err
	}

	return nil
}

func (sysDB *DataBase) QueryFirst(query *goqu.Dataset, outRows interface{}, selectEx ...interface{}) error {
	selectQuery := query
	selectQuery.Limit(1)

	if selectEx != nil {
		selectQuery = query.Select(selectEx...)
	}

	sql, args, err := selectQuery.ToSql()
	if err != nil {
		return err
	}

	// use gorm to scan rows
	err = sysDB.Raw(sql, args...).Scan(outRows).Error
	if err != nil {
		return err
	}

	return nil
}

func DebugSql(sqlBuilder *goqu.Dataset) {
	sql, args, err := sqlBuilder.ToSql()
	fmt.Println("Sql", sql, args, err)
}
