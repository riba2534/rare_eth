package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/skip2/go-qrcode"
)

func randomKey() *ecdsa.PrivateKey {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	return privateKey
}

func generateWalletWithPrefixSuffix(ctx context.Context, prefix, suffix string, found chan<- *ecdsa.PrivateKey) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			privateKey := randomKey()
			address := crypto.PubkeyToAddress(privateKey.PublicKey)
			addressStr := strings.ToLower(address.Hex()[2:])

			if strings.HasPrefix(addressStr, prefix) && strings.HasSuffix(addressStr, suffix) {
				select {
				case found <- privateKey:
					return
				case <-ctx.Done():
					return
				}
			}
		}
	}
}

func printQRCode(text string) {
	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		return
	}

	fmt.Println("二维码:")
	fmt.Println(qr.ToSmallString(false))
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: <prefix> <suffix> <num_goroutines>")
		return
	}

	prefix := os.Args[1]
	suffix := os.Args[2]
	numGoroutines, err := strconv.Atoi(os.Args[3])
	if err != nil || numGoroutines < 1 {
		fmt.Println("Error: num_goroutines must be a positive integer")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	found := make(chan *ecdsa.PrivateKey)
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			generateWalletWithPrefixSuffix(ctx, prefix, suffix, found)
			wg.Done()
		}()
	}

	privateKey := <-found
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKey))

	fmt.Printf("找到了满足条件的钱包地址：%s\n", address.Hex())
	fmt.Printf("对应的私钥是：%s\n", privateKeyHex)

	printQRCode(privateKeyHex)

	cancel()
	wg.Wait()
}
