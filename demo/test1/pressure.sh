#!/bin/bash
PKG=$(basename $(pwd)) # 获取当前路径的最后一个名字，即为文件夹的名字
echo $PKG
while true ; do
    export GOMAXPROCS=$[ 1 + $[ RANDOM % 128 ]] # 随机的GOMAXPROCS
    ./$PKG.test $@ 2>&1 # $@代表可以加入参数 2>&1代表错误输出到控制台
done