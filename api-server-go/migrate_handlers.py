#!/usr/bin/env python3
"""
Handler 目录重构迁移脚本
自动更新 package 名称和 import 路径
"""

import os
import re
from pathlib import Path

PROJECT_ROOT = "/Users/zhanglei/MyProjects/mochat/api-server-go"
HANDLER_DIR = os.path.join(PROJECT_ROOT, "internal", "handler")

def update_file_package(filepath, new_package):
    """更新文件的 package 声明"""
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 替换 package 声明
    content = re.sub(r'^package \w+$', f'package {new_package}', content, count=1)
    
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"✓ {filepath}: package -> {new_package}")

def update_file_imports(filepath, old_pattern, new_pattern):
    """更新文件的 import 路径"""
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 替换 import 路径
    content = re.sub(old_pattern, new_pattern, content)
    
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"✓ {filepath}: import updated")

def migrate_module(module_name, package_name, old_import_pattern):
    """迁移一个模块的所有文件"""
    module_dir = os.path.join(HANDLER_DIR, module_name)
    if not os.path.exists(module_dir):
        print(f"⚠ {module_dir} 不存在，跳过")
        return
    
    print(f"\n=== 迁移模块: {module_name} ===")
    
    for filepath in Path(module_dir).glob("*.go"):
        update_file_package(str(filepath), package_name)
        if old_import_pattern:
            new_pattern = f"mochat-api-server/internal/handler/{module_name}"
            update_file_imports(str(filepath), old_import_pattern, new_pattern)

def main():
    print("开始 Handler 目录重构迁移...")
    
    # 迁移 system 模块
    migrate_module(
        "dashboard/system",
        "system",
        r"mochat-api-server/internal/handler/dashboard"
    )
    
    # 迁移 contact 模块
    migrate_module(
        "dashboard/contact",
        "contact",
        r"mochat-api-server/internal/handler/dashboard"
    )
    
    # 迁移 organization 模块
    migrate_module(
        "dashboard/organization",
        "organization",
        r"mochat-api-server/internal/handler/dashboard"
    )
    
    print("\n" + "="*50)
    print("✅ 基础模块迁移完成！")
    print("\n下一步：")
    print("1. 检查编译错误")
    print("2. 继续迁移其他模块（content、marketing、analysis、platform）")
    print("3. 迁移 sidebar 到 client")
    print("4. 更新 router 引用路径")

if __name__ == "__main__":
    main()
