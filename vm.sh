#!/bin/bash

set -e

msg() {
	echo "$@" >&2
}
die() {
	msg "$@"
	exit 1
}
ensurecmd() {
	if ! which "$1" >/dev/null 2>&1; then
		die missing command $1
	fi
}

ensurecmd cpio
ensurecmd mkfs.ext4
ensurecmd mkisofs
ensurecmd qemu-system-x86_64
ensurecmd bsdtar
ensurecmd xorriso

iso_url=https://dl-cdn.alpinelinux.org/alpine/v3.20/releases/x86_64/alpine-virt-3.20.2-x86_64.iso
iso_name=alpine.iso

qemu_args=

disk_name=disk.img
disk_size=10M

# download an alpine image file
if [ ! -e "$iso_name" ]; then
	msg downloading disk image...
	curl "$iso_url" -o "$iso_name"
	xorriso -dev "$iso_name" -volid boot -commit >/dev/null 2>&1
	bsdtar -xf "$iso_name" boot
	chmod u+w -R boot
fi

# create iso from the custom/ directory
# (because squashfs does not support labels and
# erofs was not built into the default alpine kernel)
if [ -d custom/ ]; then
	qemu_args+="-drive file=custom.iso,format=raw,index=1"

	if [ ! -e custom.iso ]; then
		msg creating custom.iso...
		mkisofs -R -V root -o custom.iso custom/ >/dev/null 2>&1
	fi
fi

swap=$(mktemp)
trap 'rm "$swap"' EXIT
dd if=/dev/zero of="$swap" bs=1M count=1 status=none >/dev/null
mkswap -L swap "$swap" >/dev/null

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

if [ ! -f "$disk_name" ]; then
	msg creating a disk for permanent files...
	truncate -s "$disk_size" "$disk_name"
	mkdir "$tmp"/diskroot
	trap 'rm -rf "$tmp"/diskroot' EXIT
	mkdir "$tmp"/diskroot/enabled
	ln -sf /etc/inwit/syslogd.service "$tmp"/diskroot/enabled/syslogd
	# TODO: use guestfish?
	# TODO: this is so slow
	virt-make-fs --label=disk \
		--type=ext4 "$tmp"/diskroot \
		"$disk_name"
fi

# replace /init with ours
msg unpacking initramfs to replace /init...
zcat boot/initramfs-virt | cpio --quiet -idD "$tmp"
install -Dm755 ./init "$tmp"/init
install -Dm755 ./initctl "$tmp"/sbin/initctl
( cd "$tmp"; find -print0 | sort -z | \
	cpio -o0 --quiet -H newc | gzip
) > boot/initramfs-virt

# boot it without passing through grub
msg booting...
qemu-system-x86_64 \
	-enable-kvm \
	-m 512M \
	-smp cores=2 \
	-kernel boot/vmlinuz-virt \
	-initrd boot/initramfs-virt \
	-drive file="$iso_name",format=raw,index=0 \
	-drive file="$swap",format=raw,index=2 \
	-drive file="$disk_name",format=raw,index=3 \
	-append "console=ttyS0" \
	-nographic \
	${qemu_args[@]}
