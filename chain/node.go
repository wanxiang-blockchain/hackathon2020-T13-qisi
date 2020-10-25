package chain

import (
	"context"
	"github.com/chislab/go-fiscobcos/accounts/abi/bind"
	"github.com/chislab/go-fiscobcos/client"
	"github.com/chislab/go-fiscobcos/common"
	"github.com/chislab/go-fiscobcos/core/types"
	"github.com/gin-gonic/gin"
	"os"
)

var (
	callOpts = &bind.CallOpts{GroupId: 1, From: common.HexToAddress("0x100")}
	watchOpts = &bind.WatchOpts{Start: new(uint64), Context: context.Background()}
	GethCli  *client.Client
	tx       *types.Transaction
	err      error
	receipt  *types.Receipt
)

func init() {
	sdkPath := "./nodes/127.0.0.1/sdk/"
	conf := &client.Config{
		CAFile:     sdkPath + "ca.crt",
		CertFile:   sdkPath + "node.crt",
		KeyFile:    sdkPath + "node.key",
		Endpoint:   "127.0.0.1:20200",
		UseChannel: true,
		GroupID:    1,
	}
	GethCli, err = client.New(conf)
	if err != nil {
		println("error:", err.Error())
		os.Exit(1)
	}
}

func HandleMakeOrder(c *gin.Context) {
	height, err := GethCli.BlockNumber(context.Background())
	if err != nil {

	}
	c.JSON(200, height.String())
}

