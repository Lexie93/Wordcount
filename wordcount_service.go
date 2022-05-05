package wordcount_service

import (
	"fmt"
	"strings"
	"unicode"
)

// <key, value> pairs
type Couple []struct{
	Word string
	Num int
}

// Count service for RPC
type Count int

// Split file text in words
func split(str string) []string {
	f:= func (c rune) bool {
		return !unicode.IsNumber(c) && !unicode.IsLetter(c)
	}
	words:= strings.FieldsFunc(str, f)
	return words
}

// Sum same words occurrences in the file, reducing RPC call size
func presum(str []string) Couple {
	c:= make(Couple, len(str))
	last:=0
	if len(str)==0 {
		return nil
	}
	c[0].Word=str[0]
	c[0].Num=1
	
	for i:=1; i<len(str); i++ {
		j:=0
		for j=0 ; j<=last ; j++{
			if str[i]==c[j].Word {
				break
			}
		}
		if str[i]==c[j].Word {
			c[j].Num++
		} else {
			last++
			c[last].Num++
			c[last].Word=str[i]
		}
	}
	return c[:last+1]
}

// Determine total occurrences of words in all files
func sum(in Couple) Couple {
	var c Couple
	if len(in)==0 {
		return nil
	}
	c=append(c, in[0])
	j:=0
	for i:=1; i<len(in); i++ {
		j=0
		for j=0 ; j<len(c) ; j++{
			if in[i].Word==c[j].Word {
				j++
				break
			}
		}
		if in[i].Word==c[j-1].Word {
			c[j-1].Num+=in[i].Num
		} else {
			c=append(c, in[i])
		}
	}
	return c
}

// Map fase
func (t *Count) Map(text string, c *Couple) error {
	f:= split(strings.ToLower(text))
	*c=presum(f)
	return nil
}

// Reduce fase
func (t *Count) Reduce(in Couple, out *Couple) error {
	*out=sum(in)
	fmt.Println(*out)
	return nil
}
