## System req

### Client 
- Publish valid Message 
- When Broadcast Recieved Write it to the log with timestamp
- Can Join at any time 
- Clients can drop at anytime 
### Server 
- Broadcast all valid messages With logical timestamp. To all curent clients
- Validate messages
  - Must be less than 128 characters
  - Must be UTF-8
- Broadcast when a new client joins with current timestamp
- Broadcast when a new client drops with current timestamp


