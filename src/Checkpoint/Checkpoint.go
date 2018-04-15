package Checkpoint

import "fmt"

type Checkpoint struct {
    bookmark [] int
}

func (b *Checkpoint) compareTo(other *Checkpoint) int {
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

