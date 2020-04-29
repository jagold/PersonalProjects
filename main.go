package main

//Import necessary packages
import (
	"bufio" //Buffer i/o
	"fmt"   //standard output
	"github.com/detectlanguage/detectlanguage-go"
	"github.com/hackebrot/turtle" //emoji package
	"net"                         //networking
	"os"                          //operating system
	"regexp"                      //regular expression
	"strings"                     //string handler
)

//Hardcode host, and port number for demonstration purposes
const(
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)


//User enters their name
var user string
var client *detectlanguage.Client

func main() {
	client = detectlanguage.New("40d888a8b8c438fda468faa013bb9e4d")
	fmt.Print("Enter username: ") //Prompt user for their name
	fmt.Scan(&user)
	var client_or_server string
	fmt.Print("Do you want to be the client or server? ") //Ask about client vs. server
	fmt.Scan(&client_or_server)
	if (client_or_server == "server") { //If server
		l, err := net.Listen(CONN_TYPE,CONN_HOST+":"+CONN_PORT) //Listen at predefined ip address
		defer l.Close() //Defer closing the connection until client has disconnected
		if (err != nil){ //Error in connection
			fmt.Println("error listening for client")
			os.Exit(1)
		}
		fmt.Println("Listening...")

		for  {
			conn,err := l.Accept() //Accept connection

			if (err != nil) {
				fmt.Println("error accepting connection")
				os.Exit(2)
			}
			println("Client connected")
			//Begin goroutine to receive and send messages concurrently
			go receiveMessage(conn)
			sendMessages(conn)
		}
	}else { //User is client
		conn, err := net.Dial(CONN_TYPE,"127.0.0.1:3333") //Go to local ip address
		defer conn.Close() //defer the close
		if (err != nil) {
			fmt.Println("error connecting to server")
			os.Exit(4)
		}

		fmt.Println("Connected to Server")
		//Concurrently send and receive messages
		 go sendMessages(conn)
		 receiveMessage(conn)

	}

}
/*
sendMessages creates a buffer and gets user input
append an newline character to the input, and send it
over the connection. Recursively call send messages again
 */
func sendMessages(conn net.Conn)  {
	reader := bufio.NewReader(os.Stdin)

	message_to_send,_ := reader.ReadString('\n')

	fmt.Fprintf(conn,user+": "+message_to_send+ "\n")
	fmt.Println()
	sendMessages(conn)
}


/*
receiveMessage takes the connection as a parameter
and prints messages that are sent over the connection
 */
func receiveMessage(conn net.Conn){
	is_emoji := make(chan string) //create a channel
	message,_ := bufio.NewReader(conn).ReadString('\n') //Read the incoming message

	checkURL(message) //Check if the url is safe
	go checkLanguage(message)
	go translateToEmoji(is_emoji,message) //Translate emojis


	fmt.Println(<-is_emoji) //Print the message output from the is_emoji channel
	receiveMessage(conn) //Recursively call it again
}

/*
Uses regular expression to check whether
the url links to a secure site
 */
func checkURL(message string){
	//Matches http, but not https
	//The 's' marks a secure website
	match,_ := regexp.MatchString("(http)[^s]","\""+message+"\"hra")
	if match { //Matches http, give a warning
		fmt.Println("WARNING: URL MAY BE INSECURE")
	}
}


/*
translateToEmoji takes the incoming messages, and
checks if any of the words can be emojis, it will
change the words to their emoji counterparts if it can
 */
func translateToEmoji(emoji_chan chan string, message string){
	message_array := strings.SplitN(message," ",-1) //Turn the message into a string array
	for index,word := range message_array{
		//iterate over array and check whether the word is a valid emoji
		emojis := turtle.Emojis[strings.TrimSpace(word)]
		if emojis != nil {
			//If it is a valid emoji, replace the word in the array
			//with the emoji counterpart
			message_array[index] = turtle.Emojis[strings.TrimSpace(word)].Char
		}
	}
	//Join the array into a string
	message = strings.Join(message_array," ")
	//Send the message over a channel for printing
	emoji_chan <- message
}


func checkLanguage(message string){
	language,_ := client.DetectCode(message)
	if language != "en"  {
		fmt.Println("Detected language :",language)
	}
}