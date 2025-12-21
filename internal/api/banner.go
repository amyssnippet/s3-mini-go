package api


import "fmt"

func PrintBanner() {
    fmt.Println(`
   _____ ____      __  __ _       _ 
  / ____|___ \    |  \/  (_)     (_)
 | (___   __) |___| \  / |_ _ __  _ 
  \___ \ |__ <____| |\/| | | '_ \| |
  ____) |___) |   | |  | | | | | | |
 |_____/|____/    |_|  |_|_|_| |_|_|
    `)
    fmt.Println("  v1.0.0 - Secure P2P Object Store")
}