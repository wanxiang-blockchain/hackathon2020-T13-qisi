package check

import (
	"context"
	"fisco/build/SettledChain"
	"fmt"
	"github.com/chislab/go-fiscobcos/common"
)

var SettledAddr common.Address
var dsettled *SettledChain.SettledChain

func Settled() error {
	ctx := context.Background()
	auth := NewAuthFromPriKey()
	addr, tx, _dented, err := SettledChain.DeploySettledChain(auth, GethCli)
	if err != nil {
		return err
	}
	SettledAddr = addr
	dsettled = _dented

	fmt.Println("deploy dsettled txHash", tx.Hash().String())
	fmt.Println("deploy SettledAddr addr", SettledAddr.String())

	tx, err = dsettled.RoomRegister(auth, common.HexToAddress("0xDF12793CA392ff748adF013D146f8dA73df6E304"), []byte("lululu1936"), false, 120, []byte("30"), []byte("0"), []byte("goodplace"))
	if err != nil {
		return err
	}
	fmt.Println("call dsettled.RoomRegister txHash", tx.Hash().String())
	GethCli.CheckTx(ctx, tx)
	owner, location, isAvailable, price, area, decorateStatus, description, err := dsettled.GetRoom(callOpts, []byte("lululu1936"))
	if err != nil {
		return err
	}

	fmt.Println("get dsettled room info:::")
	fmt.Println(owner, location, isAvailable, price, area, decorateStatus, description)

	return nil
}
