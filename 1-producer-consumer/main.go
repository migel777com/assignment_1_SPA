//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

func producer(stream Stream) (tweets []*Tweet) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			return tweets
		}

		tweets = append(tweets, tweet)
	}
}

func writer(stream Stream, ch chan Tweet) { // new functiong for writing to channel
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(ch)
			return
		}

		ch <- *tweet
	}
}

func consumer(tweets []*Tweet) {
	for _, t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func listener(ch chan Tweet, wg *sync.WaitGroup) { // new function for listening tweets from channel
	defer wg.Done()
	for {
		t := <-ch
		if t != (Tweet{}) {
			if t.IsTalkingAboutGo() {
				fmt.Println(t.Username, "\ttweets about golang")
			} else {
				fmt.Println(t.Username, "\tdoes not tweet about golang")
			}
		} else {
			break
		}
	}
}

func main() {
	var wg sync.WaitGroup
	start := time.Now()
	stream := GetMockStream()
	ch := make(chan Tweet)

	// Producer
	//tweets := producer(stream)
	go writer(stream, ch)

	// Consumer
	//consumer(tweets)
	wg.Add(1)
	go listener(ch, &wg)
	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}

// I took all functions and structs from mockstream file to check my work
// GetMockStream is a blackbox function which returns a mock stream for
// demonstration purposes
func GetMockStream() Stream {
	return Stream{0, mockdata}
}

// Stream is a mock stream for demonstration purposes, not threadsafe
type Stream struct {
	pos    int
	tweets []Tweet
}

// ErrEOF returns on End of File error
var ErrEOF = errors.New("End of File")

// Next returns the next Tweet in the stream, returns EOF error if
// there are no more tweets
func (s *Stream) Next() (*Tweet, error) {

	// simulate delay
	time.Sleep(320 * time.Millisecond)
	if s.pos >= len(s.tweets) {
		return &Tweet{}, ErrEOF
	}

	tweet := s.tweets[s.pos]
	s.pos++

	return &tweet, nil
}

// Tweet defines the simlified representation of a tweet
type Tweet struct {
	Username string
	Text     string
}

// IsTalkingAboutGo is a mock process which pretend to be a sophisticated procedure to analyse whether tweet is talking about go or not
func (t *Tweet) IsTalkingAboutGo() bool {
	// simulate delay
	time.Sleep(330 * time.Millisecond)

	hasGolang := strings.Contains(strings.ToLower(t.Text), "golang")
	hasGopher := strings.Contains(strings.ToLower(t.Text), "gopher")

	return hasGolang || hasGopher
}

var mockdata = []Tweet{
	{
		"davecheney",
		"#golang top tip: if your unit tests import any other package you wrote, including themselves, they're not unit tests.",
	}, {
		"beertocode",
		"Backend developer, doing frontend featuring the eternal struggle of centering something. #coding",
	}, {
		"ironzeb",
		"Re: Popularity of Golang in China: My thinking nowadays is that it had a lot to do with this book and author https://github.com/astaxie/build-web-application-with-golang",
	}, {
		"beertocode",
		"Looking forward to the #gopher meetup in Hsinchu tonight with @ironzeb!",
	}, {
		"vampirewalk666",
		"I just wrote a golang slack bot! It reports the state of github repository. #Slack #golang",
	},
}

