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
    fmt.Println("Done is true")
    x:= <- keepgo
    fmt.Println("Keepgo",keepgo,x,"SOMETHINGELSE")
    done <- false
    fmt.Println("done",done)
    
    
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

type MachineState struct {
    value int
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
    fmt.Println("Marking")
    b.bookmark[len(b.bookmark)-1]++
    if b.req == nil {
	    fmt.Println("Checked for nil")
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

func countToTen(c *Checkpoint, s chan MachineState) {
    fmt.Println("Counting to ten")
    i := 0
    for i < 10 {
	fmt.Println("I: ",i)
	i++
        c.mark()
        s<- MachineState{value : i}
    }

}

func requestReplyLoop(c *Checkpoint, lc chan LightCheckpoint, s chan MachineState) {
    fmt.Println("RequestReplyLoop",c.fastforward)
    c.
    k := <-lc
    fmt.Println("Received request",k)
    fmt.Println("Generate reply")
    ms := MachineState{value:1}
    s<-ms
    k = <-lc
    fmt.Println("Received request",k)
    fmt.Println("Generate reply")
    ms.value=2
    s<-ms
}

func doAnything() {
    fmt.Println("Do anything")
}

func main() {
    //fmt.Println("Hello, World")
    //a := [5]int{0,1,2,3,4}
    //for i := 0; i < 5; i++ {
    //    fmt.Println("A[",i,"]",a[i])
    //}
    //h := mclose()
    //fmt.Println(h())
    //fmt.Println(h())
    //fmt.Println(h())
    //fmt.Println(h())

    //done := make(chan bool, 1)
    //keepgo := make(chan bool,1)
    //go do1(done,keepgo)
    //
    //<-done
    //fmt.Println("Received done: ", done)
    //fmt.Println("All ok, now return to do 1")
    //keepgo<-false
    //fmt.Println("After keepgo is false")
    //<-done

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

//
    req := LightCheckpoint{bookmark : make([]int,1,7)}
    req.bookmark[0] = 1
    //b := &Checkpoint{bookmark : make([]int,1,7)}
    c := make(chan LightCheckpoint)
    s := make(chan MachineState)
    go requestReplyLoop(c,s)
    lc := LightCheckpoint{bookmark : make([]int,1,7)}

    fmt.Println("Generate request")
    c <- lc
    fmt.Println("Receive state")
    ms := <-s
    fmt.Println("Received state",ms)
    lc.bookmark[0]++

    fmt.Println("Generate request")
    c <- lc
    fmt.Println("Receive state")
    ms = <-s
    fmt.Println("Received state",ms)
   
    //go countToTen(b, s)
    //fmt.Println("Sending request",req)
    //b.request <- &req
    //fmt.Println("Made it this far")
    //ms := <-s
    //fmt.Println("MS: ", ms)

    //fmt.Println(b.request)
    //t := make(chan int,1)
    //fmt.Println(t)
    //t <- 1
//

    //b.request <- &req
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
