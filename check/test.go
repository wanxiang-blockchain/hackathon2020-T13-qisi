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
	ctx := context.Background()
	auth := NewAuthFromPriKey()
	_, tx, dLeaseHold, err := leasehold.DeployLeasehold(auth, GethCli)
	if err != nil {
		return err
	}
	go parseEvt()

	_, err = dLeaseHold.WatchEvtOrderMade(watchOpts, chOrderMake)

	landlord := NewAuthFromPriKey()
	tenantry := NewAuthFromPriKey()
	property := NewAuthFromPriKey()
	factory := NewAuthFromPriKey()

	addr := str2bytes("上海市虹口区海伦路111号11栋1号1单元103室")
	fmt.Println("deploy Leasehold txHash", tx.Hash().String())
	fmt.Println("Start registar Room.");
	tx,err = dLeaseHold.RoomRegister( auth, property.From,
		factory.From, addr, big.NewInt(200),
		str2bytes("30"),uint8(2),str2bytes("这个房子很nice"));
	fmt.Println("RoomRegister txHash", tx.Hash().String())
	GethCli.CheckTx(ctx, tx)
	leaseholdRoom, err := dLeaseHold.GetRoom(callOpts, addr);
	fmt.Printf("GetRoom Info:::  %v", leaseholdRoom,bytes2str(leaseholdRoom.Description))
	fmt.Println();

	tx, err = dLeaseHold.MakeOrder(tenantry, landlord.From, addr,
		1, 12, big.NewInt(100))
	fmt.Println("MakeOrder txHash", tx.Hash().String())

	GethCli.CheckTx(context.Background(), tx)

	ok, err := dLeaseHold.GetOrderStatus(callOpts, big.NewInt(0))
	fmt.Println("IsRoomValid,room info is, \n房屋位置:",bytes2str(ok.Location),
		"\n订单金额:",ok.Funds, ",\n预定者:",ok.From.String(), ",\n物业信息:",ok.Property.String())
	//fmt.Println("GetOrderStatus",ok.)

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
