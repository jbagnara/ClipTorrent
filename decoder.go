package main

import (
	"fmt"
	"strconv"
	"bytes"
	"os"
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

func main(){

	byteArr := []byte{'6', 'i', '1', '0', '3', '9', '4', 'e', '3', ':', 'f', 'f', 'f', 't', 'g', '4', 'i', '7', '4', '3', 'e', 'u'}
	d := new(byteArr)
	err := 1
	var ret interface{}

	for err==1 {
		ret, err = decodeNext(d);
		fmt.Println(ret)
	}
}
