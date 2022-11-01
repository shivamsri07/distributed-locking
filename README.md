# distributed-locking
Distributed lock implementation in Golang

To test out, a dummy queue was created with 4 messages in it. 3 consumers (client) were run as a go routine, and they were synchronized by locking using redis. 

## Output<img width="526" alt="Screenshot 2022-11-02 at 12 42 16 AM" src="https://user-images.githubusercontent.com/12581295/199319478-65314819-621d-4d37-ad59-745bb674fe30.png">
