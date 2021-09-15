tasksOnHand := make([]*nonBlockQueue2.TaskRequest, 0)

		counter := 0

		for (*checking == val) {	//are we done
			if counter < config.BlockSize{
				taskObj := queue.Dequeue()
				if (taskObj != nil){ 	//non empty queue
					tasksOnHand = append(tasksOnHand, taskObj)
					fmt.Println("here")
					counter +=1
					continue
				} else if counter >= 0{
					break
				} else {
					cond.Wait()
				}
			}else{
				break
			}
		}

		if !(*checking == val){
			break
		}



//big for loop checking done
	//have less tasks thatn blocksize
		//attemp dequeue
		//if dequeue is non empty
			//append
		//else if dequeue empty, counter > 1
			//break
		//else
			//cond.Wait()