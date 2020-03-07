package main
 
 
func Sort(items []int) *[]int {
    if len(items) < 2 {
        return &items
    }
 
    tmp := 0
 
    for i := 0; i < len(items); i++ {
        for j := 0; j < len(items)-1; j++ {
            if items[j] > items[j+1] {
                tmp = items[j]
                items[j] = items[j+1]
                items[j+1] = tmp
            }
        }
    }
 
    return &items
}
