// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package gen

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/IRONICBo/QiYin_BE/internal/models/base"
)

func newUserBase(db *gorm.DB, opts ...gen.DOOption) userBase {
	_userBase := userBase{}

	_userBase.userBaseDo.UseDB(db, opts...)
	_userBase.userBaseDo.UseModel(&base.UserBase{})

	tableName := _userBase.userBaseDo.TableName()
	_userBase.ALL = field.NewAsterisk(tableName)
	_userBase.Id = field.NewUint(tableName, "id")
	_userBase.CreatedAt = field.NewTime(tableName, "created_at")
	_userBase.UpdatedAt = field.NewTime(tableName, "updated_at")
	_userBase.DeletedAt = field.NewTime(tableName, "deleted_at")
	_userBase.UUID = field.NewString(tableName, "uuid")
	_userBase.Email = field.NewString(tableName, "email")
	_userBase.Nickname = field.NewString(tableName, "nickname")
	_userBase.Avatar = field.NewString(tableName, "avatar")
	_userBase.Description = field.NewString(tableName, "description")
	_userBase.IsEnable = field.NewBool(tableName, "is_enable")

	_userBase.fillFieldMap()

	return _userBase
}

type userBase struct {
	userBaseDo userBaseDo

	ALL         field.Asterisk
	Id          field.Uint
	CreatedAt   field.Time
	UpdatedAt   field.Time
	DeletedAt   field.Time
	UUID        field.String
	Email       field.String
	Nickname    field.String
	Avatar      field.String
	Description field.String
	IsEnable    field.Bool

	fieldMap map[string]field.Expr
}

func (u userBase) Table(newTableName string) *userBase {
	u.userBaseDo.UseTable(newTableName)
	return u.updateTableName(newTableName)
}

func (u userBase) As(alias string) *userBase {
	u.userBaseDo.DO = *(u.userBaseDo.As(alias).(*gen.DO))
	return u.updateTableName(alias)
}

func (u *userBase) updateTableName(table string) *userBase {
	u.ALL = field.NewAsterisk(table)
	u.Id = field.NewUint(table, "id")
	u.CreatedAt = field.NewTime(table, "created_at")
	u.UpdatedAt = field.NewTime(table, "updated_at")
	u.DeletedAt = field.NewTime(table, "deleted_at")
	u.UUID = field.NewString(table, "uuid")
	u.Email = field.NewString(table, "email")
	u.Nickname = field.NewString(table, "nickname")
	u.Avatar = field.NewString(table, "avatar")
	u.Description = field.NewString(table, "description")
	u.IsEnable = field.NewBool(table, "is_enable")

	u.fillFieldMap()

	return u
}

func (u *userBase) WithContext(ctx context.Context) IUserBaseDo { return u.userBaseDo.WithContext(ctx) }

func (u userBase) TableName() string { return u.userBaseDo.TableName() }

func (u userBase) Alias() string { return u.userBaseDo.Alias() }

func (u userBase) Columns(cols ...field.Expr) gen.Columns { return u.userBaseDo.Columns(cols...) }

func (u *userBase) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := u.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (u *userBase) fillFieldMap() {
	u.fieldMap = make(map[string]field.Expr, 10)
	u.fieldMap["id"] = u.Id
	u.fieldMap["created_at"] = u.CreatedAt
	u.fieldMap["updated_at"] = u.UpdatedAt
	u.fieldMap["deleted_at"] = u.DeletedAt
	u.fieldMap["uuid"] = u.UUID
	u.fieldMap["email"] = u.Email
	u.fieldMap["nickname"] = u.Nickname
	u.fieldMap["avatar"] = u.Avatar
	u.fieldMap["description"] = u.Description
	u.fieldMap["is_enable"] = u.IsEnable
}

func (u userBase) clone(db *gorm.DB) userBase {
	u.userBaseDo.ReplaceConnPool(db.Statement.ConnPool)
	return u
}

func (u userBase) replaceDB(db *gorm.DB) userBase {
	u.userBaseDo.ReplaceDB(db)
	return u
}

