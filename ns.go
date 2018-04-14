package main

import "fmt"

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
    svalues []int
    i int
    j int
    temp int
}

func (s *MachineState) reset() {
    s.value = 0
}

func (s* MachineState) print() {
    fmt.Println("i",s.i,"j",s.j)
    for k := 0; k < len(s.svalues); k++ {
        fmt.Println(k,s.svalues[k])
    }
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
    fmt.Println("Comparing: ",b.bookmark," to ",other.bookmark)
    for i < len(b.bookmark){
        if i >= len(other.bookmark) {
	    fmt.Println("That one comes first on length comparision",i)
            return 1
	} else {
            mine := b.bookmark[i]
            theirs := other.bookmark[i]
            if mine < theirs {
	        fmt.Println("This one comes first on cell comparision",mine,theirs)
                return -1
            } else if mine > theirs {
		fmt.Println("That one comes first on cell comparision",mine,theirs)
                return 1
            }
            i++
        }
    }
    if i < len(other.bookmark) {
	fmt.Println("This one comes first on length comparison",i)
        return -1
    } else {
	fmt.Println("They're equal")
        return 0
    }
    fmt.Println("Should not get this far.")
    return 0
}

func requestReplyLoop(lc chan LightCheckpoint, s chan MachineState) {
    fmt.Println("RequestReplyLoop")

    ms := MachineState{value:1, svalues:[]int{1,6,4,2,3}}
    c := Checkpoint{request : lc, response : s, fastforward : false, bookmark : make([]int,1,7), hasActiveRequest : false, currentState : &ms}
    for {
	c.fastforward = false
	for ms.i = 0; ms.i < len(ms.svalues)-1; ms.i++ {
            for ms.j = ms.i+1; ms.j < len(ms.svalues); ms.j++ {
                if ms.svalues[ms.j] < ms.svalues[ms.i] {
                    ms.temp = ms.svalues[ms.i]
                    ms.svalues[ms.i] = ms.svalues[ms.j]
                    ms.svalues[ms.j] = ms.temp		    
                }
                c.mark()
                ms.value++
	    }
        }
    }
}

func main() {
    req := LightCheckpoint{bookmark : make([]int,1,7)}
    req.bookmark[0] = 1
    //b := &Checkpoint{bookmark : make([]int,1,7)}
    c := make(chan LightCheckpoint)
    s := make(chan MachineState)
    go requestReplyLoop(c,s)
    lc := LightCheckpoint{bookmark : make([]int,1,7)}
    lc.bookmark[0] += 1

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

}
