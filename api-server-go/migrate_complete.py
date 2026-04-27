#!/usr/bin/env python3
"""
Handler 目录重构迁移脚本 - 完整版
"""

import os
import re
import shutil
from pathlib import Path

PROJECT_ROOT = "/Users/zhanglei/MyProjects/mochat/api-server-go"
HANDLER_DIR = os.path.join(PROJECT_ROOT, "internal", "handler")

def copy_file(src, dst):
    """复制文件"""
    shutil.copy2(src, dst)
    print(f"✓ 复制: {os.path.basename(src)} -> {dst}")

def update_file_package(filepath, new_package):
    """更新文件的 package 声明"""
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    content = re.sub(r'^package \w+$', f'package {new_package}', content, count=1)
    
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"✓ package -> {new_package}")

def update_file_imports(filepath, old_pattern, new_pattern):
    """更新文件的 import 路径"""
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    content = re.sub(old_pattern, new_pattern, content)
    
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"✓ import 更新完成")

def migrate_content_module():
    """迁移 content 模块"""
    print("\n=== 迁移 content 模块 ===")
    
    # 创建目录
    os.makedirs(os.path.join(HANDLER_DIR, "dashboard/content"), exist_ok=True)
    
    # 复制文件
    copy_file(
        os.path.join(HANDLER_DIR, "dashboard/medium.go"),
        os.path.join(HANDLER_DIR, "dashboard/content/medium.go")
    )
    copy_file(
        os.path.join(HANDLER_DIR, "dashboard/plugin/greeting_handler.go"),
        os.path.join(HANDLER_DIR, "dashboard/content/greeting.go")
    )
    
    # 更新 package 和 import
    for filepath in Path(os.path.join(HANDLER_DIR, "dashboard/content")).glob("*.go"):
        update_file_package(str(filepath), "content")
        # 替换 dashboard/plugin 引用
        update_file_imports(str(filepath), 
                          r"mochat-api-server/internal/handler/dashboard/plugin",
                          "mochat-api-server/internal/handler/dashboard/content")
        # 替换 dashboard 引用
        update_file_imports(str(filepath),
                          r"mochat-api-server/internal/handler/dashboard[^/]",
                          "mochat-api-server/internal/handler/dashboard/content")

def migrate_marketing_module():
    """迁移 marketing 模块"""
    print("\n=== 迁移 marketing 模块 ===")
    
    os.makedirs(os.path.join(HANDLER_DIR, "dashboard/marketing"), exist_ok=True)
    
    # 复制文件
    copy_file(
        os.path.join(HANDLER_DIR, "dashboard/plugin/channel_code_handler.go"),
        os.path.join(HANDLER_DIR, "dashboard/marketing/channel_code.go")
    )
    copy_file(
        os.path.join(HANDLER_DIR, "dashboard/plugin/room_welcome_handler.go"),
        os.path.join(HANDLER_DIR, "dashboard/marketing/room_welcome.go")
    )
    copy_file(
        os.path.join(HANDLER_DIR, "dashboard/plugin/work_room_auto_pull_handler.go"),
        os.path.join(HANDLER_DIR, "dashboard/marketing/auto_pull.go")
    )
    copy_file(
        os.path.join(HANDLER_DIR, "dashboard/plugin/room_tag_pull_handler.go"),
        os.path.join(HANDLER_DIR, "dashboard/marketing/room_tag_pull.go")
    )
    
    # 更新 package 和 import
    for filepath in Path(os.path.join(HANDLER_DIR, "dashboard/marketing")).glob("*.go"):
        update_file_package(str(filepath), "marketing")
        update_file_imports(str(filepath),
                          r"mochat-api-server/internal/handler/dashboard/plugin",
                          "mochat-api-server/internal/handler/dashboard/marketing")

def migrate_analysis_module():
    """迁移 analysis 模块"""
    print("\n=== 迁移 analysis 模块 ===")
    
    os.makedirs(os.path.join(HANDLER_DIR, "dashboard/analysis"), exist_ok=True)
    
    copy_file(
        os.path.join(HANDLER_DIR, "dashboard/index_handler.go"),
        os.path.join(HANDLER_DIR, "dashboard/analysis/index.go")
    )
    copy_file(
        os.path.join(HANDLER_DIR, "dashboard/plugin/statistic_handler.go"),
        os.path.join(HANDLER_DIR, "dashboard/analysis/statistic.go")
    )
    
    for filepath in Path(os.path.join(HANDLER_DIR, "dashboard/analysis")).glob("*.go"):
        update_file_package(str(filepath), "analysis")
        update_file_imports(str(filepath),
                          r"mochat-api-server/internal/handler/dashboard[^/]",
                          "mochat-api-server/internal/handler/dashboard/analysis")

