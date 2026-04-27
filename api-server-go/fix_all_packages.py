#!/usr/bin/env python3
"""
彻底修复所有 package 名称
"""

import os
import re
from pathlib import Path

HANDLER_DIR = "/Users/zhanglei/MyProjects/mochat/api-server-go/internal/handler"

# 定义各目录应该的 package 名称
PACKAGE_MAPPING = {
    "dashboard/analysis": "analysis",
    "dashboard/common": "common",
    "dashboard/contact": "contact",
    "dashboard/content": "content",
    "dashboard/marketing": "marketing",
    "dashboard/organization": "organization",
    "dashboard/platform": "platform",
    "dashboard/system": "system",
    "client/contact": "contact",
    "client/content": "content",
    "client/organization": "organization",
    "client/platform": "platform",
    "client/common": "common",
}

def fix_all_packages():
    """修复所有 package 名称"""
    print("彻底修复所有 package 名称...\n")

    for dir_path, expected_package in PACKAGE_MAPPING.items():
        full_path = os.path.join(HANDLER_DIR, dir_path)
        if not os.path.exists(full_path):
            continue

        print(f"=== {dir_path} (期望: {expected_package}) ===")
        for filepath in Path(full_path).glob("*.go"):
            with open(filepath, 'r', encoding='utf-8') as f:
                lines = f.readlines()

            if not lines:
                continue

            # 查找第一个 package 声明
            for i, line in enumerate(lines):
                if line.strip().startswith("package "):
                    current_package = line.strip().replace("package ", "")
                    if current_package == expected_package:
                        print(f"  ✓ {os.path.basename(filepath)}: {current_package} (已正确)")
                    else:
                        lines[i] = f"package {expected_package}\n"
                        with open(filepath, 'w', encoding='utf-8') as f:
                            f.writelines(lines)
                        print(f"  ✅ {os.path.basename(filepath)}: {current_package} -> {expected_package}")
                    break

    print("\n✅ 所有 package 名称修复完成")

if __name__ == "__main__":
    fix_all_packages()
