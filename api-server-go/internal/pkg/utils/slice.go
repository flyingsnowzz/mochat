package utils

// InArray 检查值是否在切片中（类似 PHP 的 in_array）
func InArray[T comparable](item T, items []T) bool {
	for _, v := range items {
		if v == item {
			return true
		}
	}
	return false
}

// Column 从结构体切片中提取指定字段的值（类似 PHP 的 array_column）
func Column[T any, K comparable](items []T, fn func(T) K) []K {
	result := make([]K, 0, len(items))
	for _, item := range items {
		result = append(result, fn(item))
	}
	return result
}

// ColumnMap 从结构体切片中提取 key-value 映射（类似 PHP 的 array_column 双参数模式）
func ColumnMap[T any, K comparable, V any](items []T, keyFn func(T) K, valFn func(T) V) map[K]V {
	result := make(map[K]V, len(items))
	for _, item := range items {
		result[keyFn(item)] = valFn(item)
	}
	return result
}

// Filter 过滤切片（类似 PHP 的 array_filter）
func Filter[T any](items []T, fn func(T) bool) []T {
	result := make([]T, 0, len(items))
	for _, item := range items {
		if fn(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map 对切片每个元素执行转换（类似 PHP 的 array_map）
func Map[T any, R any](items []T, fn func(T) R) []R {
	result := make([]R, 0, len(items))
	for _, item := range items {
		result = append(result, fn(item))
	}
	return result
}

// Chunk 将切片分块（类似 PHP 的 array_chunk）
func Chunk[T any](items []T, size int) [][]T {
	if size <= 0 || len(items) == 0 {
		return nil
	}
	chunks := make([][]T, 0, (len(items)+size-1)/size)
	for i := 0; i < len(items); i += size {
		end := i + size
		if end > len(items) {
			end = len(items)
		}
		chunks = append(chunks, items[i:end])
	}
	return chunks
}

// Unique 切片去重（类似 PHP 的 array_unique）
func Unique[T comparable](items []T) []T {
	seen := make(map[T]struct{}, len(items))
	result := make([]T, 0, len(items))
	for _, item := range items {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}
