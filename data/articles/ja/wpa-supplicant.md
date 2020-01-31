wpa_supplicantでWPAの無線LANに繋ぐ---2017-12-24 19:15:39

### 暗号化方式を確かめる
Windowsでログインして、繋がってるネットワークのプロパティからWEPなのか、WPA系なのかを確かめる。
筆者環境の無線LANはWPA2形式で暗号化されていた。
以下の方法はWPA系に対応した方法。

### wpa_supplicantで繋ぐ

筆者はGentooインストールをSystemRescueCDからやっていたのでwpa_supplicantコマンドはデフォルトで入っていた。
GentooのminimalCDとかには入ってないかもなので、その場合はSystemRescueCDを使うか、一回有線でつないでemergeするなどすること。

まずは設定ファイルを作る。

```
% wpa_passphrase "各自の無線のSSID" >> /etc/wpa_supplicant/wpa_supplicant.conf

```

これで、 `/etc/wpa_supplicant/wpa_supplicant.conf` に設定ファイルができる。
念のため閲覧権限を与えておく。

```
% chmod -v 600 /etc/wpa_supplicant/wpa_supplicant.conf
```

次に、無線LANをキャッチできるドライバの名前を確かめる。

```
% iwconfig
```

iwconfigで、 `wlp...` みたいな名前のドライバがあると思うので、それの名前を覚えておく。
筆者のマシン(ThinkPad E470) では `wlp5s0` だったので、その名前で説明する。

そして、実際に繋ぐ。

```
% wpa_supplicant -Dnl80211,wext -iwlp5s0 -c/etc/wpa_supplicant/wpa_supplicant.conf -B
Successfully initialized wpa_supplicant
```

そして、DHCPからIPアドレスをもらう。

```
% dhcpcd wlp5s0
```

確認する。

```
% ifconfig wlp5s0
```

これで、 IPアドレスがちゃんと振られていれば、インターネットに繋がると思う。
pingを打つなりして確認する。

### つながらない場合
もしつながらない場合はwpa_supplicantのログを出しながら繋ぐ。

```
% wpa_supplicant -Dnl80211,wext -iwlp5s0 -c/etc/wpa_supplicant/wpa_supplicant.conf
```

-Bオプションをなくすと、その場で繋いで、ログを出してくれるのでそれを読む。

また、

```
% killall wpa_supplicant
```

で一回殺したり、

```
% ifconfig wlp5s0 down
% ifconfig wlp5s0 up
```

で、ドライバをupしなおしたりするといいという情報もある。
