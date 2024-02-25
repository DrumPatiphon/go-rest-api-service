package productsRepositories

import (
	"encoding/json"
	"fmt"

	"github.com/DrumPatiphon/go-rest-api-service/config"
	"github.com/DrumPatiphon/go-rest-api-service/modules/entities"
	"github.com/DrumPatiphon/go-rest-api-service/modules/files/filesUsecases"
	"github.com/DrumPatiphon/go-rest-api-service/modules/products"
	"github.com/DrumPatiphon/go-rest-api-service/modules/products/productPatterns"
	"github.com/jmoiron/sqlx"
)

type IProductRepository interface {
	FindOneProduct(productId string) (*products.Product, error)
	FindProduct(req *products.ProductFilter) ([]*products.Product, int)
	InsertProduct(req *products.Product) (*products.Product, error)
	UpdateProduct(req *products.Product) (*products.Product, error)
}

type productRepository struct {
	db           *sqlx.DB
	cfg          config.Iconfig
	filesUsecase filesUsecases.IFilesUsecases
}

func ProductRepository(db *sqlx.DB, cfg config.Iconfig, filesUsecase filesUsecases.IFilesUsecases) IProductRepository {
	return &productRepository{
		db:           db,
		cfg:          cfg,
		filesUsecase: filesUsecase,
	}
}

func (r *productRepository) FindOneProduct(productId string) (*products.Product, error) {
	query := `
	SELECT 
		to_jsonb(t)
	FROM (
		SELECT
			p.id,
			p.title,
			p.description,
			p.price,
			(
				SELECT
					to_jsonb(ct)
				FROM (
					SELECT
						c.id,
						c.title
					FROM categories c
					LEFT JOIN products_categories pc ON pc.category_id  = c.id
					WHERE pc.product_id = $1
				) AS ct
			) AS category,
			p.created_at,
			p.updated_at,
			p.price,
			(
				SELECT 
					COALESCE(array_to_json(array_agg(it)), '[]'::json)  
				FROM (
					SELECT 
						i.id,
						i.filename,
						i.url
					FROM images i 
					WHERE i.product_id = p.id 
				)AS it
			)AS images
		FROM products p
		WHERE p.id = $1
	) AS t;`

	productBytes := make([]byte, 0)
	product := &products.Product{
		Images: make([]*entities.Image, 0),
	}

	if err := r.db.Get(&productBytes, query, productId); err != nil {
		return nil, fmt.Errorf("get product failed %v", err)
	}
	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, fmt.Errorf("unmarshal product failed: %v", err)
	}
	return product, nil
}

func (r *productRepository) FindProduct(req *products.ProductFilter) ([]*products.Product, int) {
	builder := productPatterns.FindProductBuilder(r.db, req)
	engineer := productPatterns.FindProductEngineer(builder)

	result := engineer.FindProduct().Result()
	count := engineer.CountProduct().Count()
	return result, count

}

func (r *productRepository) InsertProduct(req *products.Product) (*products.Product, error) {
	builder := productPatterns.InsertProductBuilder(r.db, req)
	productId, err := productPatterns.InsertProductEngineer(builder).InsertProduct()
	if err != nil {
		return nil, err
	}

	product, err := r.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productRepository) UpdateProduct(req *products.Product) (*products.Product, error) {
	builder := productPatterns.UpdateProductBuilder(r.db, req, r.filesUsecase)
	engineer := productPatterns.UpdateProductProductEngineer(builder)

	if err := engineer.UpdateProduct(); err != nil {
		return nil, err
	}

	product, err := r.FindOneProduct(req.Id)
	if err != nil {
		return nil, err
	}

	return product, nil
}
