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
	"strconv"
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
	_, err = dLeaseHold.WatchEvtOrderMade(watchOpts, chOrderMake)
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
	err = GethCli.CheckTx(context.Background(), tx)
	if err != nil {
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}
	evt := <- chOrderMake
	fmt.Println("OrderId", evt.OrderId)
	c.JSON(200, gin.H{"tx_hash": tx.Hash()})
}

func HandleConfirmOrder(c *gin.Context) {
	var req ReqConfirm
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, nil)
		return
	}
	tx, err = dLeaseHold.ConfirmOrder(auths[req.Property], big.NewInt(req.OrderId))
	err = GethCli.CheckTx(context.Background(), tx)
	if err != nil {
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}
	c.JSON(200, gin.H{"tx_hash": tx.Hash()})
}

func HandleBalance(c *gin.Context) {
	user := c.PostForm("user")
	balance, _ :=dLeaseHold.BalanceOf(callOpts, auths[user].From)
	c.JSON(200, gin.H{"user": user, "balance": balance.Uint64()})
}

func HandleDeviceLog(c *gin.Context) {
	device := c.PostForm("device")
	logId := c.PostForm("log_id")
	id, err := strconv.Atoi(logId)
	if err != nil {
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}
	oneLog, _ :=dLeaseHold.GetLog(callOpts, auths[device].From, big.NewInt(int64(id)))
	c.JSON(200, oneLog)
}

// 房东注册房子
func HandleRegisterRoom(c *gin.Context) {
	var req ReqRegisterRoom
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, nil)
		return
	}

	tx, _ = dLeaseHold.RoomRegister(auths[req.Landlord], auths[req.Property].From, auths[req.Factory].From,
		str2bytes(req.Location), big.NewInt(req.Price), str2bytes(req.Area),
		req.Status, str2bytes(req.Description), []common.Address{})
	err = GethCli.CheckTx(context.Background(), tx)
	if err != nil {
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}
	c.JSON(200, gin.H{"tx_hash": tx.Hash()})
}

// 平台/厂商装修房子
func HandleFixRoom(c *gin.Context) {
	var req ReqRegisterRoom
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, nil)
		return
	}

	tx, err = dLeaseHold.UpdateRoomInfo(auths[req.Property], auths[req.Property].From, auths[req.Factory].From,
		str2bytes(req.Location), big.NewInt(req.Price), str2bytes(req.Area),
		req.Status, str2bytes(req.Description), []common.Address{})
	if err != nil {
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}
	c.JSON(200, gin.H{"tx_hash": tx.Hash()})
}
