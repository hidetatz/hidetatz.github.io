GentooのWifiがブツブツ切れる問題に対応した---2018-01-14 22:46:49

メモっておく。

何故か無線LANが12時間程度で切れて、再起動しないと二度とつながらないことが多く、色々ググッた結果、割とあっさり解決した。

### やったこと

パワーマネジメント機能をオフにした。

```
$ iwconfig
(前略)

wlp5s0    IEEE 802.11  ESSID:"My-SSID"
          Mode:Managed  Frequency:2.427 GHz  Access Point: MyAccessPoint
          Bit Rate=1 Mb/s   Tx-Power=20 dBm
          Retry short limit:7   RTS thr:off   Fragment thr:off
          Power Management:on # これ
          Link Quality=53/70  Signal level=-57 dBm
          Rx invalid nwid:0  Rx invalid crypt:0  Rx invalid frag:0
          Tx excessive retries:1  Invalid misc:498   Missed beacon:0

(後略)
```

これをオフにした。

```
$ iwconfig wlp5s0 power off
```

これをやってから快適になった。
要するに消費電力が多くなったりするらしいが、まあいいかなと思っている。
