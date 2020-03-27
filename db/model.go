package db

import (
	"reflect"
	"time"

	"github.com/jinzhu/gorm"
)

// Gorm tags:
// many2many: 							join table name
// jointable_foreignkey:				source col name
// association_jointable_foreignkey:	target col name

type ScopeFunc func(*gorm.DB) *gorm.DB

type QueryOptions struct {
	PreloadColNames []string
}

func QP(cols ...string) *QueryOptions {
	return &QueryOptions{PreloadColNames: cols}
}

type BaseModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type Model interface {
	TableName() string
}

var modelsList = map[string]Model{}

func DefineModel(m Model) {
	modelsList[m.TableName()] = m
}

func CheckModelTables(db *gorm.DB) (err error) {
	if db == nil {
		db, err = Open()
		if err != nil {
			return
		}
		defer db.Close()
	}

	for _, m := range modelsList {
		if !db.HasTable(m) {
			if err = db.CreateTable(m).Error; err != nil {
				return
			}
		}
	}

	return
}

func scope(db *gorm.DB, scopes []ScopeFunc) *gorm.DB {
	if scopes == nil || len(scopes) == 0 || scopes[0] == nil {
		return db
	}

	s := make([]func(*gorm.DB) *gorm.DB, len(scopes))
	for i, scope := range scopes {
		s[i] = scope
	}

	return db.Scopes(s...)
}

// ###############################################################################################
// Basic query

func applyQueryOptions(dbIn *gorm.DB, opt *QueryOptions) (db *gorm.DB) {
	db = dbIn

	if opt == nil {
		return
	}

	if opt.PreloadColNames != nil {
		for _, c := range opt.PreloadColNames {
			db = db.Preload(c)
		}
	}

	return
}

// FindSingle finds the first matching row.
// m should be a pointer to model and result will be of same type only if anything was found otherwise it will be nil.
func FindSingle(db *gorm.DB, m Model, opt *QueryOptions, scopes ...ScopeFunc) (result Model, err error) {
	db = applyQueryOptions(db, opt)

	err = scope(db, scopes).First(m).Error
	if gorm.IsRecordNotFoundError(err) {
		err = nil
	} else if err == nil {
		result = m
	}
	return
}

func Count(db *gorm.DB, m Model, opt *QueryOptions, scopes ...ScopeFunc) (n uint, err error) {
	db = applyQueryOptions(db, opt)

	err = scope(db.Model(m), scopes).Count(&n).Error
	if gorm.IsRecordNotFoundError(err) {
		err = nil
	}
	return
}

// FindAll finds all matching rows.
// m should be a pointer to model and result will be a slice of model type (not a slice to pointer).
func FindAll(db *gorm.DB, m Model, opt *QueryOptions, scopes ...ScopeFunc) (result interface{}, err error) {
	db = applyQueryOptions(db, opt)

	t := reflect.TypeOf(m).Elem()
	slice := reflect.MakeSlice(reflect.SliceOf(t), 0, 0)
	ptr := reflect.New(slice.Type())
	ptr.Elem().Set(slice)

	err = scope(db, scopes).Find(ptr.Interface()).Error
	if gorm.IsRecordNotFoundError(err) {
		err = nil
	} else if err == nil {
		result = ptr.Elem().Interface()
	}
	return
}

// ###############################################################################################
// CRUD

func Create(db *gorm.DB, m Model) error {
	return db.Create(m).Error
}

func Update(db *gorm.DB, m Model, updateFields []string) error {
	if updateFields != nil {
		fields := map[string]interface{}{}
		mt := reflect.TypeOf(m).Elem()
		mv := reflect.ValueOf(m).Elem()
		for _, f := range updateFields {
		ufLoop:
			for i := 0; i < mt.NumField(); i++ {
				ft := mt.Field(i)

				if f == ft.Name {
					fields[ft.Name] = mv.Field(i).Interface()
					break ufLoop
				}

				boundName := ft.Tag.Get("json")
				if boundName == "" {
					boundName = ft.Tag.Get("form")
					if boundName == "" {
						continue
					}
				}

				if f == boundName {
					fields[ft.Name] = mv.Field(i).Interface()
					break ufLoop
				}
			}
		}

		return db.Model(m).Updates(fields).Error
	}
	return db.Save(m).Error
}

func Delete(db *gorm.DB, m Model) error {
	return db.Delete(m).Error
}

// ###############################################################################################
// Scope helpers

func ModelPaging(page int, limit int) ScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((page) * limit).Limit(limit)
	}
}

func ById(id uint) ScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func OrderByIdDesc() ScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("id desc")
	}
}

func OrderByCreationDateDesc() ScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc")
	}
}

func OrderByUpdateDateDesc() ScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("updated_at desc")
	}
}
