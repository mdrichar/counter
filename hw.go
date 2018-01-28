package main

import "fmt"

func mclose() func() int {
    i := 0
    return func() int {
        i++
        return i
    }

}

func do1(done chan bool, keepgo chan bool) {
    fmt.Println("Do1 called")
    done <- true 
    x:= <- keepgo
    fmt.Println("Keepgo",keepgo,x)
    
    
}

func counter(produced chan int, request chan int) {
    item := 0
    requested := <- request
    fmt.Println("Requested: ",requested)
    for {
       item = (item + 1) % 100
       if item == requested + 1 {
           produced <- item
           fmt.Println("Produced: ",item)
           requested := <- request
           fmt.Println("Requested: ",requested)
       }

    }
}

func count(token chan int) {
    var item int = 0
    requested := <- token 
    fmt.Println("Requested: ",requested)
    for {
       item = (item + 1) % 100
       if item == requested + 1 {
           token <- item
           fmt.Println("Produced: ",item)
           requested := <- token 
           fmt.Println("Requested: ",requested)
       }

    }
}

func main() {
    fmt.Println("Hello, World")
    a := [5]int{0,1,2,3,4}
    for i := 0; i < 5; i++ {
        fmt.Println("A[",i,"]",a[i])
    }
    h := mclose()
    fmt.Println(h())
    fmt.Println(h())
    fmt.Println(h())
    fmt.Println(h())

    done := make(chan bool, 1)
    keepgo := make(chan bool,1)
    go do1(done,keepgo)
    
    <-done
    fmt.Println("All ok, now return to do 1")
    keepgo<-false
    fmt.Println("After keepgo is false")

    token := make(chan int)
    go count(token)
    token <- 13
    i := token
    fmt.Println(i)

}
