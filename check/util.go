package check

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/chislab/go-fiscobcos/accounts/abi/bind"
	"github.com/chislab/go-fiscobcos/client"
	"github.com/chislab/go-fiscobcos/common"
	"github.com/chislab/go-fiscobcos/core/types"
	"github.com/chislab/go-fiscobcos/crypto"
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
	var err error
	conf := &client.Config{
		CAFile:     "./nodes/127.0.0.1/sdk/ca.crt",
		CertFile:   "./nodes/127.0.0.1/sdk/node.crt",
		KeyFile:    "./nodes/127.0.0.1/sdk/node.key",
		Endpoint:   "127.0.0.1:20200",
		UseChannel: true,
		GroupID:    1,
	}
	GethCli, err = client.New(conf)
	if err != nil {
		println("error:", err.Error())
		os.Exit(1)
	}
	//height, err := GethCli.BlockNumber(context.Background())
	//if err != nil {
	//	println("error:", err.Error())
	//	os.Exit(1)
	//}
	//fmt.Println("Current block height is", height.String())
}

func str2Big(str string) *big.Int {
	return new(big.Int).SetBytes([]byte(str))
}
func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
func WaitMinedByHash(txHash common.Hash) *types.Receipt {
	ctx := context.Background()
	queryTicker := time.NewTicker(time.Millisecond * 200)
	defer queryTicker.Stop()
	for {
		receipt, _ := GethCli.TransactionReceipt(ctx, txHash)
		if receipt != nil {
			return receipt
		}
		// Wait for the next round.
		select {
		case <-ctx.Done():
			return nil
		case <-queryTicker.C:
		}
	}
}

func NewAuthFromPriKey(priKey ...string) *bind.TransactOpts {
	var priv *ecdsa.PrivateKey
	if len(priKey) == 0 {
		priv, _ = crypto.GenerateKey()
	} else {
		priv, _ = crypto.HexToECDSA(priKey[0])
	}
	auth := bind.NewKeyedTransactor(priv, 1, 1)
	auth.Context = context.Background()
	return auth
}

func ilog(contract string, format string, v ...interface{}) {
	log.Printf("[%s | INFO]: %s", strings.ToUpper(contract), fmt.Sprintf(format, v...))
}

func getReceiptOutput(output string) string {
	if strings.HasPrefix(output, "0x") {
		output = output[2:]
	}
	b, err := hex.DecodeString(output)
	if err != nil || len(b) < 36 {
		return output
	}
	b = b[36:]
	tail := len(b) - 1
	for ; tail >= 0; tail-- {
		if b[tail] != 0 {
			break
		}
	}
	return string(b[:tail+1])
}

func checkTx(tx *types.Transaction, err error) {
	if err != nil {
		panic(err)
	}
	receipt, err = func(tx *types.Transaction, err error) (*types.Receipt, error) {
		receipt := WaitMinedByHash(tx.Hash())
		if receipt.Status != "0x0" {
			return receipt, fmt.Errorf("receipt.Status = %s\nTxHash = %s\nOutput = %s", receipt.Status, receipt.TxHash.String(), getReceiptOutput(receipt.Output))
		}
		return receipt, nil
	}(tx, err)
	if err != nil {
		panic(err)
	}
	return
}
