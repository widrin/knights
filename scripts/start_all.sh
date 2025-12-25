#!/bin/bash

echo "启动所有服务..."

echo "1. 启动中心服务..."
./bin/server -config=configs/center.yaml &
PIDS[0]=$!
sleep 2

echo "2. 启动登录服务..."
./bin/server -config=configs/login.yaml &
PIDS[1]=$!
sleep 1

echo "3. 启动网关服务..."
./bin/server -config=configs/server.yaml &
PIDS[2]=$!
sleep 1

echo "4. 启动游戏服务..."
./bin/server -config=configs/game.yaml &
PIDS[3]=$!

echo "所有服务已启动"
echo "PIDs: ${PIDS[@]}"

# 等待中断信号
trap "echo '正在关闭所有服务...'; kill ${PIDS[@]}; exit" SIGINT SIGTERM

wait
