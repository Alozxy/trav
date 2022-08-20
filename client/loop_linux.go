package main

import "time"

func start() {
	go syn_loop()

	var external_port uint16 = 0
	for {

		request(&external_port)
		time.Sleep(time.Duration(get_conf("interval").(int)) * time.Second)
	}
}
