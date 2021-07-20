![](https://pbs.twimg.com/profile_banners/1396003728/1617122494/1500x500)

# Busy Gopher
Busy Beaver Visualizer written in Golang.


Just type
`ssh busygopher.csn3rd.com`
in your terminal and press enter.

## Cool Visualizations

2 Card 3 Symbol:
Transition for card 0 on input 0: 111
Transition for card 0 on input 1: 201
Transition for card 0 on input 2: 11-1
Transition for card 1 on input 0: 200
Transition for card 1 on input 1: 211
Transition for card 1 on input 2: 101

![](https://github.com/csn3rd/Busy-Gopher-Game/blob/master/2card3symbol.png)

4 Card 2 Symbol:
Transition for card 0 on input 0: 111
Transition for card 0 on input 1: 101
Transition for card 1 on input 0: 100
Transition for card 1 on input 1: 002
Transition for card 2 on input 0: 11-1
Transition for card 2 on input 1: 103
Transition for card 3 on input 0: 113
Transition for card 3 on input 1: 010

![](https://github.com/csn3rd/Busy-Gopher-Game/blob/master/4card2symbol.png)

## About

During the 2020-2021 spring quarter, I took CSCI 169 - Programming Languages. This course covered programming languages from different paradigms such as object-oriented programming, functional programming, and more. The final project was to pick a programming language that I was interested in, learn and understand how it works (syntax, compilation, execution, etc), and then make something of my choice with a couple hundred lines or more. 

At first, I wasn’t sure what I wanted to make. So, I decided to think back about some interesting CS problems I had encountered in previous classes. During the 2020-2021 winter quarter, I took CSCI 162 - Computational Complexity. It was a really interesting course that covered a lot of theoretical concepts in computer science. I learned about lots of different complexity classes and continued learning about different computational models following CSCI 161 - Theory of Automata and Languages. At the end of the quarter, everyone in the class was assigned a research topic to study and present. I was assigned the busy beaver problem and function. I ended up reading lots of papers on the subject and watched several lectures. It was quite interesting to see how a problem can seem very simple but can not be computed or solved. I was super surprised to see how large the numbers can get even for small cases. It was really cool to see how quickly the function grows as we increase the number of cards and symbols allowed. I then made a 30-minute presentation about the history, the actual problem, current findings, and current frontiers. Link to the presentation: [Busy Beaver Presentations.pdf](https://github.com/csn3rd/Busy-Gopher-Game/blob/master/Busy%20Beaver%20Presentation.pdf).

Since the problem is undecidable or non-computable, there are no algorithms or efficient methods for determining the most optimal machines. Instead, many researchers have written and run code on supercomputers that attempt to brute force and simulate all possible machines (up to some determined thresholds). I didn’t want to write some code to brute force this problem so I decided to look at some other interesting parts of the problem.

When I was researching the problem, I had seen some interesting visualizations and I thought it would be cool to see visualizations for any machine, not just the most optimal ones. As a result, I decided to make a Busy Beaver visualizer. My program would simulate user-defined machines and display them. First, the user inputs the number of cards and the number of symbols. Then, they define each card and transition. My code would simulate the behavior of the machine one step at a time. Since the machines can run forever, I have limited the tape to 50 cells and the overall simulation to 200 shifts.

This project would be pretty cool for people who have not heard about this problem or would like to learn more about this problem. It is a great way to visualize how different cards and symbols may lead to different results and configurations.

## Process

For my project, I made Busy Gopher, a Busy Beaver visualizer using Go. While researching ways to display the game state and some different libraries I could use, I discovered this cool game of snake that is played in the terminal over SSH and it is written in Go. The link: https://github.com/zachlatta/sshtron. Looking into the code, I found that they used the library “golang.org/x/crypto/ssh”. It seemed pretty tough to use this library and the documentation was very long and complicated. I looked into some other SSH libraries and found “github.com/gliderlabs/ssh” which was much simpler to use for setting up ssh servers. Following the examples in the repository and reading some of the documentation, I was able to create a program that starts up a server on my local machine on a port I designate that I can connect to from another terminal window. For input and output to the terminal, I looked into the “golang.org/x/term” library. I knew that I wanted to implement different colors like the sshtron game so I researched into ANSI color codes and learned how to print colored text within the terminal.

After setting up the server output code, I implemented the busy beaver game itself. Starting with the tape, I decided that I would use a String array to model the tape with the symbols on it. I created functions to initialize the tape, print the tape, and calculate the score of the tape. These weren’t particularly hard. Initializing the tape was just making a new array of size 50 and initializing each position to symbol “0”. Printing the tape is just going through the tape and printing the character at each position. I wanted to highlight the nonzero positions on the tape so I used the ANSI color codes I had found previously. To score the tape, each position is checked and the score is incremented if the position is nonzero.

Moving onto the cards or states, it seemed much harder because cards have different fields like an id and all the different transitions. Plus, each transition has four parts: input, overwrite, shift, and card transition values. I decided to create a struct card with int variable id and then three int arrays to store the transitions. I could have also made a struct for the transitions. My program would be pretty much the same with a few minor adjustments. To debug my code, I added a printCard function that prints the cards to the server logs so I could make sure that it was storing the id and transitions correctly. With a working card struct, I could then write a function to simulate a step given the current position on the tape and the current card. This was not too hard but required careful attention to make sure that the transition is followed exactly as expected. I tested making some different cards and simulating them to make sure that the transitions were correctly followed.

The final step was to actually make a game function to take user inputs. This was probably the most time-consuming part because there were lots of type conversions and error handling required to not crash the server and set up the cards correctly. After a lot of coding, testing, and debugging, I was able to have a working MVP that takes in user inputs and prints out the simulated tape. To finish the project, I experimented with different colors and ASCII art designs to make it less boring.

After submitting the project, I decided that I wanted to deploy the server publicly so others could try it without needing to install all the modules to run the code. Since I already had a VPS, I decided to use the same VPS for this project. Just running my project as it was, I could connect by sending the command “ssh ipaddress -p 2222”. However, the IP address is long and hard to remember. Since I already have a domain name for my website, I found out that I could redirect a subdomain to the IP address by adding an extra DNS A record. After some further research on how to do that, I could now connect using the command “ssh busygopher.csn3rd.com -p 2222”. I didn’t want to have to specify the port number every time so I looked into how SSH works. By default, SSH uses port 22 to connect. However, if I set the port in my code to 22, my program can not listen to that port because it is already in use. It turns out that port 22 is set up for remote login and access to the VPS. After doing a lot of research on how to change the port, I learned that this is specified within the ssh_config and sshd_config files. Changing some of the port settings within those two files solves the issue. Now, anyone can connect with just “ssh busygopher.csn3rd.com”.

Overall, this project was a lot of fun and I learned a lot about Go, the ssh and terminal libraries, setting up ssh configurations, and redirecting (sub)domain names to IP addresses using DNS records.
