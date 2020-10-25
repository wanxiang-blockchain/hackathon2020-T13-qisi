package chain

import (
	"context"
	"crypto/ecdsa"
	"github.com/chislab/go-fiscobcos/accounts/abi/bind"
	"github.com/chislab/go-fiscobcos/crypto"
	"math/big"
	"unsafe"
)

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

func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func str2Big(str string) *big.Int {
	return new(big.Int).SetBytes([]byte(str))
}
func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
