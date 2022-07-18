package main

import (
	bank "bankapi/bankcore"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

//DBの代わりにmapを用いて口座情報、送金先情報を持つ
var accounts = map[float64]*CustomAccount{}

// CustomAccount ...
type CustomAccount struct {
	*bank.Account
}

func main() {
	accounts[1001] = &CustomAccount{
		Account: &bank.Account{
			Customer: bank.Customer{
				Name:    "John",
				Address: "Los Angeles, California",
				Phone:   "(213) 555 0147",
			},
			Number: 1001,
		},
	}

	accounts[1002] = &CustomAccount{
		Account: &bank.Account{
			Customer: bank.Customer{
				Name:    "Mark",
				Address: "Irvine, California",
				Phone:   "(949) 555 0198",
			},
			Number: 1002,
		},
	}

	//「localhost:8000/statement」でハンドラ関数を実行
	http.HandleFunc("/statement", statement)
	//「localhost:8000/deposit」でハンドラ関数を実行
	http.HandleFunc("/deposit", deposit)
	//「localhost:8000/withdraw」でハンドラ関数を実行
	http.HandleFunc("/withdraw", withdraw)
	//「localhost:8000/transfer」でハンドラ関数を実行
	http.HandleFunc("/transfer", transfer)
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
			json.NewEncoder(w).Encode(bank.Statement(account))
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

//ハンドラ関数の定義
//引き出しに関する関数
//例：http://localhost:8000/withdraw?number=1001&amount=100
func withdraw(w http.ResponseWriter, req *http.Request) {

	//リクエストのクエリパラメータに書かれた口座番号を左辺の変数（文字列型）に代入
	numberqs := req.URL.Query().Get("number")
	//リクエストのクエリパラメータに書かれた預金する金額を左辺の変数（文字列型）に代入
	amountqs := req.URL.Query().Get("amount")

	//口座番号が空欄の時
	if numberqs == "" {
		fmt.Fprintf(w, "%v", "Account number is missing!")
	}

	//文字列型の各変数をfloat64型にして返す
	if number, err := strconv.ParseFloat(numberqs, 64); err != nil {
		fmt.Fprintf(w, "%v", "Invalid account number!")
	} else if amount, err := strconv.ParseFloat(amountqs, 64); err != nil {
		fmt.Fprintf(w, "%v", "Invalid amount number!")
	} else {
		account, ok := accounts[number]
		if !ok {
			fmt.Fprintf(w, "Account with number %v can't be found!", number)
		} else {
			err := account.Withdraw(amount)
			if err != nil {
				fmt.Fprintf(w, "%v", err)
			} else {
				fmt.Fprintf(w, "%v", account.Statement())
			}
		}
	}
}

//ハンドラ関数の定義
//送金に関する機能
func transfer(w http.ResponseWriter, req *http.Request) {

	//リクエストのクエリパラメータに書かれた口座番号(送金元)を左辺の変数（文字列型）に代入
	numberqs := req.URL.Query().Get("number")
	//リクエストのクエリパラメータに書かれた送金する金額を左辺の変数（文字列型）に代入
	amountqs := req.URL.Query().Get("amount")
	//リクエストのクエリパラメータに書かれた送金先を左辺の変数（文字列型）に代入
	destqs := req.URL.Query().Get("dest")

	//口座番号が空欄の時
	if numberqs == "" {
		fmt.Fprintf(w, "Account number is missing!")
		return
	}

	//文字列型の各変数をfloat64型にして返す
	if number, err := strconv.ParseFloat(numberqs, 64); err != nil {
		fmt.Fprintf(w, "Invalid account number!")
	} else if amount, err := strconv.ParseFloat(amountqs, 64); err != nil {
		fmt.Fprintf(w, "Invalid amount number!")
	} else if dest, err := strconv.ParseFloat(destqs, 64); err != nil {
		fmt.Fprintf(w, "Invalid account destination number!")
	} else {
		if accountA, ok := accounts[number]; !ok {
			fmt.Fprintf(w, "Account with number %v can't be found!", number)
		} else if accountB, ok := accounts[dest]; !ok {
			fmt.Fprintf(w, "Account with number %v can't be found!", dest)
		} else {
			err := accountA.Transfer(amount, accountB.Account)
			if err != nil {
				fmt.Fprintf(w, "%v", err)
			} else {
				fmt.Fprintf(w, "%v", accountA.Statement())
			}
		}
	}
}

// Statement ...
func (c *CustomAccount) Statement() string {
	json, err := json.Marshal(c)
	if err != nil {
		return err.Error()
	}

	return string(json)
}
