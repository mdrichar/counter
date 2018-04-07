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
     request chan LightCheckpoint
     response chan MachineState
     activeRequest LightCheckpoint
     currentState *MachineState
     hasActiveRequest bool
     fastforward bool
}

type MachineState struct {
    value int
}

func (s *MachineState) reset() {
    s.value = 0
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
    // Can only read the request once, need to save it somewhere for future use
    // As soon as we get it once, copy it over to activeRequest and use that subsequently (so we don't query the channel again) until the request has been satisfied
    if b.hasActiveRequest == false {
        fmt.Println("Wait for request")
        b.activeRequest = <- b.request
        fmt.Println("Request received", b.activeRequest)
        if b.compareTo(&b.activeRequest) > 0 {
	    fmt.Println("Rewinding")
	    b.fastforward = true
	}
        b.hasActiveRequest = true
    }
    if !b.fastforward && b.compareTo(&b.activeRequest) >= 0 {
        fmt.Println("Request satisfied")
        b.response <- *b.currentState
        b.hasActiveRequest = false // Next time through loop, we will block until request comes over the channel
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

func requestReplyLoop(lc chan LightCheckpoint, s chan MachineState) {
    fmt.Println("RequestReplyLoop")
    ms := MachineState{value:1}
    c := Checkpoint{request : lc, response : s, fastforward : false, bookmark : make([]int,1,7), hasActiveRequest : false, currentState : &ms}

    for {
	c.fastforward = false
	for i := 0; i < 6; i++ {
            c.mark()
            ms.value++
        }
    }

    //k := <-lc
    //fmt.Println("Received request",k)
    //fmt.Println("Generate reply")
    //s<-ms
    //k = <-lc
    //fmt.Println("Received request",k)
    //fmt.Println("Generate reply")
    //ms.value=2
    //s<-ms
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
    lc.bookmark[0] += 3

    fmt.Println("Generate request")
    c <- lc
    fmt.Println("Receive state")
    ms = <-s
    fmt.Println("Received state",ms)
   
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
