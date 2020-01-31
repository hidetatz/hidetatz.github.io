Goでテスト以外でもTable Drivenに書く---2018-11-09 23:36:31

Goのテストにおいては[TableDrivenTests](https://github.com/golang/go/wiki/TableDrivenTests)という技法がある。
Goの世界では非常に一般的なのでぜひ書き方を覚えたほうがいい。
主に、テストの追加/削除の容易さや、可読性の高さがメリットである。

実は、この書き方はテスト以外でも適用できる。

例えば、userモデルを保存する際に、バリデーションを定義するとしよう。
[こんなコード](https://play.golang.org/p/Tg59fLcEeiJ)を書いた。

```go
package main

import (
	"fmt"
	"os"
	"unicode/utf8"
)

type email string

func (e *email) valid() bool {
	return false // 本当はちゃんと正規表現でチェックする
}

type user struct {
	name  string
	age   int
	email email
}

func (u *user) validate() error {
	if utf8.RuneCountInString(u.name) > 20 {
		return fmt.Errorf("name must be less than 20")
	}
	if u.age < 18 {
		return fmt.Errorf("age must be more than 18")
	}
	if !u.email.valid() {
		return fmt.Errorf("email format is invalid")
	}
	return nil
}

func (u *user) save() error {
	// ここでDBアクセスする
	return nil
}

func main() {
	u := &user{name: "Jeff", age: 17, email: "jeff@example.com"}
	if err := u.validate(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
	u.save()
}
```

この `validate()` は、Table Drivenに、[以下のように](https://play.golang.org/p/74ODzAHP0bU)リファクタできる。

```
func (u *user) validate() error {
	checks := []struct {
		invalid bool
		errMsg  string
	}{
		{utf8.RuneCountInString(u.name) > 20, "name must be less than 20"},
		{u.age < 18, "age must be more than 18"},
		{!u.email.valid(), "email format is invalid"},
	}

	for _, check := range checks {
		if check.invalid {
			return fmt.Errorf("invalid: %s", check.errMsg)
		}
	}
	return nil
}
```

これによって、userの要素が増えてもバリデーションが追加しやすくなる。
オススメの書き方。
