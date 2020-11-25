package main

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerarPropuesta (total int)([]int64){
	var propuesta []int64
	for i:=0;i<=total/3;i++{
		rand.Seed(time.Now().UnixNano())
		lilprop:=rand.Perm(3)
		if i==total/3{
			sobra:=total%3
			for j,num :=range lilprop{
				if j==sobra{
					break
				}
				propuesta=append(propuesta,int64(num))
			}
		}else{
			for _,num :=range lilprop{
				propuesta=append(propuesta,int64(num))
			}
		}
	}
	return propuesta
}

func main() {
	asd:=GenerarPropuesta(54)
	fmt.Println(asd)
	fmt.Println(len(asd))
}