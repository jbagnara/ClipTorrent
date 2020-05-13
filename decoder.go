package main

import (
	"fmt"
	"strconv"
	"bytes"
	"os"
	"io/ioutil"
)

type decoder struct{	//decoder state
	s	[]byte		//source file
	pos int			//position of current read byte
}

func new(s []byte) *decoder{	//decoder state constructor
	d := decoder{s, 0}
	return &d
}

func decodeNext(d *decoder) (interface{}, int) {	//returns next decoded value and updates slice pointer
	char := d.s[d.pos]
	var ret interface{}

	for char != 'i' && char != 'l' && char!= 'd' && char != ':' && d.pos+1<len(d.s) {
		d.pos++
		char = d.s[d.pos]
	}

	switch(char){
		case 'i':
			ret = parseInt(d)
			break
		case ':':
			ret = parseString(d)
			break
		case 'l':
			ret = parseList(d)
			break
	}

	d.s=d.s[d.pos:]		//Removes parsed entries from byte slice
	d.pos=0;
	err := 0

	if ret!=nil{
		err = 1
	}

	return ret, err
}

func parseInt(d *decoder) int {		//accepts byte slice with pos pointing to 'i', returns number up to 'e'
	start := d.pos+1
	end := bytes.IndexByte(d.s, 'e')

	substr := string(d.s[start:end])
	d.pos = end+1

	ret, err := strconv.Atoi(substr)

	if err != nil {
		fmt.Printf("Invalid integer '%s'\n", substr)
		os.Exit(1)
	}

	return ret;
}

func parseString(d *decoder) string {	//accepts byte slice with pos pointing to ':"
	substrl := string(d.s[:d.pos])
	l, err := strconv.Atoi(substrl)

	if err != nil {
		fmt.Printf("Invalid length '%s'\n", substrl)
		os.Exit(1)
	}
	l++

	ret := string(d.s[d.pos+1:d.pos+l])

	d.pos = d.pos+l
	return ret
}

func parseList(d *decoder) []interface{} {
	d.pos++
	list := make([]interface{}, listSize(d));
	char := d.s[d.pos]
	i := 0

	for char != 'e' {
		list[i], _ = decodeNext(d)
		i++
		char = d.s[d.pos]
	}

	return list

}

func listSize(d *decoder) int{	//accepts byte array with pos at list 'l', returns size of list
	i := 1			//counter of current 'depth', 0 when end of initial list reached
	j := d.pos
	size := 0

	for i > 0{
		if d.s[j] == 'i' || d.s[j] == 'l' || d.s[j] == ':' ||  d.s[j] == 'd' {	//increase size each value
			if d.s[j] != ':'{
				i++
			}
			if d.s[j] == 'l'{	//if list, subtract size of sublist
				oldpos := d.pos
				d.pos = j+1
				size = size - listSize(d)
				d.pos = oldpos
			}
			size++
		} else if d.s[j] == 'e' {
			i--
		}
		j++
	}
	return size

}

func main(){
	if len(os.Args) < 2 {
		fmt.Printf("No file provided\n")
		os.Exit(1)
	}

	byteArr, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		fmt.Printf("Invalid file: %s\n", os.Args[1])
		os.Exit(1)
	}

	d := new(byteArr)
	err2 := 1
	var ret interface{}

	for err2==1 {
		ret, err2 = decodeNext(d);
		fmt.Println(ret)
	}
}
