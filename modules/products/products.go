package products

import (
	"github.com/DrumPatiphon/go-rest-api-service/modules/appInfo"
	"github.com/DrumPatiphon/go-rest-api-service/modules/entities"
)

type Product struct {
	Id          string            `json: "id"`
	Title       string            `json: "title"`
	Description string            `json: "description"`
	Category    *appInfo.Category `json: "category"`
	CreatedAt   string            `json: "created_at"`
	UpdatedAt   string            `json: "updated_at"`
	Price       float64           `json: "price"`
	Images      []*entities.Image `json: "images"`
}

type ProductFilter struct {
	Id                      string `query:"id"`
	Search                  string `query:"search"`
	*entities.PaginationReq        // like inherit class
	*entities.SortReq
}
