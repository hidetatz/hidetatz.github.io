grubでGentooとWindowsを起動できるようにする---2017-12-28 22:36:16

grubを使ってLinuxとWindowsを選択できるようにする方法を書く。

[ハンドブック](https://wiki.gentoo.org/wiki/Handbook:AMD64/Installation/Bootloader/ja)にはLinuxとWindowsのデュアルブートのやり方があまり書いてない。
このやり方はマザーボードがUEFIに対応していることを前提としている。

### 必要なライブラリ

```shell
### emerge --ask --newuse sys-boot/os-prober
```

### udevをマウント

`udev` をマウントすることで別のOSの情報にアクセスできるようになるらしい。
(これは `/mnt/gentoo` にgentooのルートをマウントしていることを前提としている。)

```shell
### mkdir -p /mnt/gentoo/run/udev
### mount -o bind /run/udev /mnt/gentoo/run/udev
### mount --make-rslave /mnt/gentoo/run/udev
```

### grubの設定ファイルを作る

```shell
### grub-mkconfig -o /boot/grub/grub.cfg
```

これで標準出力に `Found windows boot manager...` みたいに出てきたらOK。
rebootして、windowsが選択、起動できることをテストする。

### 参考
[GRUB2](https://wiki.gentoo.org/wiki/GRUB2/ja)
