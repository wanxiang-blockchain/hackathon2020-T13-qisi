package check

import (
	"context"
	"fisco/build/leasehold"
	"fmt"
	"math/big"
)

var dLeaseHold *leasehold.Leasehold
var chOrderMake = make(chan *leasehold.LeaseholdEvtOrderMade, 1)

func Test() error {
	auth := NewAuthFromPriKey()
	_, tx, dLeaseHold, err := leasehold.DeployLeasehold(auth, GethCli)
	if err != nil {
		return err
	}
	go parseEvt()

	_, err = dLeaseHold.WatchEvtOrderMade(watchOpts, chOrderMake)

	landlord := NewAuthFromPriKey()
	tenantry := NewAuthFromPriKey()
	addr := str2bytes("上海市虹口区海伦路111号11栋1号1单元101室")
	fmt.Println("deploy Leasehold txHash", tx.Hash().String())
	tx, err = dLeaseHold.MakeOrder(tenantry, landlord.From, addr,
		1, 12, big.NewInt(100))
	fmt.Println("MakeOrder txHash", tx.Hash().String())

	GethCli.CheckTx(context.Background(), tx)

	ok, err := dLeaseHold.GetOrderStatus(callOpts, big.NewInt(0))
	fmt.Println("IsRoomValid", ok)

	return nil
}

func parseEvt() {
	for {
		select {
			case evt := <- chOrderMake:
				fmt.Println("chOrderMake", evt.Amount)
		}
	}
}
