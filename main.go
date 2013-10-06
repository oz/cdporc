package main

import (
	"./api"
	"flag"
	"fmt"
	"net/url"
	"os"
)

const (
	defaultDelValue     = "ID"
	defaultListValue    = "all"
	defaultConfirmValue = "ID"
)

func InitalizeApi() *api.PaginatedQuotes {
	serverURLEnv := os.Getenv("CDP_SERVER")
	if serverURLEnv == "" {
		fmt.Println("You must set the CDP_SERVER environment variable to the server URL first!")
		os.Exit(3)
	}
	serverURL, err := url.Parse(serverURLEnv)
	if err != nil {
		fmt.Println("Invalid server URL in environment!")
		os.Exit(4)
	}

    return api.NewCoteDePorc(serverURL)
}

func DeleteQuote(api *api.PaginatedQuotes, id string) {
	fmt.Printf("Deleting quote %s... ", id)

	err := api.Delete(id)
	if err != nil {
		fmt.Println("Failed!", err.Error())
		return
	}
	fmt.Println("done.")
}

func ConfirmQuote(api *api.PaginatedQuotes, id string) {
	fmt.Printf("Publishing quote %s... ", id)

	err := api.Confirm(id)
	if err != nil {
		fmt.Println("Failed!", err.Error())
		return
	}
	fmt.Println("done.")
}

func RandomQuote(api *api.PaginatedQuotes) {
	quote, err := api.Random()
	if err != nil {
		fmt.Println("Failed!", err.Error())
		return
	}
	DisplayQuote(&quote)
}

func ListQuotes(api *api.PaginatedQuotes, kind string) {
	path := "/quotes"
	if kind != defaultListValue {
		path = path + "/" + kind
	}

	if err := api.GetAll(path); err != nil {
		panic(err)
	}
	if len(api.Entries) == 0 {
		fmt.Println("No quotes.")
		return
	}

	for _, quote := range api.Entries {
		DisplayQuote(&quote)
	}
}

func DisplayQuote(q *api.Quote) {
	fmt.Printf("%03d: %s\n", q.Id, q.Body)
}

func main() {
	var listQuotes bool
	var randomQuote bool
	var quoteType string
	var deleteId string
	var confirmId string

	flag.BoolVar(&listQuotes, "list", false, "List quotes")
	flag.BoolVar(&listQuotes, "l", false, "List quotes (shorthand)")
	flag.StringVar(&quoteType, "type", defaultListValue, "List quotes of this type")
	flag.StringVar(&quoteType, "t", defaultListValue, "List quotes of this type (shortand)")
	flag.StringVar(&deleteId, "del", defaultDelValue, "Delete quote")
	flag.StringVar(&deleteId, "d", defaultDelValue, "Delete quote (shorthand)")
	flag.StringVar(&confirmId, "publish", defaultConfirmValue, "Publish quote")
	flag.StringVar(&confirmId, "p", defaultConfirmValue, "Publish quote (shorthand)")
	flag.BoolVar(&randomQuote, "rand", false, "Get a random quote")
	flag.BoolVar(&randomQuote, "r", false, "Get a random quote (shorthand)")
	flag.Parse()
	api := InitalizeApi()

	if deleteId != defaultDelValue {
		DeleteQuote(api, deleteId)
	}

	if confirmId != defaultConfirmValue {
		ConfirmQuote(api, confirmId)
	}

	if randomQuote {
		RandomQuote(api)
	}

	if listQuotes {
		ListQuotes(api, quoteType)
	}

}
