package repository

import (
	"errors"
	"github.com/lemuzhi/cart/domain/model"
	"gorm.io/gorm"
)

type ICartRepository interface {
	InitTable() error
	FindCartByID(int64) (*model.Cart, error)
	CreateCart(cart *model.Cart) (int64, error)
	DeleteCartByID(int64) error
	UpdateCart(cart *model.Cart) error
	FindAll(userID int64) (cartAll []model.Cart, err error)

	CleanCart(userID int64) error
	IncrNum(cartID int64, num int64) error
	DecrNum(cartID int64, num int64) error
}

type CartRepository struct {
	mysqlDb *gorm.DB
}

// 创建cartRepository
func NewCartRepository(db *gorm.DB) ICartRepository {
	return &CartRepository{mysqlDb: db}
}

// 初始化表
func (c *CartRepository) InitTable() error {
	return c.mysqlDb.Migrator().CreateTable(&model.Cart{})
}

// 根据id查询Cart信息
func (c *CartRepository) FindCartByID(cartID int64) (cart *model.Cart, err error) {
	cart = &model.Cart{}
	return cart, c.mysqlDb.First(cart, cartID).Error
}

// 创建Cart信息
func (c *CartRepository) CreateCart(cart *model.Cart) (int64, error) {
	//根据条件查询是否有数据，如果没有则创建
	db := c.mysqlDb.FirstOrCreate(cart, model.Cart{ProductID: cart.ProductID, SizeID: cart.SizeID, UserID: cart.UserID})
	if db.Error != nil {
		return 0, db.Error
	}
	//处理的行数是否等于0，如果为0，则表示插入失败
	if db.RowsAffected == 0 {
		return 0, errors.New("购物车插入失败")
	}
	return cart.ID, nil
}

// 根据id删除Cart信息
func (c *CartRepository) DeleteCartByID(cartID int64) error {
	return c.mysqlDb.Where("id=?", cartID).Delete(&model.Cart{}).Error
}

// 更新Cart信息
func (c *CartRepository) UpdateCart(cart *model.Cart) error {
	return c.mysqlDb.Model(cart).Updates(cart).Error
}

// 查询所有Cart信息
func (c *CartRepository) FindAll(userID int64) (cartAll []model.Cart, err error) {
	return cartAll, c.mysqlDb.Where("user_id=?", userID).Find(&cartAll).Error
}

// 根据ID清空购物车
func (c *CartRepository) CleanCart(userID int64) error {
	return c.mysqlDb.Where("user_id = ?", userID).Delete(&model.Cart{}).Error
}

// 添加商品数量
func (c *CartRepository) IncrNum(cartID int64, num int64) error {
	cart := &model.Cart{
		ID: cartID,
	}
	return c.mysqlDb.Model(cart).UpdateColumn("num", gorm.Expr("num + ?", num)).Error
}

// 减少商品数量
func (c *CartRepository) DecrNum(cartID int64, num int64) error {
	cart := &model.Cart{
		ID: cartID,
	}
	db := c.mysqlDb.Model(cart).Where("num >= ?", num).UpdateColumn("num", gorm.Expr("num - ?", num))

	if db.Error != nil {
		return db.Error
	}

	if db.RowsAffected == 0 {
		return errors.New("减少商品失败")
	}
	return nil
}
