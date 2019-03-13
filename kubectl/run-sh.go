package main

import (
    "time"
)

func main() {
    for {
        <-time.After(time.Hour)
    }   
}
