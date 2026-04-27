#!/usr/bin/env python3
"""
修复 package 名称
"""

import os
import re
from pathlib import Path

HANDLER_DIR = "/Users/zhanglei/MyProjects/mochat/api-server-go/internal/handler"

def fix_package(filepath, expected_package):
    """修复 package 名称"""
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
    print("修复 package 名称...\n")
    
    # 修复 analysis 模块
    print("=== analysis 模块 ===")
    for filepath in Path(os.path.join(HANDLER_DIR, "dashboard/analysis")).glob("*.go"):
        fix_package(str(filepath), "analysis")
    
    # 修复 content 模块
    print("\n=== content 模块 ===")
    for filepath in Path(os.path.join(HANDLER_DIR, "dashboard/content")).glob("*.go"):
        fix_package(str(filepath), "content")
    
    # 修复 marketing 模块
    print("\n=== marketing 模块 ===")
    for filepath in Path(os.path.join(HANDLER_DIR, "dashboard/marketing")).glob("*.go"):
        fix_package(str(filepath), "marketing")
    
    # 修复 platform 模块
    print("\n=== platform 模块 ===")
    for filepath in Path(os.path.join(HANDLER_DIR, "dashboard/platform")).glob("*.go"):
        fix_package(str(filepath), "platform")
    
    # 修复 common 模块
    print("\n=== common 模块 ===")
    for filepath in Path(os.path.join(HANDLER_DIR, "dashboard/common")).glob("*.go"):
        fix_package(str(filepath), "common")
    
    print("\n✅ Package 名称修复完成")

if __name__ == "__main__":
    main()
