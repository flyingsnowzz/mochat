#!/bin/bash

# MoChat 系统启动脚本
# 用于启动整个MoChat系统，包括后端API服务和前端项目

echo "===================================="
echo "MoChat 系统启动脚本"
echo "===================================="

# 颜色定义
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
RED="\033[0;31m"
NC="\033[0m" # No Color

# 检查命令是否存在
check_command() {
    command -v "$1" >/dev/null 2>&1 || {
        echo -e "${RED}错误: 命令 $1 不存在，请安装${NC}"
        exit 1
    }
}

# 检查环境
check_environment() {
    echo -e "${YELLOW}检查环境...${NC}"
    
    # 检查PHP
    check_command php
    PHP_VERSION=$(php -v | head -n 1 | awk '{print $2}')
    echo -e "${GREEN}PHP 版本: $PHP_VERSION${NC}"
    
    # 检查Composer
    check_command composer
    echo -e "${GREEN}Composer 已安装${NC}"
    
    # 检查Node.js
    check_command node
    NODE_VERSION=$(node -v)
    echo -e "${GREEN}Node.js 版本: $NODE_VERSION${NC}"
    
    # 检查Yarn
    check_command yarn
    echo -e "${GREEN}Yarn 已安装${NC}"
    
    echo -e "${GREEN}环境检查完成${NC}"
}

# 启动后端API服务
start_backend() {
    echo -e "${YELLOW}启动后端API服务...${NC}"
    
    cd "$(dirname "$0")/api-server" || {
        echo -e "${RED}错误: 无法进入 api-server 目录${NC}"
        exit 1
    }
    
    # 检查依赖
    if [ ! -d "vendor" ]; then
        echo -e "${YELLOW}安装后端依赖...${NC}"
        composer install
        if [ $? -ne 0 ]; then
            echo -e "${RED}错误: 安装后端依赖失败${NC}"
            exit 1
        fi
    fi
    
    # 启动服务
    echo -e "${YELLOW}启动PHP服务...${NC}"
    php bin/hyperf.php start -d
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}后端API服务启动成功${NC}"
    else
        echo -e "${RED}错误: 后端API服务启动失败${NC}"
        exit 1
    fi
    
    cd ..
}

# 启动前端dashboard服务
start_dashboard() {
    echo -e "${YELLOW}启动前端dashboard服务...${NC}"
    
    cd "$(dirname "$0")/dashboard" || {
        echo -e "${RED}错误: 无法进入 dashboard 目录${NC}"
        exit 1
    }
    
    # 检查依赖
    if [ ! -d "node_modules" ]; then
        echo -e "${YELLOW}安装前端依赖...${NC}"
        yarn install
        if [ $? -ne 0 ]; then
            echo -e "${RED}错误: 安装前端依赖失败${NC}"
            exit 1
        fi
    fi
    
    # 启动开发服务器
    echo -e "${YELLOW}启动前端开发服务器...${NC}"
    yarn run dev &
    
    echo -e "${GREEN}前端dashboard服务启动成功${NC}"
    cd ..
}

# 启动前端sidebar服务
start_sidebar() {
    echo -e "${YELLOW}启动前端sidebar服务...${NC}"
    
    cd "$(dirname "$0")/sidebar" || {
        echo -e "${RED}错误: 无法进入 sidebar 目录${NC}"
        exit 1
    }
    
    # 检查依赖
    if [ ! -d "node_modules" ]; then
        echo -e "${YELLOW}安装前端依赖...${NC}"
        yarn install
        if [ $? -ne 0 ]; then
            echo -e "${RED}错误: 安装前端依赖失败${NC}"
            exit 1
        fi
    fi
    
    # 启动开发服务器
    echo -e "${YELLOW}启动前端开发服务器...${NC}"
    yarn run dev &
    
    echo -e "${GREEN}前端sidebar服务启动成功${NC}"
    cd ..
}

# 启动前端operation服务
start_operation() {
    echo -e "${YELLOW}启动前端operation服务...${NC}"
    
    cd "$(dirname "$0")/operation" || {
        echo -e "${RED}错误: 无法进入 operation 目录${NC}"
        exit 1
    }
    
    # 检查依赖
    if [ ! -d "node_modules" ]; then
        echo -e "${YELLOW}安装前端依赖...${NC}"
        yarn install
        if [ $? -ne 0 ]; then
            echo -e "${RED}错误: 安装前端依赖失败${NC}"
            exit 1
        fi
    fi
    
    # 启动开发服务器
    echo -e "${YELLOW}启动前端开发服务器...${NC}"
    yarn run dev &
    
    echo -e "${GREEN}前端operation服务启动成功${NC}"
    cd ..
}

# 显示启动信息
show_info() {
    echo -e "${GREEN}====================================${NC}"
    echo -e "${GREEN}MoChat 系统启动完成${NC}"
    echo -e "${GREEN}====================================${NC}"
    echo -e "${YELLOW}后端API服务: http://localhost:9501${NC}"
    echo -e "${YELLOW}前端dashboard: http://localhost:8080${NC}"
    echo -e "${YELLOW}前端sidebar: http://localhost:8081${NC}"
    echo -e "${YELLOW}前端operation: http://localhost:8082${NC}"
    echo -e "${GREEN}====================================${NC}"
    echo -e "${YELLOW}提示: 按 Ctrl+C 停止所有服务${NC}"
}

# 主函数
main() {
    # 检查环境
    check_environment
    
    # 启动后端服务
    start_backend
    
    # 启动前端服务
    start_dashboard
    start_sidebar
    start_operation
    
    # 显示启动信息
    show_info
    
    # 等待用户输入
    read -p "按 Enter 键停止所有服务..."
    
    # 停止所有服务
    echo -e "${YELLOW}停止所有服务...${NC}"
    pkill -f "php bin/hyperf.php"
    pkill -f "yarn run dev"
    echo -e "${GREEN}所有服务已停止${NC}"
}

# 执行主函数
main
