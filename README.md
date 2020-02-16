## blog

https://dtyler.io

## How to write Article

English article:

```shell
$ echo title of article---$(date "+%F %T") > ./data/articles/url_path.md
```

Japanese article:

```shell
$ echo title of article---$(date "+%F %T") > ./data/articles/ja/url_path.md
```

External link:

```shell
$ echo title of article---$(date "+%F %T")---https://url.of.the.link > ./data/articles/url_path.md
```

## How to generate blog

```shell
$ go run main.go gen
```
