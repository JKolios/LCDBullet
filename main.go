package main

import (
	"github.com/JKolios/goLcd/lcd"
	_ "github.com/kidoman/embd/host/rpi"
	"time"
)

func main() {
	display := lcd.NewDisplay(21, 16, 6, 13, 19, 26, 5, true)
	defer display.Close()

	display.FlashDisplay(3, 5*time.Second)
	display.DisplayMessage("Hello!", 5*time.Second)
	display.DisplayMessage("0123456789ABCDEF0123456789ABCDEF", 5*time.Second)
	display.DisplayMessage("Has Anyone Really Been Far Even as Decided to Use Even Go Want to do Look More Like", 5*time.Second)

}
