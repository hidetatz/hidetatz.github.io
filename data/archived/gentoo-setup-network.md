Gentooのネットワーク周りのセットアップについて---2017-12-26 12:29:07

### 問題
筆者は有線LANの環境がない家に住んでいる(デフォルトで無線LANがあるのが売りのマンション)ので、無線に繋ぐ必要がある。
Gentooのインストールが終わり、ネットワーク設定をしようと思ったが、無線のドライバが認識されない。

具体的には、以下のように `wlp5s0` という無線のドライバが表示されるはずが、

```
$ ip addr
1: lo: # localhost
2: enp4s0: # 有線
3: wlp5s0: # 無線
```

インストール直後は `wlp5s0` が表示されなかった。 `ifconfig -a` しても表示されないので、完全にOSが認識できていないのだろう。と判断した。

* `ifconfig -a` に表示されない以上、カーネルが認識できていないはず
* ということは、カーネルコンフィグが足りてない
* カーネルコンフィグを追加しよう
* どのように追加すればいいか調べる
* カーネルを再ビルドし、認識されるか試す

と考えた。

### 調査方法

`lspci` コマンドで、ドライバについて理解する。

```
% lspci
(略)
05:00.0 Network controller: Qualcomm Atheros QCA9377 802.11ac Wireless Network Adapter (rev 31)
```

この `Qualcomm Atheros QCA9377 802.11ac Wireless Network Adapter` がマシンに組み込まれているだろう、と判断できる。

おもむろにググると、[このような](https://wiki.gentoo.org/wiki/Qualcomm_Atheros_QCA6174) ページにたどり着いた。

型番は若干違うが、これでいけるのでは？と判断し、書いてある通り、

```
[*] Networking support  --->
    [*] Wireless  --->
        <*>   cfg80211 - wireless configuration API
        [ ]     nl80211 testmode command
        [ ]     enable developer warnings
        [ ]     cfg80211 regulatory debugging
        [ ]     cfg80211 certification onus
        [*]     enable powersave by default
        [ ]     cfg80211 DebugFS entries
        [ ]     use statically compiled regulatory rules database
        [ ]     cfg80211 wireless extensions compatibility
        <*>   Generic IEEE 802.11 Networking Stack (mac80211)
        [*]   Minstrel
        [*]     Minstrel 802.11n support
        [ ]       Minstrel 802.11ac support
              Default rate control algorithm (Minstrel)  --->
        [ ]   Enable mac80211 mesh networking (pre-802.11s) support
        -*-   Enable LED triggers
        [ ]   Export mac80211 internals in DebugFS
        [ ]   Trace all mac80211 debug messages
        [ ]   Select mac80211 debugging features  ----
```

```

Device Drivers  --->
    [*] Network device support  --->
        [*]   Wireless LAN  --->
            [*]   Atheros/Qualcomm devices
            <M>     Atheros 802.11ac wireless cards support
            <M>       Atheros ath10k PCI support
```

の設定と、

```
% emerge sys-kernel/linux-firmware
```

をやった。(書いてある通りやっただけ)

これで `make && make modules_install` して、 `make install` して再起動したらイケた。公式は最高。

### 起動時に自動で無線に繋ぐ設定

まずは、[こちら](https://yagi5.com/2017/12/24/wpa_supplicant/)の通り `/etc/wpa_supplicant/wpa_supplicant.conf` を作成する。

ネットワーク設定を登録。

```
% vi /etc/conf.d/net

module_wlp5s0="dhcp"
modules="wpa_supplicant"
wpa_supplicant_wlp5s0="-Dnl80211,wext"
```

起動時に繋ぐよう設定する。

```
% cd /etc/init.d
% ln -s net.lo net.wlp5s0
% /etc/init.d/net.wlp5s0 up
% rc-update add net.wlp5s0 default
```

rebootして、つながるか確認する。
