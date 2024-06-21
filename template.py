article_content = """[<- home](/)
    
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

index_page_md = """I'm Hidetatz, a software and AI enthusiast.
I'm currently developing [whale](https://github.com/hidetatz/whale), a deep learning framework inspired by PyTorch but entirely written in Go.
Currently AI technology is for researchers rathar than for developers.
As building AI into the systems becomes common, it must be changed. Building AI model must be more diverse, and that's why I develop whale.

See some other open source software authored by me:

* [shiba](https://github.com/hidetatz/shiba)

shiba is a programming language which is plain like Python, but modern like Go or Rust.
Python is great, but some parts (e.g. package manager, code formatting) must be updated and that's why I create shiba.

* [kubecolor](https://github.com/hidetatz/kubecolor) (publicly archived)

kubecolor is a CLI tool which colorizes the kubectl output for readability.
It's archived because I am not using Kubernetes for my work. I'll restart maitaining it if I have to run some containers on Kubernetes again (though I don't hope so).
You can use the [community fork](https://github.com/kubecolor/kubecolor) if wanted.

* [incdb](https://github.com/hidetatz/incdb)

This is an incrementally developed RDBMS from scratch. This is work-in-progress project.

* [rv](https://github.com/hidetatz/rv)

RISC-V software emulator. WIP.

Visit my [GitHub](https://github.com/hidetatz) for more information.

[Atom/RSS feed](/feed.xml)

---

## 書いたもの

$articles

---

© 2024 Hidetatz Yaginuma. Unless otherwise noted, these posts are made available under a [Creative Commons Attribution License](https://creativecommons.org/licenses/by/4.0/)."""

diary_content = """[<- 戻る](/)
    
# $title

$content"""

diary_index_page_md = """日記です。  
$diaries

---

© 2024 Hidetatz Yaginuma. Unless otherwise noted, these posts are made available under a [Creative Commons Attribution License](https://creativecommons.org/licenses/by/4.0/)."""

not_found_page_md = """
# 404: Page Not Found

This page doesn't exist. Recent artices are below or go to the [top page](/).

$recent_articles
"""
