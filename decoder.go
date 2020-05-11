package main

import (
	"fmt"
	"strconv"
	"bytes"
)

type decoder struct{	//decoder state
	s	[]byte		//source file
	pos int			//position of current read byte
}

func new(s []byte) *decoder{	//decoder state constructor
	d := decoder{s, 0}
	return &d
}

func decodeNext(d decoder) interface{} {	//decode next value
	char := d.s[d.pos]
	fmt.Println(char)
	var ret interface{}

	switch(char){
		case 'i':
			ret = parseInt(d)
			break
	}
	d.pos++
	return ret
}

func parseInt(d decoder) int {
	start := d.pos+1
	end := bytes.IndexByte(d.s[d.pos:], 'e')
	substr := string(d.s[start:end])
	d.pos = end+1

	ret, _ := strconv.Atoi(substr)
	return ret;
}

func main(){
	byteArr := []byte{'i', '1', '0', '3', '9', '4', 'e'}
	d := new(byteArr)
	fmt.Println(decodeNext(*d))
}
