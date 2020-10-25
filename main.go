package main

import (
	"github.com/gin-gonic/gin"
	"leasehold/chain"
)

func main() {
	r := gin.Default()

	r.POST("/make_order", chain.HandleMakeOrder)

	r.Run()
}

