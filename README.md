# distributed-locking
Distributed lock implementation in Golang

To test out, a dummy queue was created with 4 messages in it. 3 consumers (client) were run as a go routine, and they were synchronized by locks using redis. 

## Output

<img width="467" alt="Screenshot 2022-11-02 at 12 14 02 PM" src="https://user-images.githubusercontent.com/12581295/199419092-8d669c3d-24cb-4757-8f42-800aff5d5cb7.png">

