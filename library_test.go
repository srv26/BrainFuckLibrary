package brainfucklibrary

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

const (
	helloworld = "++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++."
	WorldLine  = "+]"
)

func Test_PositiveBF(t *testing.T) {
	r := strings.NewReader(helloworld)
	buf := make([]byte, 1)
	output := new(bytes.Buffer)
	mydata := GetData()
	mydata.Input = os.Stdin
	mydata.Output = output
	mydata.Memory = make(map[int]byte, 0)
	var err error
	var got string
	for {
		_, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		mydata.Program = buf[0]
		err = mydata.Run()
	}
	if err == nil {
		got = string(output.Bytes())
	}
	want := "Hello World!\n"
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func Test_Error(t *testing.T) {
	r := strings.NewReader(WorldLine)
	buf := make([]byte, 1)
	output := new(bytes.Buffer)
	mydata := GetData()
	mydata.Input = os.Stdin
	mydata.Output = output
	mydata.Idx = 0
	mydata.Memory = make(map[int]byte, 0)
	var got error
	want := "unable to run program: unbalanced loop end at index 1"
	for {
		_, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		mydata.Program = buf[0]
		got = mydata.Run()
	}
	if got.Error() != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

}
