## Purpose
The goal of this project is to produce a network-forming node that is capable of
searching both itself and the network for a given file.

## Design

### Overview

Each FileSearch node can send and receive messages to and from other FileSearch
nodes in its network.

Each node implements an HTTP client and server. Nodes connect to each other via
HTTP. Each connection returns a pair of channels, one to send messages with and
the other to receive responses.

For a node to join a network it must first connect to an existing node and
authenticate. Once a node has successfully connected to a node, its target node
will form a connection in return. Successful connections result in a node's
address added to its `neighboringNodes` array. Each time an edge is added, the
node performs an `append()` on the slice and resets the `neighboringNodes`
field. This strikes me as inefficient, but I'm not yet sure how best to treat a
list of indeterminate length. 

When a node receives a search message, it saves the sender's address along with
a channel that's used to stop the search. The node then forwards the search
request to all neighboring nodes while kicking off a search of its own file
system. Before a node begins searching its own file system, it adds an entry to
its `searchCache` so searches for the same file will not happen twice. The
`searchCache` contains a channel that can be used to stop the search if needed
(This is WIP because the project uses Go's builtin `WalkDir` function, which
[does not utilize goroutines](https://golang.org/src/path/filepath/path.go?s=12067:12112#L388).)

If the file is found, the node that finds encrypts the file with the public key
and sends a RelayMessage to the node that requested the search. If a node
receives a RelayMessage, it forwards it to the address in the message.

### Usage

A user will have access to a CLI (still WIP) that can send a file's hash and a
public key to the address of a node.

There are 2 types of messages at each node's disposal:

1. Search Messages
    - These contain the hash of the requested file and a public RSA key. When a
      node receives a SearchMessage it simultaneously begins a search of its own
      file system and broadcasts the SearchMessage out to neighboring nodes.
      ```
        type SearchMessage struct {
            FileHash      string
            PublicKey     string
            SenderAddress string
        }
      ```

2. Relay Messages
    - These messages are sent once a file is found. Each contains the encrypted
      file data as well as the http address they are to be sent to. When a node
      receives a relay message, it simply forwards it to the address in the
      message.
    ```
    type RelayMessage struct {
        FileData string
        Address  string
    }
    ```

## Discussion
My original plan was to use docker’s Libchan library as transport between nodes.
Libchan uses Google's [SPDY](https://www.chromium.org/spdy/spdy-whitepaper)
protocol as transport between nodes. The protocol's focus on low latency,
multiple streams over a single connection seemed a natural fit with Go's focus
on concurrency. But Libchan relies on an implementation of SPDY that was not as
straightforward as I would’ve liked. The protocol itself being experimental, I
thought it best to fall back onto http as a way to send messages between nodes.
For this reason, communication between file search nodes is a bit rough, and
relies on a pair of request and response channels per connection. This is
something that I alluded to in the “Discussion” portion of the repository.

Since the original design was meant to expect a push from the server (SPDY
allows for this) and had to be changed to accommodate a request/response
communication model, the code that records the network edges (channels that
allow message-passing between nodes) is not 100% there. But I hope the code
outlines how I’d implement this given the time.

Right now the node's relay functionality is WIP. The currently implementation
stores the edges of the network graph on a node by node basis, so each node has
knowledge only of its neighboring nodes. When a RelayMessage is received, a node
could, in theory, lookup the node that first requested a search to begin with.
This still needs to be fleshed out.

Messages are serialized into wrappers called `Unpackable` messages. The reason
for choosing a struct to wrap the messages instead of a map was that a struct
was easier to unserialize into a `NodeMessage`. The switch-case statement that
handle messages can then type match on the message instead of having to use an
enum or worse, a string comparison, to decide what handler method to call.

## Security Considerations
The hash function to use is still up in the air. For testing I've been using
SHA256 hash function, but im not clear on the advantages of this hash function
vs others (apart from it using a 256 bit signature).

All node activity happens over HTTP. Right now none of the nodes need or require
SSL certificates. Nodes should authenticate with each other before successfully
connecting, and all nodes in the network should eventually pass messages over
TLS connections.

The file system should be restricted to a directory that a given node allows to
be searched. Access to the file system could be further restricted to certain
permissions.

When the file is found, its encrypted using the public key sent by the user.
When the file is finally returned to the user it needs to be verified. One way
to do this is to decrypt the file, hash it, and compared the output hash to the
input hash. More consideration needs to be put in here.