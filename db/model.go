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
	ModelName() string
}

var modelsList = map[string]Model{}

func DefineModel(m Model) {
	modelsList[m.ModelName()] = m
}

func CheckModelTables(dc *gorm.DB) (err error) {
	if dc == nil {
		dc, err = Open()
		if err != nil {
			return
		}
		defer dc.Close()
	}

	for _, m := range modelsList {
		if !dc.HasTable(m) {
			if err = dc.CreateTable(m).Error; err != nil {
				return
			}
		}
	}

	return
}

func scope(dc *gorm.DB, scopes []ScopeFunc) *gorm.DB {
	if scopes == nil || len(scopes) == 0 || scopes[0] == nil {
		return dc
	}

	s := make([]func(*gorm.DB) *gorm.DB, len(scopes))
	for i, scope := range scopes {
		s[i] = scope
	}

	return dc.Scopes(s...)
}

func applyQueryOptions(dcIn *gorm.DB, opt *QueryOptions) (dc *gorm.DB) {
	dc = dcIn

	if opt == nil {
		return
	}

	if opt.PreloadColNames != nil {
		for _, c := range opt.PreloadColNames {
			dc = dc.Preload(c)
		}
	}

	return
}

// ###############################################################################################
// Basic query

// Count counts the number of rows.
func Count(dc *gorm.DB, m Model, opt *QueryOptions, scopes ...ScopeFunc) (n uint, err error) {
	dc = applyQueryOptions(dc, opt)

	err = scope(dc.Model(m), scopes).Count(&n).Error
	if gorm.IsRecordNotFoundError(err) {
		err = nil
	}
	return
}

// QuerySingle finds the first matching row.
// m should be a pointer to model and result will be of same type only if anything was found otherwise it will be nil.
func QuerySingle(dc *gorm.DB, m Model, opt *QueryOptions, scopes ...ScopeFunc) (result Model, err error) {
	dc = applyQueryOptions(dc, opt)

	err = scope(dc, scopes).First(m).Error
	if gorm.IsRecordNotFoundError(err) {
		err = nil
	} else if err == nil {
		result = m
	}
	return
}

// QueryAll finds all matching rows.
// m should be a pointer to model and result will be a slice of model type (not a slice to pointer).
func QueryAll(dc *gorm.DB, m Model, opt *QueryOptions, scopes ...ScopeFunc) (result interface{}, err error) {
	dc = applyQueryOptions(dc, opt)

	t := reflect.TypeOf(m).Elem()
	slice := reflect.MakeSlice(reflect.SliceOf(t), 0, 0)
	ptr := reflect.New(slice.Type())
	ptr.Elem().Set(slice)

	err = scope(dc, scopes).Find(ptr.Interface()).Error
	if gorm.IsRecordNotFoundError(err) {
		err = nil
	} else if err == nil {
		result = ptr.Elem().Interface()
	}
	return
}

// ###############################################################################################
// CRUD

func Create(dc *gorm.DB, m Model) error {
	return dc.Create(m).Error
}

func Update(dc *gorm.DB, m Model, updateFields ...string) error {
	if len(updateFields) > 0 {
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

		return dc.Model(m).Updates(fields).Error
	}
	return dc.Save(m).Error
}

func Delete(dc *gorm.DB, m Model) error {
	return dc.Delete(m).Error
}

// ###############################################################################################
// Scope helpers

func ModelPaging(page int, limit int) ScopeFunc {
	return func(dc *gorm.DB) *gorm.DB {
		return dc.Offset((page) * limit).Limit(limit)
	}
}

func ById(id uint) ScopeFunc {
	return func(dc *gorm.DB) *gorm.DB {
		return dc.Where("id = ?", id)
	}
}

func OrderByIdDesc() ScopeFunc {
	return func(dc *gorm.DB) *gorm.DB {
		return dc.Order("id desc")
	}
}

func OrderByCreationDateDesc() ScopeFunc {
	return func(dc *gorm.DB) *gorm.DB {
		return dc.Order("created_at desc")
	}
}

func OrderByUpdateDateDesc() ScopeFunc {
	return func(dc *gorm.DB) *gorm.DB {
		return dc.Order("updated_at desc")
	}
}
