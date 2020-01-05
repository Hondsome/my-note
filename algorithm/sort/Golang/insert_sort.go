package main
import "fmt"

func InsertSorted(nums []int) []int{
	for j:=1; j < len(nums); j++{
		key := nums[j]
		i := j-1
		// 把小的数放前面。大于key数往后移一位，找到key合适位置
		for i >= 0 && nums[i]>key {
			nums[i+1] = nums[i]
			i -=1
		}
		nums[i+1] = key
	}
	return  nums
}

func main() {
	a := []int{2,15,6,7,3,1,2}
	fmt.Printf("sorted nums :%v \n",InsertSorted(a))
}