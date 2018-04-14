package main

import "fmt"

type Checkpoint struct {
    bookmark [] int
}

type MachineState struct {
    request chan Checkpoint
    internalCounter Checkpoint
    state int
}

func requestReplyLoop(msChan chan MachineState) {
    ms := <- msChan
    fmt.Println("Recieved ms")
    requestChan := ms.request
    fmt.Println("Waiting for request")
    request := <- requestChan
    fmt.Println("Received request",request)
    ms.state = 1
    msChan <- ms
      

}

func main() {
    ms := MachineState{request : make(chan Checkpoint), internalCounter : Checkpoint{bookmark : make([]int,1,7)}}
    request := Checkpoint{bookmark : make([]int,1,7)}
    requestChan := ms.request
    msChan := make(chan MachineState)
    go requestReplyLoop(msChan)
    msChan <- ms
    fmt.Println("Sending request")
    requestChan <- request
    fmt.Println("Waiting for state to come back")
    ms = <- msChan
    fmt.Println("Received state value ",ms.state)
}
