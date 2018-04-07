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
           requested = <- token 
           fmt.Println("Requested: ",requested)
       }

    }
}

type LightCheckpoint struct {
     bookmark []int
}

type Checkpoint struct {
     bookmark []int 
     request chan *LightCheckpoint
     fastforward bool
     req *LightCheckpoint
}

func (b *Checkpoint) push() {
    if len(b.bookmark) < cap(b.bookmark) {
        b.bookmark = b.bookmark[:len(b.bookmark)+1]
    } else {
        fmt.Println("Unexpected length in push",len(b.bookmark),cap(b.bookmark))
    }
}

func (b *Checkpoint) pop() {
    if len(b.bookmark) > 1 {
        b.bookmark = b.bookmark[:len(b.bookmark)-1]
    } else {
        fmt.Println("Unexpected length in pop",len(b.bookmark),cap(b.bookmark))
    }

}

func (b *Checkpoint) p() {
    fmt.Println(b.bookmark)
}

func (b *Checkpoint) mark() {
    b.bookmark[len(b.bookmark)-1]++
    if b.req == nil {
        b.req = <- b.request
    }
    if b.compareTo(b.req) >= 0 {
        fmt.Println("Request satisfied")
        
    }
}

func (b *Checkpoint) reset() {
    b.bookmark = b.bookmark[0:1]
    b.bookmark[0] = 0
    b.fastforward = false
}

func (b *Checkpoint) compareTo(other *LightCheckpoint) int {
    i := 0
    for i < len(b.bookmark){
        if i >= len(other.bookmark) {
            return 1
	} else {
            mine := b.bookmark[i]
            theirs := other.bookmark[i]
            if mine < theirs {
                return -1
            } else if mine > theirs {
                return 1
            }
            i++
        }
    }
    if i < len(other.bookmark) {
        return -1
    } else {
        return 0
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

    //token := make(chan int)
    //go count(token)
    //token <- 13
    //i := <- token
    //fmt.Println(i)
    //token <- 53
    //i = <- token
    //fmt.Println("I", i)
    //token <- 7
    //i = <- token
    //fmt.Println("I", i)

    req := LightCheckpoint{bookmark : make([]int,1,7)}
    req.bookmark[0] = 1
    b := &Checkpoint{bookmark : make([]int,1,7)}
    fmt.Println(b.request)
    t := make(chan int,1)
    fmt.Println(t)
    t <- 1
    b.request <- &req
    //b.mark()
    //b.p()
    //b.mark()
    //b.p()
    //b.push()
    //b.p()
    //b.mark()
    //b.p()
    ////b.pop()
    ////b.p()

    //c := make([]int,1,7)
    //fmt.Println(c)
    //c = c[:2]
    //fmt.Println(c)

    //d := &Checkpoint{bookmark : []int{1,2}}
    //fmt.Println(b.compareTo(d))
    //fmt.Println(d.compareTo(b))
    //fmt.Println(d.compareTo(d))
    //fmt.Println(b.compareTo(b))
    //b.pop()
    //fmt.Println(b.compareTo(d))
    //fmt.Println(d.compareTo(b))
    //fmt.Println(d.compareTo(d))
    //fmt.Println(b.compareTo(b))

}
