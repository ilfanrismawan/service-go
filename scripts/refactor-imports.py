#!/usr/bin/env python3
"""
Smart import path updater for domain-based refactoring.
This script intelligently updates import paths based on file location and content.
"""

import os
import re
import sys
from pathlib import Path

# Mapping of old imports to new imports based on domain
IMPORT_MAPPINGS = {
    # Shared resources
    'service/internal/config': 'service/internal/shared/config',
    'service/internal/database': 'service/internal/shared/database',
    'service/internal/middleware': 'service/internal/shared/middleware',
    'service/internal/utils': 'service/internal/shared/utils',
    'service/internal/notification': 'service/internal/shared/notification',
    'service/internal/monitoring': 'service/internal/shared/monitoring',
    
    # Payment legacy
    'service/internal/payment': 'service/internal/payments/legacy_payment',
}

# Domain-specific mappings
DOMAIN_MAPPINGS = {
    'internal/users': {
        'service/internal/core': 'service/internal/users/dto',
        'service/internal/repository': 'service/internal/users/repository',
    },
    'internal/orders': {
        'service/internal/core': 'service/internal/orders/dto',
        'service/internal/service': 'service/internal/orders/service',
        'service/internal/repository': 'service/internal/orders/repository',
    },
    'internal/payments': {
        'service/internal/core': 'service/internal/payments/dto',
        'service/internal/service': 'service/internal/payments/service',
        'service/internal/repository': 'service/internal/payments/repository',
    },
    'internal/branches': {
        'service/internal/core': 'service/internal/branches/dto',
        'service/internal/service': 'service/internal/branches/service',
        'service/internal/repository': 'service/internal/branches/repository',
    },
    'internal/shared': {
        'service/internal/core': 'service/internal/shared/model',
    },
}

# Shared model types that should stay in shared/model
SHARED_MODEL_TYPES = [
    'APIResponse', 'ErrorResponse', 'PaginationResponse', 'PaginatedResponse',
    'SuccessResponse', 'CreateErrorResponse', 'PaginatedSuccessResponse',
    'ErrUnauthorized', 'ErrForbidden', 'ErrNotFound', 'ErrInvalidInput',
    'ErrInternalError', 'CalculateDistance',
]

def detect_shared_model_usage(content):
    """Detect if file uses shared model types"""
    for model_type in SHARED_MODEL_TYPES:
        if model_type in content:
            return True
    return False

def update_imports_in_file(filepath):
    """Update import paths in a single file"""
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
        
        original_content = content
        updated = False
        
        # Determine domain from file path
        filepath_str = str(filepath)
        domain_mapping = None
        for domain, mappings in DOMAIN_MAPPINGS.items():
            if domain in filepath_str:
                domain_mapping = mappings
                break
        
        # Update domain-specific imports
        if domain_mapping:
            for old_import, new_import in domain_mapping.items():
                # Only replace if it's a full import line
                pattern = rf'(\s+)"{re.escape(old_import)}"'
                if re.search(pattern, content):
                    content = re.sub(pattern, rf'\1"{new_import}"', content)
                    updated = True
        
        # Update shared resource imports
        for old_import, new_import in IMPORT_MAPPINGS.items():
            pattern = rf'(\s+)"{re.escape(old_import)}"'
            if re.search(pattern, content):
                content = re.sub(pattern, rf'\1"{new_import}"', content)
                updated = True
        
        # Special handling: core references that should be shared/model
        if 'shared' in filepath_str and 'service/internal/core' in content:
            # Check if it's using shared models
            if detect_shared_model_usage(content):
                pattern = r'"service/internal/core"'
                if pattern in content:
                    content = content.replace(pattern, '"service/internal/shared/model"')
                    updated = True
        
        # Update core references to dto/shared based on context
        # This is a heuristic - may need manual review
        if domain_mapping and 'core.' in content:
            # Replace core. references with dto. for domain files
            if 'internal/users' in filepath_str:
                content = re.sub(r'\bcore\.', 'dto.', content)
                updated = True
            elif 'internal/orders' in filepath_str:
                content = re.sub(r'\bcore\.', 'dto.', content)
                updated = True
            elif 'internal/payments' in filepath_str:
                content = re.sub(r'\bcore\.', 'dto.', content)
                updated = True
            elif 'internal/branches' in filepath_str:
                content = re.sub(r'\bcore\.', 'dto.', content)
                updated = True
        
        if updated and content != original_content:
            with open(filepath, 'w', encoding='utf-8') as f:
                f.write(content)
            return True
        
        return False
    except Exception as e:
        print(f"Error processing {filepath}: {e}", file=sys.stderr)
        return False