type userBaseDo struct{ gen.DO }

type IUserBaseDo interface {
	gen.SubQuery
	Debug() IUserBaseDo
	WithContext(ctx context.Context) IUserBaseDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IUserBaseDo
	WriteDB() IUserBaseDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IUserBaseDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IUserBaseDo
	Not(conds ...gen.Condition) IUserBaseDo
	Or(conds ...gen.Condition) IUserBaseDo
	Select(conds ...field.Expr) IUserBaseDo
	Where(conds ...gen.Condition) IUserBaseDo
	Order(conds ...field.Expr) IUserBaseDo
	Distinct(cols ...field.Expr) IUserBaseDo
	Omit(cols ...field.Expr) IUserBaseDo
	Join(table schema.Tabler, on ...field.Expr) IUserBaseDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IUserBaseDo
	RightJoin(table schema.Tabler, on ...field.Expr) IUserBaseDo
	Group(cols ...field.Expr) IUserBaseDo
	Having(conds ...gen.Condition) IUserBaseDo
	Limit(limit int) IUserBaseDo
	Offset(offset int) IUserBaseDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IUserBaseDo
	Unscoped() IUserBaseDo
	Create(values ...*base.UserBase) error
	CreateInBatches(values []*base.UserBase, batchSize int) error
	Save(values ...*base.UserBase) error
	First() (*base.UserBase, error)
	Take() (*base.UserBase, error)
	Last() (*base.UserBase, error)
	Find() ([]*base.UserBase, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*base.UserBase, err error)
	FindInBatches(result *[]*base.UserBase, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*base.UserBase) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IUserBaseDo
	Assign(attrs ...field.AssignExpr) IUserBaseDo
	Joins(fields ...field.RelationField) IUserBaseDo
	Preload(fields ...field.RelationField) IUserBaseDo
	FirstOrInit() (*base.UserBase, error)
	FirstOrCreate() (*base.UserBase, error)
	FindByPage(offset int, limit int) (result []*base.UserBase, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IUserBaseDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (u userBaseDo) Debug() IUserBaseDo {
	return u.withDO(u.DO.Debug())
}

func (u userBaseDo) WithContext(ctx context.Context) IUserBaseDo {
	return u.withDO(u.DO.WithContext(ctx))
}

func (u userBaseDo) ReadDB() IUserBaseDo {
	return u.Clauses(dbresolver.Read)
}

func (u userBaseDo) WriteDB() IUserBaseDo {
	return u.Clauses(dbresolver.Write)
}

func (u userBaseDo) Session(config *gorm.Session) IUserBaseDo {
	return u.withDO(u.DO.Session(config))
}

func (u userBaseDo) Clauses(conds ...clause.Expression) IUserBaseDo {
	return u.withDO(u.DO.Clauses(conds...))
}

func (u userBaseDo) Returning(value interface{}, columns ...string) IUserBaseDo {
	return u.withDO(u.DO.Returning(value, columns...))
}

func (u userBaseDo) Not(conds ...gen.Condition) IUserBaseDo {
	return u.withDO(u.DO.Not(conds...))
}

func (u userBaseDo) Or(conds ...gen.Condition) IUserBaseDo {
	return u.withDO(u.DO.Or(conds...))
}

func (u userBaseDo) Select(conds ...field.Expr) IUserBaseDo {
	return u.withDO(u.DO.Select(conds...))
}

func (u userBaseDo) Where(conds ...gen.Condition) IUserBaseDo {
	return u.withDO(u.DO.Where(conds...))
}

func (u userBaseDo) Order(conds ...field.Expr) IUserBaseDo {
	return u.withDO(u.DO.Order(conds...))
}

func (u userBaseDo) Distinct(cols ...field.Expr) IUserBaseDo {
	return u.withDO(u.DO.Distinct(cols...))
}

func (u userBaseDo) Omit(cols ...field.Expr) IUserBaseDo {
	return u.withDO(u.DO.Omit(cols...))
}

func (u userBaseDo) Join(table schema.Tabler, on ...field.Expr) IUserBaseDo {
	return u.withDO(u.DO.Join(table, on...))
}

func (u userBaseDo) LeftJoin(table schema.Tabler, on ...field.Expr) IUserBaseDo {
	return u.withDO(u.DO.LeftJoin(table, on...))
}

func (u userBaseDo) RightJoin(table schema.Tabler, on ...field.Expr) IUserBaseDo {
	return u.withDO(u.DO.RightJoin(table, on...))
}

func (u userBaseDo) Group(cols ...field.Expr) IUserBaseDo {
	return u.withDO(u.DO.Group(cols...))
}

func (u userBaseDo) Having(conds ...gen.Condition) IUserBaseDo {
	return u.withDO(u.DO.Having(conds...))
}

func (u userBaseDo) Limit(limit int) IUserBaseDo {
	return u.withDO(u.DO.Limit(limit))
}

func (u userBaseDo) Offset(offset int) IUserBaseDo {
	return u.withDO(u.DO.Offset(offset))
}

func (u userBaseDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IUserBaseDo {
	return u.withDO(u.DO.Scopes(funcs...))
}

func (u userBaseDo) Unscoped() IUserBaseDo {
	return u.withDO(u.DO.Unscoped())
}

func (u userBaseDo) Create(values ...*base.UserBase) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Create(values)
}

func (u userBaseDo) CreateInBatches(values []*base.UserBase, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (u userBaseDo) Save(values ...*base.UserBase) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Save(values)
}

func (u userBaseDo) First() (*base.UserBase, error) {
	if result, err := u.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*base.UserBase), nil
	}
}

