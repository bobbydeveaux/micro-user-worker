package main

import (
	"encoding/json"
	"github.com/nats-io/go-nats"
	"log"
	"os"

	"time"
)

// need to get this from ENV, because GitHub public project will expose this. Oops.
const passphrase string = "fbac-FJfxeMQCzXBPqrIY8Hhk"

type person struct {
	Id          int64
	Name        string
	Valid       bool
	Jwt         string
	AccessToken accessToken
}

type accessToken struct {
	Value  string
	Expiry int
}

type jwtAuth struct {
	User   person
	Expiry int
}

type jwtToken struct {
	Value string
}

func main() {
	//NATS_HOST = nats://localhost:4222
	nc, _ := nats.Connect(os.Getenv("NATS_HOST"))
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer ec.Close()

	var p person
	ec.QueueSubscribe("user.login", "job_workers", func(msg *nats.Msg) {
		log.Printf("logging in: %s\n", msg.Data)
		err := json.Unmarshal(msg.Data, &p)
		if err != nil {
			log.Println(err.Error())
		}

		err = ec.Request("user.auth", p, &p, 100*time.Millisecond)
		if err != nil {
			if nc.LastError() != nil {
				log.Println("Error in Request: %v\n", nc.LastError())
			}
			log.Println("Error in Request: %v\n", err)
		} else {
			log.Printf("Published [%s] : '%s'\n", "user.auth", p.Name)
			log.Printf("Received Auth [%v] : '%s'\n", p.Id, p.Name)
		}

		if p.Valid == false {
			ec.Publish(msg.Reply, p)
		}

		var at accessToken
		err = ec.Request("auth.generateaccesstoken", p, &at, 100*time.Millisecond)
		if err != nil {
			if nc.LastError() != nil {
				log.Println("Error in Request: %v\n", nc.LastError())
			}
			log.Println("Error in Request: %v\n", err)
			p.AccessToken.Value = "error with auth.generateaccesstoken service"
		} else {
			log.Printf("Published [%s] : '%s'\n", "auth.generateaccesstoken", p.Name)
			log.Printf("Received Token [%v] : '%s'\n", at.Value)
			p.AccessToken = at
		}

		var jwt jwtToken
		err = ec.Request("auth.jwt", p, &jwt, 100*time.Millisecond)
		if err != nil {
			if nc.LastError() != nil {
				log.Println("Error in Request: %v\n", nc.LastError())
			}
			log.Println("Error in Request: %v\n", err)
			p.Jwt = "error with auth.jwt service"
		} else {
			log.Printf("Published Access Token[%s] : '%s'\n", "auth.jwt", p.AccessToken)
			log.Printf("Received JWT [%v] : '%s'\n", jwt.Value)
			p.Jwt = jwt.Value
		}

		if err != nil {
			log.Println(err.Error())
		}
		ec.Publish(msg.Reply, p)
	})

	ec.QueueSubscribe("user.auth", "job_workers", func(msg *nats.Msg) {
		log.Printf("Authenticating: %s\n", msg.Data)
		err := json.Unmarshal(msg.Data, &p)
		if err != nil {
			log.Println(err.Error())
		}

		// @TODO check against database
		p.Id = int64(time.Now().UnixNano())
		p.Valid = false
		if p.Name == "bobby" {
			p.Valid = true
		}

		if err != nil {
			log.Println(err.Error())
		}
		ec.Publish(msg.Reply, p)
	})

	ec.QueueSubscribe("user.getuser", "job_workers", func(msg *nats.Msg) {
		log.Printf("Finding user: %s\n", msg.Data)
		err := json.Unmarshal(msg.Data, &p)
		if err != nil {
			log.Println(err.Error())
		}

		// @TODO check against database
		p.Id = int64(time.Now().UnixNano())
		p.Name = "Username"

		if err != nil {
			log.Println(err.Error())
		}
		ec.Publish(msg.Reply, p)
	})

	ec.QueueSubscribe("user.createuser", "job_workers", func(msg *nats.Msg) {
		log.Printf("Creating user: %s\n", msg.Data)
		err := json.Unmarshal(msg.Data, &p)
		if err != nil {
			log.Println(err.Error())
		}

		// @TODO save against database
		p.Id = int64(time.Now().UnixNano())
		p.Name = "Username"
		p.Valid = true
		if err != nil {
			log.Println(err.Error())
		}
		ec.Publish(msg.Reply, p)
	})

	select {}
}
