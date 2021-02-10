package main

import "fmt"

func main() {
    data, err := loadGameContent()
    if err != nil {
        fmt.Println("Error in loading game content")
    }
    for _, wonder := range data.Wonders {
        wonder.print()
        fmt.Println()
    }
}
