# GGPO-Go - GGPO Port Into Go
GGPO-Go is a port of the GGPO rollback netcode library into Go. Currently unfinished. 

## Usage 
General usage would be best explained by looking at the code in the example folder.

But to view the example, follow the following steps: 

- Clone the repository. 
- Travel to the example folder 
- enter:
    `go run . <local_port> <num_players>  <local|remote_ip:port> <local|remote_ip:port> <current_player>`

An example usage would be to open one command line input with the following command 
`go run . 7000 2 local 127.0.0.1:7001 1`
and another with this command 
`go run . 7001 2 127.0.0.1:7000 local 2`

