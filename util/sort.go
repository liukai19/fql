package util

import (
	"sort"
	"strconv"
)

const (
	SORT_METHOD_DESC = "DESC"
	SORT_METHOD_ASC = "ASC"
)

func ArraySort(data[]map[string]string, key string, method string) []map[string]string {
	if len(data) <= 1 {
		return data
	}
	//if firstIndex < 0 || firstIndex > len(nums[0])-1 {
	//	fmt.Println("Warning: Param firstIndex should between 0 and len(nums)-1. The original array is returned.")
	//	return nums
	//}
	//排序
	mArray := &Array{data, key, method}
	sort.Sort(mArray)
	return mArray.mArr
}
type Array struct {
	mArr       []map[string]string
	firstIndex string
	method string
}
//IntArray实现sort.Interface接口
func (arr *Array) Len() int {
	return len(arr.mArr)
}
func (arr *Array) Swap(i, j int) {
	arr.mArr[i], arr.mArr[j] = arr.mArr[j], arr.mArr[i]
}
func (arr *Array) Less(i, j int) bool {
	arr1 := arr.mArr[i]
	arr2 := arr.mArr[j]
	//for index := arr.firstIndex; index < len(arr1); index++ {
	//	if arr1[index] < arr2[index] {
	//		return true
	//	} else if arr1[index] > arr2[index] {
	//		return false
	//	}
	//}
	arr1Val , _ := strconv.Atoi(arr1[arr.firstIndex])
	arr2Val , _ := strconv.Atoi(arr2[arr.firstIndex])
	if SORT_METHOD_DESC == arr.method {
		if arr1Val < arr2Val {
			return false
		} else if arr1Val > arr2Val {
			return true
		}
	} else {
		if arr1Val < arr2Val {
			return true
		} else if arr1Val > arr2Val {
			return false
		}
	}
	return i < j
}

