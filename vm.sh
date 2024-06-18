#!/bin/sh

set -e

msg() {
	echo "$@" >&2
}

iso_url=https://dl-cdn.alpinelinux.org/alpine/v3.20/releases/x86_64/alpine-standard-3.20.0-x86_64.iso
iso_name=alpine.iso

# download an alpine image file
if [ ! -e "$iso_name" ]; then
	msg downloading disk image...
	curl "$iso_url" -o "$iso_name"
	bsdtar -xf "$iso_name" boot
	chmod u+w -R boot
fi

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

# replace /init with ours
msg unpacking initramfs to replace /init...
zcat boot/initramfs-lts | cpio --quiet -idD "$tmp"
install -Dm755 ./init "$tmp"/init
( cd "$tmp"; find -print0 | sort -z | \
	cpio -o0 --quiet -H newc | gzip
) > boot/initramfs-lts

# boot it without passing through grub
msg booting...
qemu-system-x86_64 \
	-enable-kvm \
	-m 512M \
	-smp cores=2 \
	-kernel boot/vmlinuz-lts \
	-initrd boot/initramfs-lts \
	-drive file="$iso_name",format=raw \
	-append "console=ttyS0" \
	-nographic