func (u userBaseDo) Take() (*base.UserBase, error) {
	if result, err := u.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*base.UserBase), nil
	}
}

func (u userBaseDo) Last() (*base.UserBase, error) {
	if result, err := u.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*base.UserBase), nil
	}
}

func (u userBaseDo) Find() ([]*base.UserBase, error) {
	result, err := u.DO.Find()
	return result.([]*base.UserBase), err
}

func (u userBaseDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*base.UserBase, err error) {
	buf := make([]*base.UserBase, 0, batchSize)
	err = u.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (u userBaseDo) FindInBatches(result *[]*base.UserBase, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return u.DO.FindInBatches(result, batchSize, fc)
}

func (u userBaseDo) Attrs(attrs ...field.AssignExpr) IUserBaseDo {
	return u.withDO(u.DO.Attrs(attrs...))
}

func (u userBaseDo) Assign(attrs ...field.AssignExpr) IUserBaseDo {
	return u.withDO(u.DO.Assign(attrs...))
}

func (u userBaseDo) Joins(fields ...field.RelationField) IUserBaseDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Joins(_f))
	}
	return &u
}

func (u userBaseDo) Preload(fields ...field.RelationField) IUserBaseDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Preload(_f))
	}
	return &u
}

func (u userBaseDo) FirstOrInit() (*base.UserBase, error) {
	if result, err := u.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*base.UserBase), nil
	}
}

func (u userBaseDo) FirstOrCreate() (*base.UserBase, error) {
	if result, err := u.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*base.UserBase), nil
	}
}

func (u userBaseDo) FindByPage(offset int, limit int) (result []*base.UserBase, count int64, err error) {
	result, err = u.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = u.Offset(-1).Limit(-1).Count()
	return
}

func (u userBaseDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = u.Count()
	if err != nil {
		return
	}

	err = u.Offset(offset).Limit(limit).Scan(result)
	return
}

func (u userBaseDo) Scan(result interface{}) (err error) {
	return u.DO.Scan(result)
}

func (u userBaseDo) Delete(models ...*base.UserBase) (result gen.ResultInfo, err error) {
	return u.DO.Delete(models)
}

func (u *userBaseDo) withDO(do gen.Dao) *userBaseDo {
	u.DO = *do.(*gen.DO)
	return u
}
