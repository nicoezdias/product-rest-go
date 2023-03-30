package handler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"clase19/internal/domain"
	"clase19/internal/product"
	"clase19/pkg/web"

	"github.com/gin-gonic/gin"
)

type productHandler struct {
	s product.Service
}

// NewProductHandler crea un nuevo controller de productos
func NewProductHandler(s product.Service) *productHandler {
	return &productHandler{
		s: s,
	}
}

// GetAll godoc
// @Summary      Get all products
// @Description  Get all products from repository
// @Tags         products
// @Produce      json
// @Param        token header string true "token"
// @Success      200 {object}  web.response
// @Router       /products [get]
func (h *productHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		products, _ := h.s.GetAll()
		web.Success(c, 200, products)
	}
}

// GetByID godoc
// @Summary      Get a product by Id
// @Description  Get a product by Id from repository
// @Tags         products
// @Produce      json
// @Param        token header string true "token"
// @Param        id   path      int  true  "Product Id"
// @Success      200 {object}  web.response
// @Failure      400 {object}  web.errorResponse
// @Failure      404 {object}  web.errorResponse
// @Router       /products/:id [get]
func (h *productHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid id"))
			return
		}
		product, err := h.s.GetByID(id)
		if err != nil {
			web.Failure(c, 404, errors.New("product not found"))
			return
		}
		web.Success(c, 200, product)
	}
}

// Search godoc
// @Summary      Get  products by price
// @Description  Get  products whose price is greater than a value from repository
// @Tags         products
// @Produce      json
// @Param        token header string true "token"
// @Param        priceGt   query      float64  true  "Price Gt"
// @Success      200 {object}  web.response
// @Failure      400 {object}  web.errorResponse
// @Failure      404 {object}  web.errorResponse
// @Router       /products/search [get]
func (h *productHandler) Search() gin.HandlerFunc {
	return func(c *gin.Context) {
		priceParam := c.Query("priceGt")
		price, err := strconv.ParseFloat(priceParam, 64)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid price"))
			return
		}
		products, err := h.s.SearchPriceGt(price)
		if err != nil {
			web.Failure(c, 404, errors.New("product not found"))
			return
		}
		web.Success(c, 200, products)
	}
}

// ConsumerPrice godoc
// @Summary      Returns a price and a list
// @Description  Returns the price of a list of products and the list
// @Tags         products
// @Produce      json
// @Param        token header string true "token"
// @Param        list   query      []int  true  "List Ids"
// @Success      200 {object}  web.response
// @Failure      400 {object}  web.errorResponse
// @Failure      404 {object}  web.errorResponse
// @Router       /products/consumer_price [get]
func (h *productHandler) ConsumerPrice() gin.HandlerFunc {
	return func(c *gin.Context) {
		type response struct {
			products    []domain.Product
			total_price float64
		}
		ids := c.Query("list")
		ids = strings.Replace(ids, "[", "", -1)
		ids = strings.Replace(ids, "]", "", -1)
		listIds := strings.Split(string(ids), ",")
		var listIdsInt []int
		for _, v := range listIds {
			id, err := strconv.Atoi(v)
			if err != nil {
				web.Failure(c, 400, errors.New("invalid id"))
				return
			}
			listIdsInt = append(listIdsInt, id)
		}
		products, price, err := h.s.ConsumerPrice(listIdsInt)
		if err != nil {
			web.Failure(c, 400, err)
			return
		}
		// data := response{products, price}
		// web.Success(c, 200, data)
		c.JSON(200, gin.H{
			"products":    products,
			"total_price": price,
		})
	}
}

// Post godoc
// @Summary      Create a new product
// @Description  Create a new product in repository
// @Tags         products
// @Produce      json
// @Param        token header string true "token"
// @Param        body body domain.Product true "Product"
// @Success      201 {object}  web.response
// @Failure      400 {object}  web.errorResponse
// @Failure      404 {object}  web.errorResponse
// @Router       /products [post]
func (h *productHandler) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		var product domain.Product
		err := c.ShouldBindJSON(&product)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid json"))
			return
		}
		valid, err := validateEmptys(&product)
		if !valid {
			web.Failure(c, 400, err)
			return
		}
		valid, err = validateExpiration(&product)
		if !valid {
			web.Failure(c, 400, err)
			return
		}
		p, err := h.s.Create(product)
		if err != nil {
			web.Failure(c, 400, err)
			return
		}
		web.Success(c, 201, p)
	}
}

