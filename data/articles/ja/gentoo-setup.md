Gentooの初期設定---2017-12-26 12:05:00

Gentooの初期設定をメモる。

ネットワーク周りについては、色々ハマったので[こちら](https://yagi5.com/2017/12/26/gentoo-setup-network/)に書いた。

### ユーザ作成

```
% useradd -m -G users,wheel,audio -s /bin/bash ユーザー名
% passwd ユーザー名
```

### sudo導入

sudoがないとなにもできない。

```
% emerge sudo

### 設定
% visudo

### ここのコメントアウトを外す
%wheel   ALL=(ALL)   NOPASSWD: ALL
```

これにより、wheelグループに入っているユーザはsudoできるようになる。

### sshdを導入

sshdを入れて、Macからsshで入って作業できるようにする。(色々と便利なので。)
sshd自体はデフォルトで入っていた。

```
### 起動時にsshdをあげる
% rc-update add sshd default
% rc-service sshd start

### 作っておく
% mkdir ~/.ssh
```

まずはMac側で作業。

```
$ cd ~/.ssh
$ ssh-keygen -t rsa -f gentoo_rsa
$ scp ./gentoo_rsa.pub ユーザー名@gentooホストのIP:~/.ssh
```

Gentoo側で作業

```
$ cd ~
$ chmod 700 .ssh
$ cd .ssh
$ cat gentoo_rsa.pub >> authorized_keys
$ chmod 600 authorized_keys
$ rm -fv gentoo_rsa.pub
```

```
$ sudo vi /etc/ssh/sshd_config
```

以下のようにする。

```
PubkeyAuthentication yes
AuthorizedKeysFile     .ssh/authorized_keys
PermitRootLogin no
PasswordAuthentication no
```

sshd再起動

```
$ rc-service sshd restart
```

Mac側でエイリアスを設定する。

```
$ vi ~/.ssh/config
```

```
Host gentoo
  User ユーザ名
  HostName IP
  IdentityFile ~/.ssh/gentoo_rsa
```

これで、Macから `ssh gentoo` で入れるようになる。

### NTP導入

ntpdをあげておく。

```
$ rc-service ntpd start
$ rc-update add ntpd default
```

### zsh導入

zshにしておく。

```
$ emerge --ask app-shells/zsh
$ emerge --ask app-shells/zsh-completions
```

初期設定

```
### .zshrcはお好みで。
$ touch ~/.zshrc
$ chsh -s /bin/zsh
```

### git導入

```
emerge --ask dev-vcs/git
```
