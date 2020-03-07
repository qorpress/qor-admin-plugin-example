package main
 
import "math/rand"
 
 
func Sort(items []int) *[]int {
    if len(items) < 2 {
        return &items
    }
 
    peg := items[rand.Intn(len(items))]
 
    below := make([]int, 0, len(items))
    above := make([]int, 0, len(items))
    middle := make([]int, 0, len(items))
 
    for _, item := range items {
        switch {
        case item < peg:
            below = append(below, item)
        case item == peg:
            middle = append(middle, item)
        case item > peg:
            above = append(above, item)
        }
    }
 
    below = *Sort(below)
    above = *Sort(above)
 
    sorted := append(append(below, middle...), above...)
    return &sorted
}
