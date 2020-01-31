ThinkPadE470にGentooをインストールした その2---2017-12-26 11:39:37


以下について書く。

6 . tarballを展開したディレクトリにchrootする  
7 . portageをセットアップする  
8 . カーネルをビルドする  
9 . ブートローダをセットアップする  
10 . Gentooが起動することを確認する  

### 6. tarballを展開したディレクトリにchrootする

```shell
### /etc/resolv.confのコピー
% cp -L /etc/resolv.conf /mnt/gentoo/etc/

### 必要なファイルシステムをマウント
% mount -t proc proc /mnt/gentoo/proc
% mount --rbind /dev /mnt/gentoo/dev
% mount --rbind /sys /mnt/gentoo/sys

### /dev/sda3に入る
% chroot /mnt/gentoo /bin/bash
% source /etc/profile
```

### 7. portageをセットアップする

portageとは、Gentooのパッケージ管理ツール。MacでいうHomebrew。
設定の前に、コンパイルオプションを設定する。

```shell
% nano -w /etc/portage/make.conf
```

デフォルトに加えて、以下のように編集する。

```
### デフォルトに追加
MAKEOPTS="-march=native"

### 行として追加
GRUB_PLATFORMS="efi-64"
```

Gentooは、基本的に全てのツールをソースコードからビルドする。その際に、ファイルをコンパイルする時のオプションをここで指定できる。と思えばOK。
詳しくは[これ](https://wiki.gentoo.org/wiki/Safe_CFLAGS)を見る。筆者もあまり詳しくない。

そしてportageをセットアップする。

```shell
% emerge-webrsync
% emerge --sync
% eselect profile set 1
```

まあまあ時間がかかった気がする。

### 8. カーネルをビルドする

ここがけっこう大変。
どのようにビルドするかの方針を決める必要がある。
ここは環境に合わせるしかないので、省略する。

### 9. ブートローダをセットアップする

ブートローダをセットアップする。grubを使った。

```shell
% mkdir -p /boot/efi
% emerge grub efibootmgr
% grub-install --target=x86_64-efi --efi-directory=/boot --bootloader-id=gentoo_grub /dev/sda
% mkdir -p /boot/efi/boot
% cp /boot/efi/gentoo_grub/grubx64.efi /boot/efi/boot/bootx64.efi
% grub-mkconfig -o /boot/grub/grub.cfg
```

筆者は、 `grub-install` のときに `Could not prepare Boot variable: Read-only file system` のエラーが出たので、[こちら](https://forums.gentoo.org/viewtopic-t-1069106-start-0.html)を参考に以下のようにした。

```shell
% mount -o remount,rw /sys/firmware/efi/efivars
```

また、 `grub-mkconfig` の時のメッセージで、先程ビルドしたカーネルが認識されていることを確認すること。

### 10. Gentooが起動することを確認する

これでたぶんできているはず。
以下のようにして再起動。

```
% exit

### chrootは抜けている
% cd
% umount -l /mnt/gentoo/dev{/shm,/pts,}
% umount -l /mnt/gentoo{/boot,/proc,}
% reboot
```

Gentooが上がってきたら優勝です。
