## 排序算法

###1.插入排序

    insert-sort(A):
        for j =2 to A.length:
            key = A[j]
            i = j-1
            while i>0 and  A[i]>key :
                A[i+1] = A[j]
                i--
        A[i+1] = key 
        // 把key插入到合适的位置去
        
###2.并归排序
    

###3.快排


###4.冒泡

###5.桶排序

###各排序算法特点，优劣对比，使用场景