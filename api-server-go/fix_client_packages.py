#!/usr/bin/env python3
"""
修复 client 目录的 package 名称
"""

import os
import re
from pathlib import Path

HANDLER_DIR = "/Users/zhanglei/MyProjects/mochat/api-server-go/internal/handler/client"

def fix_client_package(filepath, expected_package):
    """修复 client 子模块的 package 名称"""
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 查找 package 声明
    match = re.search(r'^package (\w+)$', content, re.MULTILINE)
    if not match:
        print(f"⚠ {filepath}: 未找到 package 声明")
        return
    
    current_package = match.group(1)
    if current_package == expected_package:
        print(f"✓ {filepath}: package 已正确 ({current_package})")
        return
    
    # 替换 package
    content = re.sub(r'^package \w+$', f'package {expected_package}', content, count=1)
    
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"✅ {filepath}: {current_package} -> {expected_package}")

def main():
    print("修复 client 目录的 package 名称...\n")
    
    # 各子模块对应的 package 名称
    modules = {
        'contact': 'contact',
        'organization': 'organization',
        'platform': 'platform',
        'content': 'content',
        'common': 'common',
    }
    
    for module_name, package_name in modules.items():
        print(f"=== {module_name} 模块 ===")
        module_dir = os.path.join(HANDLER_DIR, module_name)
        if os.path.exists(module_dir):
            for filepath in Path(module_dir).glob("*.go"):
                fix_client_package(str(filepath), package_name)
    
    print("\n✅ Client 目录的 package 名称修复完成")

if __name__ == "__main__":
    main()
