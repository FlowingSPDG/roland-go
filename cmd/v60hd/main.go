package main

import (
	"fmt"
	"os"

	"github.com/FlowingSPDG/roland-go/v60hd"
)

func main() {
	// V-60HDのIPアドレスとポート番号
	ipAddress := "192.168.2.254" // ご使用のV-60HDのIPアドレスに置き換えてください
	port := "8023"

	// 切り替え先のチャンネル (例: SDI IN 1 は 0, SDI IN 2 は 1)
	channel := 1

	v, err := v60hd.NewV60HD(ipAddress, port)
	if err != nil {
		fmt.Println("V-60HDの接続エラー:", err)
		os.Exit(1)
	}
	defer v.Close()

	if err := v.PGM(channel); err != nil {
		fmt.Println("PGMコマンドの送信エラー:", err)
		os.Exit(1)
	}

	fmt.Println("処理を完了しました")
}