def migrate_platform_module():
    """迁移 platform 模块"""
    print("\n=== 迁移 platform 模块 ===")
    
    os.makedirs(os.path.join(HANDLER_DIR, "dashboard/platform"), exist_ok=True)
    
    # official_account_handler.go 存在
    if os.path.exists(os.path.join(HANDLER_DIR, "dashboard/official_account_handler.go")):
        copy_file(
            os.path.join(HANDLER_DIR, "dashboard/official_account_handler.go"),
            os.path.join(HANDLER_DIR, "dashboard/platform/official_account.go")
        )
    
    # agent.go 可能不存在，需要从 sidebar 迁移或创建空文件
    # 注意：agent 会在 sidebar -> client 迁移时处理
    
    for filepath in Path(os.path.join(HANDLER_DIR, "dashboard/platform")).glob("*.go"):
        update_file_package(str(filepath), "platform")
        update_file_imports(str(filepath),
                          r"mochat-api-server/internal/handler/dashboard[^/]",
                          "mochat-api-server/internal/handler/dashboard/platform")

def migrate_common_module():
    """迁移 common 模块"""
    print("\n=== 迁移 common 公共接口 ===")
    
    os.makedirs(os.path.join(HANDLER_DIR, "dashboard/common"), exist_ok=True)
    
    if os.path.exists(os.path.join(HANDLER_DIR, "dashboard/common.go")):
        copy_file(
            os.path.join(HANDLER_DIR, "dashboard/common.go"),
            os.path.join(HANDLER_DIR, "dashboard/common/common.go")
        )
        
        for filepath in Path(os.path.join(HANDLER_DIR, "dashboard/common")).glob("*.go"):
            update_file_imports(str(filepath),
                              r"mochat-api-server/internal/handler/dashboard[^/]",
                              "mochat-api-server/internal/handler/dashboard/common")

def migrate_sidebar_to_client():
    """迁移 sidebar 到 client"""
    print("\n=== 迁移 sidebar 到 client ===")
    
    # 创建 client 子目录
    os.makedirs(os.path.join(HANDLER_DIR, "client/contact"), exist_ok=True)
    os.makedirs(os.path.join(HANDLER_DIR, "client/organization"), exist_ok=True)
    os.makedirs(os.path.join(HANDLER_DIR, "client/platform"), exist_ok=True)
    os.makedirs(os.path.join(HANDLER_DIR, "client/content"), exist_ok=True)
    os.makedirs(os.path.join(HANDLER_DIR, "client/common"), exist_ok=True)
    
    # 复制文件
    copy_file(
        os.path.join(HANDLER_DIR, "sidebar/work_contact_handler.go"),
        os.path.join(HANDLER_DIR, "client/contact/contact.go")
    )
    copy_file(
        os.path.join(HANDLER_DIR, "sidebar/work_room_handler.go"),
        os.path.join(HANDLER_DIR, "client/organization/room.go")
    )
    copy_file(
        os.path.join(HANDLER_DIR, "sidebar/work_agent_handler.go"),
        os.path.join(HANDLER_DIR, "client/platform/agent.go")
    )
    copy_file(
        os.path.join(HANDLER_DIR, "sidebar/medium_handler.go"),
        os.path.join(HANDLER_DIR, "client/content/medium.go")
    )
    copy_file(
        os.path.join(HANDLER_DIR, "sidebar/common_handler.go"),
        os.path.join(HANDLER_DIR, "client/common/common.go")
    )
    
    # 更新各子模块
    for filepath in Path(os.path.join(HANDLER_DIR, "client/contact")).glob("*.go"):
        update_file_package(str(filepath), "contact")
        update_file_imports(str(filepath),
                          r"mochat-api-server/internal/handler/sidebar",
                          "mochat-api-server/internal/handler/client/contact")
    
    for filepath in Path(os.path.join(HANDLER_DIR, "client/organization")).glob("*.go"):
        update_file_package(str(filepath), "organization")
        update_file_imports(str(filepath),
                          r"mochat-api-server/internal/handler/sidebar",
                          "mochat-api-server/internal/handler/client/organization")
    
    for filepath in Path(os.path.join(HANDLER_DIR, "client/platform")).glob("*.go"):
        update_file_package(str(filepath), "platform")
        update_file_imports(str(filepath),
                          r"mochat-api-server/internal/handler/sidebar",
                          "mochat-api-server/internal/handler/client/platform")
    
    for filepath in Path(os.path.join(HANDLER_DIR, "client/content")).glob("*.go"):
        update_file_package(str(filepath), "content")
        update_file_imports(str(filepath),
                          r"mochat-api-server/internal/handler/sidebar",
                          "mochat-api-server/internal/handler/client/content")
    
    for filepath in Path(os.path.join(HANDLER_DIR, "client/common")).glob("*.go"):
        update_file_package(str(filepath), "common")
        update_file_imports(str(filepath),
                          r"mochat-api-server/internal/handler/sidebar",
                          "mochat-api-server/internal/handler/client/common")

def main():
    print("="*60)
    print("Handler 目录重构迁移 - 完整版")
    print("="*60)
    
    try:
        migrate_content_module()
        migrate_marketing_module()
        migrate_analysis_module()
        migrate_platform_module()
        migrate_common_module()
        migrate_sidebar_to_client()
        
        print("\n" + "="*60)
        print("✅ 所有模块迁移完成！")
        print("="*60)
        print("\n下一步：")
        print("1. 检查编译: cd " + PROJECT_ROOT + " && go build ./...")
        print("2. 更新 router 引用路径")
        print("3. 运行测试验证")
        
    except Exception as e:
        print(f"\n❌ 迁移出错: {e}")
        import traceback
        traceback.print_exc()

if __name__ == "__main__":
    main()
