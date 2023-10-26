package pkg

import "html/template"

func checkTemplate(t string) *template.Template {
	return template.Must(template.New("").Parse(t))
}

var hookTemplate = checkTemplate(`
package dao

import (
	"context"

	"github.com/IRONICBo/QiYin_BE/internal/conn/db"
	"github.com/IRONICBo/QiYin_BE/internal/dal/cache"
	"github.com/IRONICBo/QiYin_BE/internal/dal/gen"
	{{.NormPkgName}} "github.com/IRONICBo/QiYin_BE/internal/models/{{.PkgName}}"
)

// {{.ModelName}}Dao {{.NormModelName}} dao.
type {{.ModelName}}Dao struct {
	Dao
}

// New{{.ModelName}}Dao return a {{.NormModelName}} dao.
func New{{.ModelName}}Dao() *{{.ModelName}}Dao {
	query := gen.Use(db.GetMysqlDB())
	cache := cache.Use(db.GetRedis())

	return &{{.ModelName}}Dao{
		Dao: Dao{
			ctx:   context.Background(),
			query: query,
			cache: &cache,
		},
	}
}

// Create create one or multi models.
func (d *{{.ModelName}}Dao) Create(m ...*{{.NormPkgName}}.{{.ModelName}}) error {
	return d.query.WithContext(d.ctx).{{.ModelName}}.Create(m...)
}

// First get first matched result.
func (d *{{.ModelName}}Dao) First() (*{{.NormPkgName}}.{{.ModelName}}, error) {
	return d.query.WithContext(d.ctx).{{.ModelName}}.First()
}

// FindAll get all matched results.
func (d *{{.ModelName}}Dao) FindAll() ([]*{{.NormPkgName}}.{{.ModelName}}, error) {
	return d.query.WithContext(d.ctx).{{.ModelName}}.Find()
}
{{$modelName := .ModelName}}{{$normPkgName := .NormPkgName}}{{range .Fields}}
// FindFirstBy{{.FieldName}} get first matched result by {{.NormFieldName}}.
func (d *{{$modelName}}Dao) FindFirstBy{{.FieldName}}({{.NormFieldName}} {{.FieldType}}) (*{{$normPkgName}}.{{$modelName}}, error) {
	m := d.query.{{$modelName}}

	return m.WithContext(d.ctx).Where(m.{{.FieldName}}.Eq({{.NormFieldName}})).First()
}

// FindBy{{.FieldName}}Page get page by {{.FieldName}}.
func (d *{{$modelName}}Dao) FindBy{{.FieldName}}Page({{.NormFieldName}} {{.FieldType}}, offset int, limit int) ([]*{{$normPkgName}}.{{$modelName}}, int64, error) {
	m := d.query.{{$modelName}}

	result, count, err := m.WithContext(d.ctx).Where(m.{{.FieldName}}.Eq({{.NormFieldName}})).FindByPage(offset, limit)

	return result, count, err
}
{{end}}
// Update update model.
func (d *{{.ModelName}}Dao) Update(m *{{.NormPkgName}}.{{.ModelName}}) error {
	res, err := d.query.WithContext(d.ctx).{{.ModelName}}.Updates(m)
	if err != nil && res.Error != nil {
		return err
	}

	return nil
}

// Delete delete model.
func (d *{{.ModelName}}Dao) Delete(m ...*{{.NormPkgName}}.{{.ModelName}}) error {
	res, err := d.query.WithContext(d.ctx).{{.ModelName}}.Delete(m...)
	if err != nil && res.Error != nil {
		return err
	}

	return nil
}

// Count count matched records.
func (d *{{.ModelName}}Dao) Count() (int64, error) {
	return d.query.WithContext(d.ctx).{{.ModelName}}.Count()
}

///////////////////////////////////////////////////////////
//              Append your code here.                   //
///////////////////////////////////////////////////////////
`)
