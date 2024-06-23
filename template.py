article_content = """[<- ホーム](/)
    
# $title

#### $timestamp

$content"""

html_page = """<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>$title</title>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta name="author" content="Hidetatz Yaginuma">
  <meta name="viewport" content="width=device-width, initial-scale=1, minimal-ui">

  <style>
    body {
      box-sizing: border-box;
      min-width: 200px;
      max-width: 980px;
      margin: 0 auto;
      padding: 45px;
    }
  </style>

  <link href="/markdown.css" rel="stylesheet"></link>
  <link href="/syntax.css" rel="stylesheet"></link>
  <script type="text/javascript" async
    src="https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.1/MathJax.js?config=TeX-MML-AM_CHTML">
  </script>
</head>
<body class="markdown-body">
$body

<script src="/highlight.pack.js"></script>
<script>hljs.highlightAll();</script>
</body>
</html>"""

index_page_md = """## Hidetatz Web Page

Hidetatzは計算機、ソフトウェア、AI領域のエンジニアです。

[GitHub](https://github.com/hidetatz)

[Atom/RSSフィード](/feed.xml)

## プロジェクト

### [whale](https://github.com/hidetatz/whale)

Goで書かれたディープラーニングフレームワーク。

### [shiba](https://github.com/hidetatz/shiba)

Pythonのようにプレーンで、RustやGoのようにモダンなプログラミング言語。

### [incdb](https://github.com/hidetatz/incdb)

フルスクラッチで実装されたRDBMS。

### [kubecolor](https://github.com/hidetatz/kubecolor) (アーカイブ)

kubectlを100倍便利にするソフトウェア。

---

## 書いたもの

$articles

[日記](/diary)

---

© 2024 Hidetatz Yaginuma. Unless otherwise noted, these posts are made available under a [Creative Commons Attribution License](https://creativecommons.org/licenses/by/4.0/)."""

diary_content = """[<- 日記](/diary)
    
# $title

$content"""

diary_index_page_md = """[<- ホーム](/)

日記です。  

$diaries

---

© 2024 Hidetatz Yaginuma. Unless otherwise noted, these posts are made available under a [Creative Commons Attribution License](https://creativecommons.org/licenses/by/4.0/)."""

not_found_page_md = """
# 404: Page Not Found

[トップに戻る](/)

## 最近の記事

$recent_articles
"""
