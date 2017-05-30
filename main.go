
package main

import "fmt"
import "net/http"
import "io/ioutil"
import "regexp"
import "os"
import "bufio"
import "time"



func main() {
	// set threads
	threads := 0
	total := 0
	// go through stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// sleep if too many threads
		for threads >= 5 {
			time.Sleep(time.Millisecond * 100)
		}
		// prepare url
		r, _ := regexp.Compile("\n+$")
		line := r.ReplaceAllString(scanner.Text(), "")
		// launch goroutine
		threads++
		go do_request( line, threads, &threads, &total)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	// wait all threads done
	for threads > 0 {
		time.Sleep(time.Millisecond * 100)
	}
	fmt.Println("TOTAL:", total)
}

func do_request( url string, thread_num int, threads *int, total *int ) {

	// do http request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Bad URL: " + url)
		return
	}
	defer resp.Body.Close()

	// read body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Bad URL resp: ", resp)
		return
	}

	// search body
	r, _ := regexp.Compile("Go")
	go_matches := len(r.FindAll(body, -1))
	*total = *total + go_matches

	// output result
	fmt.Println(thread_num, "URL:", url,", matches:", go_matches)

	// decrease threads
	*threads = *threads - 1
}

