#!/bin/bash
set -ex

mkdir /chroot
mkdir /chroot/bin
mkdir /chroot/lib
mkdir /chroot/lib64
mkdir /chroot/dev
mkdir /chroot/tmp
mkdir /chroot/var

# # debootstrap
# debootstrap saucy /chroot

# busybox
cp /bin/busybox /chroot/bin/sh
cp /lib64/ld-linux-x86-64.so.2 /chroot/lib64/ld-linux-x86-64.so.2
cp /lib/x86_64-linux-gnu/libc.so.6 /chroot/lib/libc.so.6

# legacy-bridge
cp /src/sandstorm-master/bin/legacy-bridge /chroot/
cp /usr/local/lib/libcapnp-rpc-0.5-dev.so /chroot/lib/libcapnp-rpc-0.5-dev.so
cp /usr/local/lib/libkj-async-0.5-dev.so /chroot/lib/libkj-async-0.5-dev.so
cp /usr/local/lib/libcapnp-0.5-dev.so /chroot/lib/libcapnp-0.5-dev.so
cp /usr/local/lib/libkj-0.5-dev.so /chroot/lib/libkj-0.5-dev.so
cp /usr/lib/x86_64-linux-gnu/libstdc++.so.6 /chroot/lib/libstdc++.so.6
cp /lib/x86_64-linux-gnu/libm.so.6 /chroot/lib/libm.so.6
cp /lib/x86_64-linux-gnu/libgcc_s.so.1 /chroot/lib/libgcc_s.so.1

# shell
go build -o /chroot/shell github.com/kevinwallace/sandstorm-shell/shell
cp /lib/x86_64-linux-gnu/libpthread.so.0 /chroot/lib/libpthread.so.0

# manifest
capnp eval -I /src/sandstorm-master/src -b /root/manifest.capnp manifest > /chroot/sandstorm-manifest

# package
spk pack /chroot /root/secret.key /output/shell.spk