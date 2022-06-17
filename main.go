package main

import (
	bank "bankapi/bankcore"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

//DBの代わりにmapを用いて口座情報を持つ
var accounts = map[float64]*bank.Account{}

func main() {
	accounts[1001] = &bank.Account{
		Customer: bank.Customer{
			Name:    "John",
			Address: "Los Angeles, California",
			Phone:   "(213) 555 0147",
		},
		Number: 1001,
	}

	//「localhost:8000/statement」でハンドラ関数を実行
	http.HandleFunc("/statement", statement)
	//上記のハンドラ関数が実行されているエンドポイントでサーバーを起動
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

//ハンドラ関数の定義
func statement(w http.ResponseWriter, req *http.Request) {
	//リクエストのクエリパラメータに書かれた口座番号を左辺の変数（文字列型）に代入
	numberqs := req.URL.Query().Get("number")

	if numberqs == "" {
		fmt.Fprintf(w, "Account number is missing!")
		return
	}

	if number, err := strconv.ParseFloat(numberqs, 64); err != nil {
		fmt.Fprintf(w, "Invalid account number!")
	} else {
		account, ok := accounts[number]
		if !ok {
			fmt.Fprintf(w, "Account with number %v can't be found!", number)
		} else {
			fmt.Fprintf(w, "%v", account.Statement())
		}
	}
}