// Put godoc
// @Summary      Update a product by id
// @Description  Update a product by id in repository
// @Tags         products
// @Produce      json
// @Param        token header string true "token"
// @Param        body body domain.Product true "Product"
// @Param        id   path      int  true  "Product Id"
// @Success      200 {object}  web.response
// @Failure      400 {object}  web.errorResponse
// @Failure      404 {object}  web.errorResponse
// @Router       /products/:id [put]
func (h *productHandler) Put() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid id"))
			return
		}
		var product domain.Product
		err = c.ShouldBindJSON(&product)
		if err != nil {
			web.Failure(c, 400, err)
			return
		}
		valid, err := validateEmptys(&product)
		if !valid {
			web.Failure(c, 400, err)
			return
		}
		valid, err = validateExpiration(&product)
		if !valid {
			web.Failure(c, 400, err)
			return
		}
		p, err := h.s.UpdateProduct(id, product)
		if err != nil {
			web.Failure(c, 400, err)
			return
		}
		web.Success(c, 200, p)
	}
}

// Patch godoc
// @Summary      Update a product
// @Description  Update a product by id in repository
// @Tags         products
// @Produce      json
// @Param        token header string true "token"
// @Param        body body domain.Product true "Product"
// @Param        id   path      int  true  "Product Id"
// @Success      200 {object}  web.response
// @Failure      400 {object}  web.errorResponse
// @Failure      404 {object}  web.errorResponse
// @Router       /products/:id [patch]
func (h *productHandler) Patch() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid id"))
			return
		}
		var product domain.Product
		err = c.BindJSON(&product)
		if err != nil {
			web.Failure(c, 400, err)
			return
		}
		valid, err := validateExpiration(&product)
		if !valid {
			web.Failure(c, 400, err)
			return
		}
		p, err := h.s.UpdateProduct(id, product)
		if err != nil {

			web.Failure(c, 400, err)
			return
		}
		web.Success(c, 200, p)
	}
}

// Delete elimina un producto por su id
// Delete godoc
// @Summary      Delete a product
// @Description  Delete a product by id in repository
// @Tags         products
// @Produce      json
// @Param        token header string true "token"
// @Param        id   path      int  true  "Product Id"
// @Success      200 {object}  web.response
// @Failure      400 {object}  web.errorResponse
// @Failure      404 {object}  web.errorResponse
// @Router       /products/:id [delete]
func (h *productHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid id"))
			return
		}
		err = h.s.Delete(id)
		if err != nil {
			web.Failure(c, 404, err)
			return
		}
		web.Success(c, 204, fmt.Sprintf("user %d deleted", id))
	}
}

/* ---------------------------------- Utils --------------------------------- */

// validateEmptys valida que los campos no esten vacios
func validateEmptys(product *domain.Product) (bool, error) {
	switch {
	case product.Name == "":
		return false, errors.New("name can't be empty")
	case product.CodeValue == "":
		return false, errors.New("code_value can't be empty")
	case product.Expiration == "":
		return false, errors.New("expiration can't be empty")
	case product.Quantity <= 0 || product.Price <= 0:
		if product.Quantity <= 0 {
			return false, errors.New("quantity must be greater than 0")
		}
		if product.Price <= 0 {
			return false, errors.New("price must be greater than 0")
		}
	}
	return true, nil
}

// validateExpiration valida que la fecha de expiracion sea valida
func validateExpiration(product *domain.Product) (bool, error) {
	dates := strings.Split(product.Expiration, "/")
	list := []int{}
	if len(dates) != 3 {
		return false, errors.New("invalid expiration date, must be in format: dd/mm/yyyy")
	}
	for value := range dates {
		number, err := strconv.Atoi(dates[value])
		if err != nil {
			return false, errors.New("invalid expiration date, must be numbers")
		}
		list = append(list, number)
	}
	condition := (list[0] < 1 || list[0] > 31) && (list[1] < 1 || list[1] > 12) && (list[2] < 1 || list[2] > 9999)
	if condition {
		return false, errors.New("invalid expiration date, date must be between 1 and 31/12/9999")
	}
	return true, nil
}
