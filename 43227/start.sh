#!/bin/bash

echo "=========================================="
echo "  养蜂管理与蜂蜜交易平台 - 启动脚本"
echo "=========================================="

echo ""
echo "[1/3] 检查 Go 环境..."
if ! command -v go &> /dev/null; then
    echo "错误: 未检测到 Go 环境，请先安装 Go 1.21+"
    exit 1
fi
echo "Go 版本: $(go version)"

echo ""
echo "[2/3] 检查 Node.js 环境..."
if ! command -v node &> /dev/null; then
    echo "错误: 未检测到 Node.js 环境，请先安装 Node.js 18+"
    exit 1
fi
echo "Node.js 版本: $(node --version)"
echo "npm 版本: $(npm --version)"

echo ""
echo "[3/3] 检查 PostgreSQL..."
if ! command -v psql &> /dev/null; then
    echo "警告: 未检测到 PostgreSQL，请确保数据库服务已启动"
else
    echo "PostgreSQL 客户端: $(psql --version)"
fi

echo ""
echo "=========================================="
echo "  后端服务启动"
echo "=========================================="
cd backend

echo "安装 Go 依赖..."
go mod download

echo "启动后端服务 (端口: 8080)..."
go run main.go &
BACKEND_PID=$!
echo "后端服务已启动，PID: $BACKEND_PID"

cd ..

echo ""
echo "=========================================="
echo "  前端服务启动"
echo "=========================================="
cd frontend

echo "安装前端依赖..."
npm install

echo "启动前端开发服务器 (端口: 3000)..."
npm run dev &
FRONTEND_PID=$!
echo "前端服务已启动，PID: $FRONTEND_PID"

cd ..

echo ""
echo "=========================================="
echo "  服务启动完成"
echo "=========================================="
echo ""
echo "后端 API: http://localhost:8080"
echo "前端页面: http://localhost:3000"
echo ""
echo "按 Ctrl+C 停止所有服务"

trap "kill $BACKEND_PID $FRONTEND_PID 2>/dev/null" EXIT

wait
