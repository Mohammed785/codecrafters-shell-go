package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)
var builtin []string = []string{"exit","echo","type","pwd"}

func read() (string,[]string,error){
	fmt.Fprint(os.Stdout, "$ ")
	inp,err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err!=nil{
		return "",nil,err
	}
	cmd,args_s,_ := strings.Cut(strings.TrimSpace(inp)," ")
	args := strings.Split(args_s, " ")
	return cmd,args,nil
}
func handle_commands(cmd string,args []string){
	switch cmd {
	case "exit":
		status,_ := strconv.Atoi(args[0])
		os.Exit(status)
	case "echo":
		fmt.Fprintf(os.Stdout,"%s\n",strings.Join(args, " "))
	case "pwd":
		pwd,_:=os.Getwd()
		fmt.Fprintln(os.Stdout,pwd)
	case "cd":
		if args[0]=="~"{
			args[0],_=os.UserHomeDir()
		}
		if err:=os.Chdir(args[0]);err!=nil{
			fmt.Fprintf(os.Stdout,"cd: %s: No such file or directory\n",args[0])
		}
	case "type":
		PATH,_:=os.LookupEnv("PATH")
		if slices.Contains(builtin,args[0]){
			fmt.Fprintf(os.Stdout,"%s is a shell builtin\n",args[0])
		}else{
			for _,p:=range strings.Split(PATH,":"){
				fp := filepath.Join(p,args[0])
				if _,err:=os.Stat(fp);err==nil{
					fmt.Fprintf(os.Stdout,"%s is %s\n",args[0],fp)
					return
				}
			}
			fmt.Fprintf(os.Stdout,"%s: not found\n",args[0])
		}
	default:
		command:=exec.Command(cmd,args...)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err:= command.Run();err!=nil&&errors.Is(err,exec.ErrNotFound){
			fmt.Fprintf(os.Stdout,"%v: command not found\n",strings.TrimSpace(cmd))
		}
	}
}
func eval(){
	cmd,args,err:= read()
	if err!=nil{
		fmt.Fprintln(os.Stderr,err.Error())
		return
	}
	handle_commands(cmd,args)
}

func main() {
	for {
		eval()
	}
}
