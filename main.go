package main

import (
	"fmt"
	"log"
	"os"

	config "github.com/KrishKoria/GatorConfig"
)

type state struct {
    Config *config.Config   
}

type command struct {
    Name string
    Args []string
}

type commands struct {
    Handlers map[string]func(*state, command) error
}

func main()  {
    cfg, err := config.Read()
    if err != nil {
        log.Fatalf("Error reading config file: %v", err)
    }
    
    s := &state{Config: cfg}
   
    cmds := &commands{}
    cmds.register("login", handlerLogin)    

    if len(os.Args) < 2 {
        fmt.Println("Error: No command provided")
        os.Exit(1)
    }

    cmd := command{
        Name: os.Args[1],
        Args: os.Args[2:],
    }

    err = cmds.run(s, cmd)
    if err != nil {
        log.Fatalf("Error running command: %v", err)
    }
}

func handlerLogin(s *state, cmd command) error {
    if len(cmd.Args) == 0 {
        return fmt.Errorf("enter a username")
    }
    username := cmd.Args[0]
    err:= s.Config.SetUser(username)
    if err != nil {
        return fmt.Errorf("error setting user: %v", err)
    }
    fmt.Printf("User set to %s\n", username)
    return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
    if c.Handlers == nil {
        c.Handlers = make(map[string]func(*state, command) error)
    }
    c.Handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error{
    if handler, exists := c.Handlers[cmd.Name]; exists {
        return handler(s, cmd)
    }
    return fmt.Errorf("command not found: %s", cmd.Name)
}
