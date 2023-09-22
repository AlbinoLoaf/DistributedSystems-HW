# a)
 we use channels to transfer our data, which are strings, and our metadata, which are integers, between 2 threads that act as a client and server.

# b) 
our implementation uses threads, this is unrealistic because you would practically never have a client-server relationship that runs on the same computer. this eliminates the possibility of data being lost

# c)
- Message sequencing:
We ensure that the entire process of the three-way handshake are in sequence of a step is not in sequence the request does not continue.

- Buffering sequencing and  Reorder Logic:
to this issue adding a buffer to the server would be a very effective idea this way we can "hold" a (or multiple) message that comes in the wrong order until the correct one arrives 

# d)
Our implementation do account for it by ensuring data integrity with sequences and checking that the received answer or request is indeed the expected step. We could add logic to respond when it is out of sequence or to save the data and insert it into sequence of discard it if it is not possible

# e)
It ensures proper synchronization and data integrity and let us establish a reliable connection 
