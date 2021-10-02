Docker x Capybara x Selenium x Headless ChromeでE2Eテストを書く---2017-09-07 09:18:13

Docker上にrubyが動くコンテナを作って、E2Eテストを書くやり方を残しておく。

### Dockerの用意

ベースイメージには Ubuntu を使用。

```
FROM ubuntu:14.04

### set timezone JST
RUN /bin/cp -p /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
```

まずはJSTをセット。

```
### install tools
RUN apt-get update && \
    apt-get install -y \
    git curl wget unzip vim language-pack-ja-base language-pack-ja \
    # used to apt-add-repository
    software-properties-common \
    # used to bundle install
    make zlib1g-dev \
    # used to dpkg google-chrome
    gconf-service libasound2 libatk1.0-0 libcairo2 libcups2 \
    libfontconfig1 libgdk-pixbuf2.0-0 libgtk-3-0 libnspr4 libx11-xcb1 \
    libxss1 fonts-liberation libappindicator1 libnss3  xdg-utils
```

software-properties-common は、rubyを速くインストールするために入れている。

make と zlib1g-dev は後で bundle install を叩くために必要っぽかったので入れた。使用するgemによっては他にも必要だと思うので、エラーメッセージを見ながら直してほしい。

gconf-service 以下は、google-chromeをdpkgする時に必要な依存パッケージらしい。

```
### install ruby
RUN apt-add-repository ppa:brightbox/ruby-ng && \
    apt-get update && \
    yes | apt-get install ruby2.4 ruby2.4-dev && \
    gem install bundler
```

rubyをインストール。bundlerもついでにインストールしている。

```
### install headless chrome
RUN curl -O https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb && \
    sudo dpkg -i google-chrome-stable_current_amd64.deb && \
    rm -rf google-chrome-stable_current_amd64.deb

### install chromedriver
RUN curl -O https://chromedriver.storage.googleapis.com/2.32/chromedriver_linux64.zip && \
     unzip chromedriver_linux64.zip && \
     chmod +x chromedriver && \
     mv -f chromedriver /usr/local/share/chromedriver && \
     ln -s /usr/local/share/chromedriver /usr/local/bin/chromedriver && \
     ln -s /usr/local/share/chromedriver /usr/bin/chromedriver
```

headless chromeとchromedriverを用意。
こちらを参考にした。

### コンテナに入ってみる

```
### ruby -v
ruby 2.4.1p111 (2017-03-22 revision 58053) [x86_64-linux-gnu]

### chromedriver -v
ChromeDriver 2.32.498513 (2c63aa53b2c658de596ed550eb5267ec5967b351)

### which google-chrome
/usr/bin/google-chrome
```

ちゃんとインストールされている。

### E2Eテストの書き方

こちらを参考に、自分なりに書いてみた。
以下は spec_helper.rb .

```
require ''selenium-webdriver''
require ''capybara''
require ''nokogiri''

Capybara.register_driver :headless_chrome do |app|

  ua = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_4) "
  ua  \
      {binary: ''/usr/bin/google-chrome'', \
      args: [ \
        "--headless",  \
        "--no-sandbox", \
        "--disable-gpu", \
        "--user-agent=#{ua}", \
        "window-size=1280x800"
      ]})

  driver = Capybara::Selenium::Driver.new(
    app,
    browser: :chrome,
    desired_capabilities: caps
  )
end

Capybara.default_driver = :headless_chrome
```

上記をrequireすれば、あとは実コードを書くだけ。
今回はrspecを使わなかったので、こんな感じで書ける。

```
require ''./spec_helper''
require ''test/unit''
```

### 社内プロキシを突破する

社内のプロキシを突破する必要があった。

```
caps = Selenium::WebDriver::Remote::Capabilities.chrome(
    “chromeOptions” => \
      {binary: ’/usr/bin/google-chrome’, \
      args: [ \
        “–headless”,  \
        “–no-sandbox”, \
        “–disable-gpu”, \
        “–user-agent=#{ua}”, \
        “window-size=1280x800”
      ]})
```

の `args` にプロキシとクレデンシャル情報を指定すれば良いと思っていたが、
[公式ドキュメント](https://www.chromium.org/developers/design-documents/network-settings#TOC-Command-line-options-for-proxy-settings)を見ても、
認証は書けないようだった。

### やり方

以下のライブラリを使用した。
[proxy-login-automator](https://github.com/sjitech/proxy-login-automator)


Dockerfileで、任意のディレクトリにこのライブラリををcloneした。(中身はjsファイルひとつなので楽。)
このライブラリ動かすために、Dockerfileに `node` のインストール処理を追加した。
また、Dockerfileの最後で,

```
CMD node /usr/local/bin/proxy-login-automator/proxy-login-automator.js \
    -local_port 18080 \
    -remote_host プロキシのホスト(IP) \
    -remote_port プロキシのポート \
    -usr プロキシのユーザ \
    -pwd プロキシのパスワード
```

を追加した。

これによって、localhost:18080に来たアクセスはプロキシにフォワードされ、認証もきっちりやってくれる、ということになる。(CMDで動いてるので、コンテナ起動時にはこのコマンドが叩かれている。)

あとは、google-chromeからlocalhost:18080にアクセスしてやればいいので、

```
caps = Selenium::WebDriver::Remote::Capabilities.chrome(
    “chromeOptions” => \
      {binary: ’/usr/bin/google-chrome’, \
      args: [ \
        “–headless”,  \
        “–no-sandbox”, \
        “–disable-gpu”, \
        “–user-agent=#{ua}”, \
        “window-size=1280x800”, \
        “–proxy-server=localhost:18080”
      ]})
“`
```

と追加することで、無事通った。
