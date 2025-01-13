<div align="center">
  <img src="newt.png" alt="Newt Logo" width="200" />
</div>

# Newt 🦎

An attempt to make a smaller, cuter, and more user-friendly HTTP client for Go.

## Sample Usage

```golang
package main

import (
	"fmt"
	"io"

	newt "github.com/douglasgreyling/newt/client"
)

func main() {
	client, err := newt.NewClient("https://jsonplaceholder.typicode.com")
	if err != nil {
		fmt.Println("Error creating client:", err)
		return
	}

	params := newt.Params{"title": "Newts 101", "body": "A newt is a salamander", "userId": "1"}

	resp, err := client.Post("/posts", params)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println(string(body))
}
```