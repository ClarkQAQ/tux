package main

import (
	"context"
	"fmt"
	"tux/rest"
)

type Services struct{}

func (*Services) AdminService() rest.Groupor {
	return rest.Group("/admin",
		new(AdminService), // 管理接口
	).Tags("管理员接口")
}

type AdminService struct{}

func (*AdminService) BaseService() rest.Groupor {
	return rest.Group("/base",
		new(BaseService), // 基础服务
	).Tags("基础服务")
}

func (*AdminService) AdminMiddleware() rest.Middlewareor {
	return rest.Middleware(func(ctx context.Context) {
		fmt.Println("QAQ")
	})
}

type BaseService struct{}

func (*BaseService) UserService() rest.Groupor {
	return rest.Group("/user",
		new(UserService), // 用户管理
	).Tags("用户管理")
}

type UserService struct{}

func (*UserService) Middleware() rest.Middlewareor {
	return rest.Middleware(func(ctx context.Context) {
	})
}

func (*UserService) Hack() rest.Hackor {
	return rest.Hack(rest.ResponseHack, func(ctx context.Context) {
		fmt.Println("QAQ")
	})
}

type PostReq struct {
}

func (*UserService) Post() rest.Methodor {
	return rest.Method(func(ctx context.Context, req *struct {
		Name string `json:"name" form:"name" description:"名称"`
	}) (bool, error) {
		fmt.Println(req.Name)

		return req.Name == "admin", nil
	})
}

type GetReq struct {
	Id int64 `json:"id" form:"id" query:"id" validate:"required" description:"编号"`

	Limit int64  `json:"limit" form:"limit" query:"limit" description:"单页数量"`
	Page  int64  `json:"page" form:"page" query:"page" description:"页码"`
	Sort  string `json:"sort" form:"sort" query:"sort" default:"id" enum:"id,date" description:"排序字段"`
}

func (*UserService) Get() rest.Methodor {
	return rest.Method(func(ctx context.Context, req *GetReq) (*string, error) {
		return nil, nil
	}).Summary("获取信息").Description("获取用户信息")
}

func main() {
	rest.Test(new(Services))
}
