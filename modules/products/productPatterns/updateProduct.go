package productPatterns

import (
	"context"
	"fmt"

	"github.com/DrumPatiphon/go-rest-api-service/modules/entities"
	"github.com/DrumPatiphon/go-rest-api-service/modules/files"
	"github.com/DrumPatiphon/go-rest-api-service/modules/files/filesUsecases"
	"github.com/DrumPatiphon/go-rest-api-service/modules/products"
	"github.com/jmoiron/sqlx"
)

type IUpdateProductBuilder interface {
	initTransaction() error
	initQuery()
	updateTitleQuery()
	updateDescriptionQuery()
	updatePriceQuery()
	updateCategory() error
	insertImages() error
	getOldImages() []*entities.Image
	deleteOldImages() error
	closeQuery()
	updateProduct() error
	getQueryFields() []string
	getvalue() []any
	getQuery() string
	setQuery(query string)
	getImagesLen() int
	commit() error
}

type updateProductbuilder struct {
	db             *sqlx.DB
	tx             *sqlx.Tx // for roll-back
	req            *products.Product
	filesUsecases  filesUsecases.IFilesUsecases
	query          string
	queryFields    []string
	lastStackIndex int
	value          []any
}

func UpdateProductBuilder(db *sqlx.DB, req *products.Product, filesUsecases filesUsecases.IFilesUsecases) IUpdateProductBuilder {
	return &updateProductbuilder{
		db:            db,
		req:           req,
		filesUsecases: filesUsecases,
		queryFields:   make([]string, 0),
		value:         make([]any, 0),
	}
}

type updateProductProductEngineer struct {
	builder IUpdateProductBuilder
}

func (b *updateProductbuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}

	b.tx = tx
	return nil
}
func (b *updateProductbuilder) initQuery() {
	b.query += `
	UPDATE products SET`
}
func (b *updateProductbuilder) updateTitleQuery() {
	if b.req.Title != "" {
		b.value = append(b.value, b.req.Title)
		b.lastStackIndex = len(b.value)

		b.queryFields = append(b.queryFields, fmt.Sprintf(`
		title = $%d`, b.lastStackIndex))
	}
}
func (b *updateProductbuilder) updateDescriptionQuery() {
	if b.req.Description != "" {
		b.value = append(b.value, b.req.Description)
		b.lastStackIndex = len(b.value)

		b.queryFields = append(b.queryFields, fmt.Sprintf(`
		description = $%d`, b.lastStackIndex))
	}
}
func (b *updateProductbuilder) updatePriceQuery() {
	if b.req.Price > 0 {
		b.value = append(b.value, b.req.Price)
		b.lastStackIndex = len(b.value)

		b.queryFields = append(b.queryFields, fmt.Sprintf(`
		price = $%d`, b.lastStackIndex))
	}
}
func (b *updateProductbuilder) updateCategory() error {
	if b.req.Category == nil {
		return nil
	}
	if b.req.Category.Id == 0 {
		return nil
	}

	query := `
	UPDATE products_categories SET
		category_id = $1
	WHERE product_id = $2;`

	if _, err := b.tx.ExecContext(
		context.Background(),
		query,
		b.req.Category.Id,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update products_categories failed: %v", err)
	}
	return nil
}
func (b *updateProductbuilder) insertImages() error {

	query := `
	INSERT INTO images (
		filename,
		url,
		product_id
	)
	VALUES
	`
	valueStack := make([]any, 0)
	var index int
	for i := range b.req.Images {
		valueStack = append(valueStack,
			b.req.Images[i].FileName,
			b.req.Images[i].Url,
			b.req.Id,
		)
		if i != len(b.req.Images)-1 {
			query += fmt.Sprintf(`
			($%d, $%d, $%d),`, index+1, index+2, index+3)
		} else {
			query += fmt.Sprintf(`
			($%d, $%d, $%d);`, index+1, index+2, index+3)
		}
		index += 3
	}

	if _, err := b.tx.ExecContext(
		context.Background(),
		query,
		valueStack...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert images failed %v", err)
	}

	return nil
}
func (b *updateProductbuilder) getOldImages() []*entities.Image {
	query := `
	SELECT
		id,
		filename,
		url
	FROM images 
	WHERE product_id = $1;`

	images := make([]*entities.Image, 0)
	if err := b.db.Select(
		&images,
		query,
		b.req.Id,
	); err != nil {
		return make([]*entities.Image, 0)
	}
	return images
}
func (b *updateProductbuilder) deleteOldImages() error {
	query := `
	DELETE FROM images
	WHERE product_id = $1;`

	images := b.getOldImages()
	if len(images) > 0 {
		deleteFileReq := make([]*files.DeleteFileReq, 0)
		for _, img := range images {
			deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
				Destination: fmt.Sprintf("image/test/%s", img.FileName),
			})
		}
		if err := b.filesUsecases.DeleteFileOnGCP(deleteFileReq); err != nil {
			b.tx.Rollback()
			return err
		}
	}

	if _, err := b.tx.ExecContext(
		context.Background(),
		query,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("failed to delete images: %v", err)
	}
	return nil
}
func (b *updateProductbuilder) closeQuery() {
	b.value = append(b.value, b.req.Id)
	b.lastStackIndex = len(b.value)

	b.query += fmt.Sprintf(`
	WHERE id = $%d`, b.lastStackIndex)
}
func (b *updateProductbuilder) updateProduct() error {
	if _, err := b.tx.ExecContext(
		context.Background(),
		b.query,
		b.value...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("uapdte products failed: %v", err)
	}
	return nil
}
func (b *updateProductbuilder) getQueryFields() []string {
	return b.queryFields
}
func (b *updateProductbuilder) getvalue() []any {
	return b.value
}
func (b *updateProductbuilder) getQuery() string {
	return b.query
}
func (b *updateProductbuilder) setQuery(query string) {
	b.query = query
}
func (b *updateProductbuilder) getImagesLen() int {
	return len(b.req.Images)
}
func (b *updateProductbuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdateProductProductEngineer(b IUpdateProductBuilder) *updateProductProductEngineer {
	return &updateProductProductEngineer{
		builder: b,
	}
}

func (en *updateProductProductEngineer) sumQueryFields() {
	en.builder.updateTitleQuery()
	en.builder.updateDescriptionQuery()
	en.builder.updatePriceQuery()

	fields := en.builder.getQueryFields()

	for i := range fields {
		query := en.builder.getQuery()
		if i != len(fields)-1 {
			en.builder.setQuery(query + fields[i] + ",")
		} else {
			en.builder.setQuery(query + fields[i])
		}
	}
}

func (en *updateProductProductEngineer) UpdateProduct() error {
	en.builder.initTransaction()

	en.builder.initQuery()
	en.sumQueryFields()
	en.builder.closeQuery()

	// Update Product
	if err := en.builder.updateProduct(); err != nil {
		return nil
	}

	// update category
	if err := en.builder.updateCategory(); err != nil {
		return err
	}

	if en.builder.getImagesLen() > 0 {
		if err := en.builder.deleteOldImages(); err != nil {
			return err
		}
		if err := en.builder.insertImages(); err != nil {
			return err
		}
	}

	// Commit
	if err := en.builder.commit(); err != nil {
		return err
	}

	return nil
}
