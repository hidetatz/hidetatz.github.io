karabiner-elements(公式)が進化してた---2017-10-17 15:09:56

自分がkarabinerを使って実現したいのは、左右のCommandを英数/かなにするのと、SandS。
今までは、これを参考に、本家のkarabiner-elementsをforkしたPRを使っていた。
いつしか公式でこれらができるようになっていたのでやり方をメモしておく。

### やり方

公式からDLする。(執筆時点で11.0.0)

DLして展開したら、以下のコマンドで設定ファイルを編集する。

```
vi ~/.config/karabiner/karabiner.json
```

```
"manipulators": [
  {
      "from": {
          "key_code": "spacebar",
          "modifiers": {
              "optional": [
                  "any"
              ]
          }
      },
      "to": [
          {
              "key_code": "left_shift"
          }
      ],
      "to_if_alone": [
          {
              "key_code": "spacebar"
          }
      ],
      "type": "basic"
  },
  {
      "from": {
          "key_code": "left_command",
          "modifiers": {
              "optional": [
                  "any"
              ]
          }
      },
      "to": [
          {
              "key_code": "left_command"
          }
      ],
      "to_if_alone": [
          {
              "key_code": "japanese_eisuu"
          }
      ],
      "type": "basic"
  },
  {
      "from": {
          "key_code": "right_command",
          "modifiers": {
              "optional": [
                  "any"
              ]
          }
      },
      "to": [
          {
              "key_code": "right_command"
          }
      ],
      "to_if_alone": [
          {
              "key_code": "japanese_kana"
          }
      ],
      "type": "basic"
  }
```
