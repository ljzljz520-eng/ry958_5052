package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"campus-creative-store"
)

func main() {
	service := campusstore.NewFixture()
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		address := ":8080"
		if len(os.Args) > 2 {
			address = os.Args[2]
		}
		log.Printf("校园文创商品站运行于 http://localhost%s", address)
		log.Fatal(http.ListenAndServe(address, campusstore.NewHandler(service)))
	}
	fmt.Println("校园文创商品站")
	for _, product := range service.Products() {
		fmt.Printf("%s | ¥%s\n", product.Name, product.Price.StringFixed(2))
	}
	fmt.Println("使用 go run ./cmd/campus-store serve 启动 HTTP 入口")
}
