package chain

import (
	"context"
	"fmt"
	"github.com/chislab/go-fiscobcos/accounts/abi/bind"
	"github.com/chislab/go-fiscobcos/client"
	"github.com/chislab/go-fiscobcos/common"
	"github.com/chislab/go-fiscobcos/core/types"
	"github.com/gin-gonic/gin"
	"leasehold/contracts/leasehold"
	"math/big"
	"os"
)

var (
	GethCli *client.Client

	callOpts  = &bind.CallOpts{GroupId: 1, From: common.HexToAddress("0x100")}
	watchOpts = &bind.WatchOpts{Start: new(uint64), Context: context.Background()}
	tx        *types.Transaction
	err       error
	receipt   *types.Receipt
	auths     = make(map[string]*bind.TransactOpts)

	dLeaseHold  *leasehold.Leasehold
	chOrderMake = make(chan *leasehold.LeaseholdEvtOrderMade)
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

	auths["admin"] = NewAuthFromPriKey()
	auths["landlord"] = NewAuthFromPriKey()
	auths["tenantry"] = NewAuthFromPriKey()
	auths["property"] = NewAuthFromPriKey()
	auths["device"] = NewAuthFromPriKey()
	auths["factory"] = NewAuthFromPriKey()

	_, tx, dLeaseHold, err = leasehold.DeployLeasehold(auths["admin"], GethCli)
	fmt.Println("DeployLeasehold", tx.Hash().String())
}

func HandleMakeOrder(c *gin.Context) {
	var req ReqOrder
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, nil)
		return
	}
	// from, to, location, startAt, endAt, funds
	tx, err = dLeaseHold.MakeOrder(auths[req.From], auths[req.To].From,
		str2bytes(req.Location), req.StartAt, req.EndAt, big.NewInt(req.Funds))
	c.JSON(200, gin.H{"tx_hash": tx.Hash()})
}
