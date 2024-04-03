package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfBuburs = 10

var buburMade, buburFailed, total int

// Producer is a type for structs that holds two channels: one for pizzas, with all
// information for a given pizza order including whether it was made
// successfully, and another to handle end of processing (when we quit the channel)
type Producer struct {
	data chan BuburOrder
	quit chan chan error
}

// PizzaOrder is a type for structs that describes a given pizza order. It has the order
// number, a message indicating what happened to the order, and a boolean
// indicating if the order was successfully completed.
type BuburOrder struct {
	buburNumber int
	message     string
	success     bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	color.Cyan("The Bubur Asli Depok is open for business!")
	color.Cyan("----------------------------------")
	fmt.Println()

	orderData := make(chan BuburOrder)
	quit := make(chan chan error)

	go buburProducer(orderData, quit)
	go buburConsumer(orderData, quit)

	<-quit
	close(quit)
}

func buburProducer(orderData chan BuburOrder, quitTrigger chan chan error) {
	i := 0

	for {
		currentBubur := makeBubur(i)
		if currentBubur != nil {
			i = currentBubur.buburNumber
			select {
			case orderData <- *currentBubur:

			case quitChan := <-quitTrigger:
				close(orderData)
				close(quitChan)
			}
		}
	}
}

func makeBubur(buburNumber int) *BuburOrder {
	buburNumber++

	if buburNumber <= NumberOfBuburs {
		var msg string
		var success bool
		delay := rand.Intn(3) + 1
		fmt.Printf("Received order #%d!\n", buburNumber)

		rnd := rand.Intn(10) + 1

		if rnd < 5 {
			buburFailed++
		} else {
			buburMade++
		}

		total++

		fmt.Printf("making bubur number %d, it will take %d seconds...\n", buburNumber, delay)
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("chef is too tired while making bubur order number %d", buburNumber)
		} else if rnd <= 5 {
			msg = fmt.Sprintf("bubur order number %d is out of ingredient", buburNumber)
		} else {
			success = true
			msg = fmt.Sprintf("bubur order number %d is ready", buburNumber)
		}

		b := &BuburOrder{
			success:     success,
			buburNumber: buburNumber,
			message:     msg,
		}

		return b

	}

	return &BuburOrder{
		success:     false,
		buburNumber: buburNumber,
		message:     "target completed",
	}
}

func buburConsumer(orderData chan BuburOrder, quitTrigger chan chan error) {
	for i := range orderData {
		if i.buburNumber <= NumberOfBuburs {
			if i.success {
				color.Green(i.message)
				color.Green("order number %d ready to serve", i.buburNumber)
			} else {
				color.Red(i.message)
				color.Red("The customer is really mad!")
			}
		} else {
			color.Cyan("done making bubur")
			color.Cyan(i.message)
			break
		}
	}

	color.Cyan("-----------------")
	color.Cyan("Done for the day.")

	color.Cyan("We made %d buburs, but failed to make %d, with %d attempts in total.", buburMade, buburFailed, total)

	switch {
	case buburFailed > 9:
		color.Red("It was an awful day...")
	case buburFailed >= 6:
		color.Red("It was not a very good day...")
	case buburFailed >= 4:
		color.Yellow("It was an okay day....")
	case buburFailed >= 2:
		color.Yellow("It was a pretty good day!")
	default:
		color.Green("It was a great day!")
	}

	ch := make(chan error)
	quitTrigger <- ch
	close(ch)
}
