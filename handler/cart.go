package handler

import (
	"context"
	"github.com/lemuzhi/cart/domain/model"
	"github.com/lemuzhi/cart/domain/service"
	"github.com/lemuzhi/cart/proto"
	"github.com/lemuzhi/common"
)

type Cart struct {
	CartDataService service.ICartDataService
}

// 添加购物车
func (c *Cart) AddCart(ctx context.Context, req *proto.CartInfo, resp *proto.ResponseAdd) (err error) {
	cart := &model.Cart{}
	err = common.SwapTo(req, cart)
	if err != nil {
		return err
	}
	resp.CartId, err = c.CartDataService.AddCart(cart)
	return err
}

// 清空购物车
func (c *Cart) CleanCart(ctx context.Context, req *proto.Clean, resp *proto.Response) error {
	err := c.CartDataService.CleanCart(req.UserId)
	if err != nil {
		return err
	}
	resp.Msg = "清空购物车成功"
	return nil
}

// 添加购物车数量成功
func (c *Cart) Incr(ctx context.Context, req *proto.Item, resp *proto.Response) error {
	err := c.CartDataService.IncrNum(req.Id, req.ChangeNum)
	if err != nil {
		return err
	}
	resp.Msg = "购物车数量添加成功"
	return nil
}

// 减少购物车失败
func (c *Cart) Decr(ctx context.Context, req *proto.Item, resp *proto.Response) error {
	err := c.CartDataService.DecrNum(req.Id, req.ChangeNum)
	if err != nil {
		return err
	}
	resp.Msg = "购物车数量减少成功"
	return nil
}

// 删除购物车
func (c *Cart) DeleteItemByID(ctx context.Context, req *proto.CartID, resp *proto.Response) error {
	err := c.CartDataService.DeleteCart(req.Id)
	if err != nil {
		return err
	}
	resp.Msg = "购物车删除成功"
	return nil
}

// 查询用户所有的购物车信息
func (c *Cart) GetAll(ctx context.Context, req *proto.CartFindAll, resp *proto.CartAll) error {
	allCart, err := c.CartDataService.FindAllCart(req.UserId)
	if err != nil {
		return err
	}
	for _, v := range allCart {
		car := &proto.CartInfo{}
		err = common.SwapTo(v, car)
		if err != nil {
			return err
		}
		resp.CartInfo = append(resp.CartInfo, car)
	}
	return nil
}
