package service

import (
	"github.com/lemuzhi/cart/domain/model"
	"github.com/lemuzhi/cart/domain/repository"
)

type ICartDataService interface {
	AddCart(cart *model.Cart) (int64, error)
	DeleteCart(int64) error
	UpdateCart(cart *model.Cart) error
	FindCartByID(int64) (*model.Cart, error)
	FindAllCart(userID int64) ([]model.Cart, error)

	CleanCart(userID int64) error
	IncrNum(cartID int64, num int64) error
	DecrNum(cartID int64, num int64) error
}

type CartDataService struct {
	CartRepository repository.ICartRepository
}

func NewCartDataService(cartRepository repository.ICartRepository) ICartDataService {
	return &CartDataService{CartRepository: cartRepository}
}

func (c *CartDataService) AddCart(cart *model.Cart) (int64, error) {
	return c.CartRepository.CreateCart(cart)
}

func (c *CartDataService) DeleteCart(cartID int64) error {
	return c.CartRepository.DeleteCartByID(cartID)
}

func (c *CartDataService) UpdateCart(cart *model.Cart) error {
	return c.CartRepository.UpdateCart(cart)
}

func (c *CartDataService) FindCartByID(cartID int64) (*model.Cart, error) {
	return c.CartRepository.FindCartByID(cartID)
}

func (c *CartDataService) FindAllCart(userID int64) ([]model.Cart, error) {
	return c.CartRepository.FindAll(userID)
}

func (c *CartDataService) CleanCart(userID int64) error {
	return c.CartRepository.CleanCart(userID)
}

func (c *CartDataService) IncrNum(cartID int64, num int64) error {
	return c.CartRepository.IncrNum(cartID, num)
}

func (c *CartDataService) DecrNum(cartID int64, num int64) error {
	return c.CartRepository.DecrNum(cartID, num)
}
