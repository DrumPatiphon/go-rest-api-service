package productsUsecases

import (
	"math"

	"github.com/DrumPatiphon/go-rest-api-service/modules/entities"
	"github.com/DrumPatiphon/go-rest-api-service/modules/products"
	"github.com/DrumPatiphon/go-rest-api-service/modules/products/productsRepositories"
)

type IProductUseCase interface {
	FindOneProduct(productId string) (*products.Product, error)
	FindProduct(req *products.ProductFilter) *entities.PageRes
	InsertProduct(req *products.Product) (*products.Product, error)
	UpdateProduct(req *products.Product) (*products.Product, error)
	DeleteProduct(productId string) error
}

type productsUsecases struct {
	productRepository productsRepositories.IProductRepository
}

func ProductsUsecases(productRepository productsRepositories.IProductRepository) IProductUseCase {
	return &productsUsecases{
		productRepository: productRepository,
	}
}

func (u *productsUsecases) FindOneProduct(productId string) (*products.Product, error) {
	product, err := u.productRepository.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (u *productsUsecases) FindProduct(req *products.ProductFilter) *entities.PageRes {
	products, count := u.productRepository.FindProduct(req)

	return &entities.PageRes{
		Data:       products,
		Page:       req.Page,
		Limit:      req.Page,
		TotalItems: count,
		TotalPage:  int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *productsUsecases) InsertProduct(req *products.Product) (*products.Product, error) {
	product, err := u.productRepository.InsertProduct(req)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (u *productsUsecases) UpdateProduct(req *products.Product) (*products.Product, error) {
	product, err := u.productRepository.UpdateProduct(req)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (u *productsUsecases) DeleteProduct(productId string) error {
	if err := u.productRepository.DeleteProduct(productId); err != nil {
		return err
	}
	return nil
}
