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
	//「localhost:8000/deposit」でハンドラ関数を実行
	http.HandleFunc("/deposit", deposit)
	//上記のハンドラ関数が実行されているエンドポイントでサーバーを起動
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

//ハンドラ関数の定義
//HTTPの呼び出しを通じて口座番号を送信する関数
//例：http://localhost:8000/statement?number=1001
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

//ハンドラ関数の定義
//預金に関する関数
//例：http://localhost:8000/deposit?number=1001&amount=100
func deposit(w http.ResponseWriter, req *http.Request) {

	//リクエストのクエリパラメータに書かれた口座番号を左辺の変数（文字列型）に代入
	numberqs := req.URL.Query().Get("number")
	//リクエストのクエリパラメータに書かれた預金する金額を左辺の変数（文字列型）に代入
	amountqs := req.URL.Query().Get("amount")

	if numberqs == "" {
		fmt.Fprintf(w, "Account number is missing!")
		return
	}

	if number, err := strconv.ParseFloat(numberqs, 64); err != nil {
		fmt.Fprintf(w, "Invalid account number!")
	} else if amount, err := strconv.ParseFloat(amountqs, 64); err != nil {
		fmt.Fprintf(w, "Invalid amount number!")
	} else {
		account, ok := accounts[number]
		if !ok {
			fmt.Fprintf(w, "Account with number %v can't be found!", number)
		} else {
			err := account.Deposit(amount)
			if err != nil {
				fmt.Fprintf(w, "%v", err)
			} else {
				fmt.Fprintf(w, "%v", account.Statement())
			}
		}
	}
}
