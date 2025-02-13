package main

import (
	v1 "challenge-yuno/cmd/api/v1"
	"challenge-yuno/internal/business/usecases/order"
	"challenge-yuno/internal/platform/repositories/kvstore"
	"challenge-yuno/internal/platform/repositories/sql"
	"challenge-yuno/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=postgres user=user password=password dbname=postgres port=5432 sslmode=disable TimeZone=America/Argentina/Mendoza"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Errorf("error connecting to db %v", err)
		panic(err)
	}
	db.Debug()

	kvsOrderRepo := kvstore.NewOrderRepository()
	sqlOrderRepo := sql.NewOrderRepository(db)

	notificationService := services.NewNotificationService("whatsapp")
	orderUsecase := order.NewOrderUsecase(kvsOrderRepo, sqlOrderRepo, notificationService)

	e := echo.New()

	e.Debug = true
	e.HideBanner = true

	v1.NewOrderHandler(e, orderUsecase)

	e.Logger.Fatal(e.Start(":8080"))
}
