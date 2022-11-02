# distributed-locking
Distributed lock implementation in Golang

To test out, a dummy queue was created with 4 messages in it. 3 consumers (client) were run as a go routine, and they were synchronized by locks using redis. 

## Output

<img width="499" alt="Screenshot 2022-11-02 at 6 27 05 PM" src="https://user-images.githubusercontent.com/12581295/199495861-7f24ee9b-a7cf-4b9a-ba45-f071c17dc6bf.png">


