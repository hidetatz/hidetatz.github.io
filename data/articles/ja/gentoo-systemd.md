GentooをSystemdに移行した---2018-02-26 20:11:32

OpenRCからSystemdに移行した。

ウインドウマネージャにFluxboxを使っていたが、GNOMEが使ってみたくなり、Systemdを入れることにした。
オーバーレイを導入すればOpenRCでもGNOMEが使えるようだったが、そんなにOpenRCにこだわりがなかったし、
最近カーネルのビルドもやってなかったので、やってみた。

### カーネル設定

カーネルのリビルドが必要だったので、以下のように対応。ついでにアップデートした。

```
$ emerge gentoo-sources

$ eselect kernel list
### eselect kernel set n でインストールしたカーネルにセット

$ cd /usr/src/linux
$ cp /usr/src/前のkernelのディレクトリ/.config ./.config
$ make oldconfig
```

以下のように設定。ソースは[公式](https://wiki.gentoo.org/wiki/Systemd)。

```
Gentoo Linux --->
   Support for init systems, system and service managers --->
      [*] systemd

General setup  --->
	[*] open by fhandle syscalls
	[*] Control Group support --->
	[ ] Enable deprecated sysfs features to support old userspace tools
	[*] Configure standard kernel features (expert users)  --->
		[*] Enable eventpoll support
		[*] Enable signalfd() system call
		[*] Enable timerfd() system call
[*] Networking support --->
Device Drivers  --->
	Generic Driver Options  --->
		[*] Maintain a devtmpfs filesystem to mount at /dev
File systems  --->
	[*] Inotify support for userspace
	Pseudo filesystems  --->
		[*] /proc file system support
		[*] sysfs file system support

General setup  --->
        [*] Checkpoint/restore support
	[*] Namespaces support  --->
		[*] Network namespace
[*] Enable the block layer  --->
	[*] Block layer SG support v4
Processor type and features  --->
	[*] Enable seccomp to safely compute untrusted bytecode
Networking support --->
	Networking options --->
		<*> The IPv6 protocol
Device Drivers  --->
	Generic Driver Options  --->
		()  path to uevent helper
		[ ] Fallback user-helper invocation for firmware loading
Firmware Drivers  --->
	[*] Export DMI identification via sysfs to userspace
File systems --->
	<*> Kernel automounter version 4 support (also supports v3)
	Pseudo filesystems --->
		[*] Tmpfs virtual memory file system support (former shm fs)
		[*]   Tmpfs POSIX Access Control Lists
		[*]   Tmpfs extended attributes
```

これでカーネルをリビルドする。

```
$ make && make modules_install
$ make install
$ mount /boot
$ grub-mkconfig -o /boot/grub/grub.cfg

$ vi /etc/default/grub
### 以下の行のコメントアウトを外す
GRUB_CMDLINE_LINUX="init=/lib/systemd/systemd"
```

再起動する前に、公式に従って以下を対応。

```
$ ln -sf /proc/self/mounts /etc/mtab
$ eselect profile list
### systemdのものにset

$ emerge --unmerge sys-fs/udev
$ emerge -avDN @world
```

あとは再起動する。

```
$ sudo reboot
```

これでうまく立ち上がるのを祈る。

### systemctl 設定

今までrc-serviceに登録してたものをsystemctlに登録する。

```
$ systemctl enable cronie.service
$ systemctl enable sshd.service
$ systemctl enable ntpd.service
$ systemctl enable wpa_supplicant@wlp5s0
$ systemctl enable docker.service
```

再起動して、ちゃんとサービスが上がるか、インターネットにつながるか確認する。
