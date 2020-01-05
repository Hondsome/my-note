package main

// 快速排序
// 快速排序思路
// 快速排序通过分支法的思想，从一个数组中选取一个基准元素pivot，把这个数组中小于pivot的移动到左边，把大于pivot的移动到右边。然后再分别对左右两边数组进行快速排序。

// 双边循环法
// 思路
// 设置两个指针left和right，最初分别指向数组的左右两端。比较right指针指向元素和pivot元素，如果right元素大于pivot元素，right指针左移一位，再和pivot进行比较，如果right元素小于pivot元素的话停止移动，换到left指针。
// left指针的操作是，left指针指向的元素和pivot元素比较，如果left指向元素小于或等于pivot，left指针右移，如果left元素大于pivot元素，停止移动。
// 左右都停止移动后，交换left和right指向的元素，这样left指针指向的是一个小于pivot的元素，right指向的是一个大于pivot的元素。
// 当left和right重叠的时候结束比较，将第一个元素和left，right指向的元素做交换，完成一轮排序

// 代码
func partition(arr []int, startIndex, endIndex int) int {
    var (
        pivot = arr[startIndex]
        left  = startIndex
        right = endIndex
    )

    for left != right {
        for left < right && pivot < arr[right] {
            right--
        }

        for left < right && pivot >= arr[left] {
            left++
        }

        if left < right {
            arr[left], arr[right] = arr[right], arr[left]
        }
    }

    arr[startIndex], arr[left] = arr[left], arr[startIndex]

    return left
}
// 单边循环法
// 思路
// 单边循环代码实现简单。
// 通过一个mark指针，指向小于pivot的集合的最后一个元素，最后把第一个元素和mark指向的元素做交换，进行下一轮。
// mark指针开始指向第一个元素，然后开始遍历数组，如果当前元素比pivot大，继续遍历，如果比pivot小，mark指针右移，将mark指向元素和当前遍历元素交换。

// 代码
func partitionv2(arr []int, startIndex, endIndex int) int {
    var (
        mark  = startIndex
        pivot = arr[startIndex]
        point = startIndex + 1
    )

    for point < len(arr) {
        if arr[point] < pivot {
            mark++
            arr[mark], arr[point] = arr[point], arr[mark]
        }
        point++
    }

    arr[startIndex], arr[mark] = arr[mark], arr[startIndex]
    return mark
}
// pivot选择
// 有数组5，4，3，2，1要进行排序，如果选择第一个元素作为pivot的话，，每次选择的都是该数组中的最大值或最小值，每次进行排序只确定了一个元素的位置，导致时间复杂度退化成O(n^2)
// 在选择pivot时，可以用随机选择的方式选择，即在当前数组中随机选择一个元素来作为pivot，减少选择到最大值或最小值的几率。