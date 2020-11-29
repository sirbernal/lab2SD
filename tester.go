package main

import (
	"fmt"
	"math/rand"
	"time"
	"reflect"
	"os"
)
var activos=[]int{1,2}
func GenerarPropuestaNueva (total int, conectados int)([]int64){
	var propuesta []int64
	for i:=0;i<=total/conectados;i++{
		rand.Seed(time.Now().UnixNano())
		lilprop:=rand.Perm(conectados)
		if i==total/conectados{
			sobra:=total%conectados
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
	switch conectados{
	case 1:
		if activos[0]==0{
			return propuesta
		}else{
			for i,_:= range propuesta{
				propuesta[i]=int64(activos[0])
			}
			return propuesta
		}
	case 2:
		if reflect.DeepEqual(activos,[]int{0,2}){
			for i,j:=range propuesta{
				if j==1{
					propuesta[i]=2
				}
			}
			return propuesta
		}else if reflect.DeepEqual(activos,[]int{1,2}){
			for i,j:=range propuesta{
				if j==0{
					propuesta[i]=2
				}
			}
			return propuesta
		}else{
			return propuesta
		}
	default:
		return propuesta
	}
}

func main() {
	asd:=GenerarPropuestaNueva(54,2)
	fmt.Println(asd)
	fmt.Println(len(asd))
	
	file, err := os.Open("datanode1/lel.pdf_0")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} 
	defer file.Close()
	fmt.Println(reflect.TypeOf(file))
}