package main

import (
	"fmt"
	"github.com/go-micro/plugins/v4/registry/consul"
	"github.com/go-micro/plugins/v4/wrapper/ratelimiter/uber"
	"github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/lemuzhi/cart/domain/repository"
	"github.com/lemuzhi/cart/domain/service"
	"github.com/lemuzhi/cart/handler"
	pb "github.com/lemuzhi/cart/proto"
	"github.com/lemuzhi/common"
	"go-micro.dev/v4/registry"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"

	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

var (
	QPS = 100
)

func main() {
	//配置中心
	consulConfig, err := common.GetConsulConfig("121.40.63.97", 8500, "micro/config")
	if err != nil {
		log.Println(err)
	}
	//注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"121.40.63.97:8500",
		}
	})

	//获取mysql配置， 路径中不带前缀
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")

	//初始化数据库连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlInfo.User, mysqlInfo.Pwd, mysqlInfo.Host, mysqlInfo.Port, mysqlInfo.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: &schema.NamingStrategy{
			SingularTable: true,
		},
		//禁用事物
		//SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Println(err)
	}

	//初始化表，只执行一次
	err = repository.NewCartRepository(db).InitTable()
	if err != nil {
		log.Println(err)
	}

	//链路追踪
	jargerTracer, i, err := common.NewTracer("cart", "121.40.63.97:8500")
	if err != nil {
		return
	}
	defer i.Close()

	// Create service
	srv := micro.NewService()
	srv.Init(
		micro.Name("cart"),
		micro.Version("latest"),
		micro.Address("127.0.0.1:8087"),
		//添加注册中心
		micro.Registry(consulRegistry),
		//添加链路追踪
		micro.WrapHandler(opentracing.NewHandlerWrapper(jargerTracer)),
		//添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
	)

	srv.Init()

	cartDataService := service.NewCartDataService(repository.NewCartRepository(db))

	// Register handler
	if err = pb.RegisterCartHandler(srv.Server(), &handler.Cart{CartDataService: cartDataService}); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err = srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
