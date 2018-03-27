package main
import (
"fmt"
"time"
)
func main() {
     message := make(chan string) // no buffer
     count := 3

     go func() {
          for i := 1; i <= count; i++ {
               fmt.Println("send message")
               message <- fmt.Sprintf("delete")
          }
     }()

   time.Sleep(time.Minute * 1)

    // for i := 1; i <= count; i++ {
          fmt.Println(<-message)
    // }
}