def update_package_name(filepath):
    """Update package name based on directory structure"""
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
        
        original_content = content
        filepath_str = str(filepath)
        
        # Determine new package name
        new_package = None
        if 'internal/users/auth' in filepath_str:
            new_package = 'auth'
        elif 'internal/users/handler' in filepath_str:
            new_package = 'handler'
        elif 'internal/users/service' in filepath_str:
            new_package = 'service'
        elif 'internal/users/repository' in filepath_str:
            new_package = 'repository'
        elif 'internal/users/dto' in filepath_str:
            new_package = 'dto'
        elif 'internal/orders/handler' in filepath_str:
            new_package = 'handler'
        elif 'internal/orders/service' in filepath_str:
            new_package = 'service'
        elif 'internal/orders/repository' in filepath_str:
            new_package = 'repository'
        elif 'internal/orders/dto' in filepath_str:
            new_package = 'dto'
        elif 'internal/payments/handler' in filepath_str:
            new_package = 'handler'
        elif 'internal/payments/service' in filepath_str:
            new_package = 'service'
        elif 'internal/payments/repository' in filepath_str:
            new_package = 'repository'
        elif 'internal/payments/dto' in filepath_str:
            new_package = 'dto'
        elif 'internal/payments/legacy_payment' in filepath_str:
            new_package = 'legacy_payment'
        elif 'internal/branches/handler' in filepath_str:
            new_package = 'handler'
        elif 'internal/branches/service' in filepath_str:
            new_package = 'service'
        elif 'internal/branches/repository' in filepath_str:
            new_package = 'repository'
        elif 'internal/branches/dto' in filepath_str:
            new_package = 'dto'
        elif 'internal/shared' in filepath_str:
            # Extract subdirectory name
            parts = filepath_str.split('internal/shared/')
            if len(parts) > 1:
                subdir = parts[1].split('/')[0]
                if subdir == 'handlers':
                    new_package = 'handlers'
                else:
                    new_package = subdir
        
        if new_package:
            # Update package declaration
            pattern = r'^package\s+\w+'
            if re.match(pattern, content):
                content = re.sub(pattern, f'package {new_package}', content, count=1)
        
        if content != original_content:
            with open(filepath, 'w', encoding='utf-8') as f:
                f.write(content)
            return True
        
        return False
    except Exception as e:
        print(f"Error updating package in {filepath}: {e}", file=sys.stderr)
        return False

def main():
    """Main function"""
    base_dir = Path(__file__).parent.parent
    
    # Find all Go files in new domain structure
    go_files = []
    for domain_dir in ['internal/users', 'internal/orders', 'internal/payments', 
                       'internal/branches', 'internal/shared']:
        domain_path = base_dir / domain_dir
        if domain_path.exists():
            go_files.extend(domain_path.rglob('*.go'))
    
    print(f"Found {len(go_files)} Go files to process")
    
    updated_count = 0
    for go_file in go_files:
        # Update package name first
        if update_package_name(go_file):
            updated_count += 1
            print(f"✓ Updated package in {go_file.relative_to(base_dir)}")
        
        # Then update imports
        if update_imports_in_file(go_file):
            updated_count += 1
            print(f"✓ Updated imports in {go_file.relative_to(base_dir)}")
    
    print(f"\n✅ Updated {updated_count} files")

if __name__ == '__main__':
    main()

