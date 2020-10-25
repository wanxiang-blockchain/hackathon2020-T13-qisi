package main

import (
	"github.com/gin-gonic/gin"
	"leasehold/chain"
)

func main() {
	r := gin.Default()

	r.POST("/make_order", chain.HandleMakeOrder)
	r.POST("/confirm_order", chain.HandleConfirmOrder)
	r.GET("/balance", chain.HandleBalance)
	r.POST("/register_room", chain.HandleRegisterRoom)
	r.POST("/update_room", chain.HandleFixRoom)

	r.Run()
}

