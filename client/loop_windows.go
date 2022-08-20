package main

import "time"

func start() {

	var external_port uint16 = 0
	for {

		send_syn()

		request(&external_port)

		for i := 0; i < get_conf("interval").(int); i++ {
			send_syn()
			time.Sleep(1 * time.Second)
		}
	}
}
