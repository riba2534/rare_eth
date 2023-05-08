package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
)

func genETHWallet() *ecdsa.PrivateKey {
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
			privateKey := genETHWallet()
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
	var prefix, suffix string
	var numGoroutines int

	var rootCmd = &cobra.Command{
		Use:   "./rare_eth",
		Short: "ETH 钱包靓号生成器，可以指定钱包地址的 前缀 和 后缀，支持指定线程数",
		Run: func(cmd *cobra.Command, args []string) {
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
		},
	}

	rootCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "需要的钱包地址的前缀, 不指定则为不限制")
	rootCmd.Flags().StringVarP(&suffix, "suffix", "s", "", "需要的钱包地址的后缀, 不指定则为不限制")
	rootCmd.Flags().IntVarP(&numGoroutines, "numGoroutines", "n", 100, "线程数量, 不指定默认为  100")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
