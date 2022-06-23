package brainfucklibrary

import (
	"fmt"
	"io"
	"log"
)

type mydata struct {
	Program    byte         // indivisual command
	Memory     map[int]byte // data
	PMemory    int
	Input      io.Reader
	Output     io.Writer
	Shouldomit bool   // this is to inform the function to store commands in []byte for looping operation
	Programs   []byte // slice of commands(need when there is a loop)
	Idx        int    // to keep thetrack of ']' bracket
	Insptr     int    // to keep a track of intsruction number
}

func GetData() *mydata {
	return &mydata{}
}

//This function is executed when there is a loop
//Does the same operation as Run but in loop
func RunLoop(mytype *mydata) error {
	var PInstr = 0

	// Scans the loop to store the begining and end position  of the loop  to jump to and fro
	// Also returs an error if the loop is not balanced
	loopPositions, err := scanLoops(mytype)
	if err != nil {
		return fmt.Errorf("unable to run program: %w", err)
	}
	loops := NewStack()
	if mytype.Memory[mytype.PMemory] != 0 {
		for PInstr < len(mytype.Programs) {

			switch mytype.Programs[PInstr] {
			case '+':
				if mytype.Memory[mytype.PMemory] == 255 {
					mytype.Memory[mytype.PMemory] = 0
				} else {
					mytype.Memory[mytype.PMemory]++
				}
			case '-':
				if mytype.Memory[mytype.PMemory] == 0 {
					mytype.Memory[mytype.PMemory] = 255
				} else {
					mytype.Memory[mytype.PMemory]--
				}
			case '>':
				if len(mytype.Memory) == mytype.PMemory {
					mytype.Memory[mytype.PMemory+1] = 0
				}
				mytype.PMemory++
			case '<':
				if mytype.PMemory != 0 {
					mytype.PMemory--
				}
			case '.':
				_, err := mytype.Output.Write([]byte{mytype.Memory[mytype.PMemory]})
				if err != nil {
					return err
				}
			case ',':
				buf := make([]byte, 1)
				_, err := mytype.Input.Read(buf)

				switch err {
				case nil:
					mytype.Memory[mytype.PMemory] = buf[0]
				case io.EOF:
					mytype.Memory[mytype.PMemory] = 0
				default:
					return fmt.Errorf("unable to read from Input: %w", err)
				}

			case '[':
				if mytype.Memory[mytype.PMemory] == 0 {
					PInstr = loopPositions[PInstr]
				} else {
					loops.Push(PInstr)
				}
			case ']':
				PInstr, err = loops.Pop()
				if err != nil {
					return err
				}
				continue
			}

			PInstr++
		}
	}
	mytype.Programs = []byte{}

	return nil
}

// 1) 8 bit cell size is chosen. Programs that use extra range are likely to be slow, since storing the value {\displaystyle n}n into a cell requires {\displaystyle O(n)}O(n) time as a cell's value may only be changed by incrementing and decrementing.
// 2) Run executes a given command and perform the required action.
// 3) for each non matching '[' or ']' error is sent
// 4) error is updated for each command execution
// 5)+ Increment the value at data pointer and wraps to 0 if value has reached 255.
// 6)- Decrement the value at data pointer and wraps to 255 if value has reached 0.
// 7)> Increment (shift one right) the byte and adds an extra byte if its equal to length.
// 8)< Decrement (shift one left) the byte if idx not zero.
// 9). Output the byte at the data pointer.
// 10), Accept one byte(8 bit) of input, storing its value in the byte at the data pointer.
// 11)[ If the input command is '[' and current value is not zero, then the successive command is stored in a []byte so as to execute them in a loop till the ']' corresponding to '[' is entered
// If the byte at the data pointer is zero, then instead of moving the instruction pointer forward to the next command, skip till the command after the matching ] command.
// 12)] If the byte at the data pointer is nonzero, then instead of moving the instruction pointer forward to the next command, jump it back to the command after the matching [ command.`

func (mytype *mydata) Run() error {

	switch mytype.Program {
	case '+':
		if mytype.Shouldomit {
			if mytype.Memory[mytype.PMemory] != 0 {
				mytype.Programs = append(mytype.Programs, mytype.Program)
			}
		} else {
			if mytype.Memory[mytype.PMemory] == 255 {
				mytype.Memory[mytype.PMemory] = 0
			} else {
				mytype.Memory[mytype.PMemory]++
			}
			mytype.Insptr++
		}
	case '-':
		if mytype.Shouldomit {
			if mytype.Memory[mytype.PMemory] != 0 {
				mytype.Programs = append(mytype.Programs, mytype.Program)
			}
		} else {
			if mytype.Memory[mytype.PMemory] == 0 {
				mytype.Memory[mytype.PMemory] = 255
			} else {
				mytype.Memory[mytype.PMemory]--
			}
			mytype.Insptr++
		}
	case '>':
		if mytype.Shouldomit {
			if mytype.Memory[mytype.PMemory] != 0 {
				mytype.Programs = append(mytype.Programs, mytype.Program)
			}
		} else {
			if len(mytype.Memory) == mytype.PMemory {
				mytype.Memory[mytype.PMemory+1] = 0
			}
			mytype.PMemory++
			mytype.Insptr++
		}
	case '<':
		if mytype.Shouldomit {
			if mytype.Memory[mytype.PMemory] != 0 {
				mytype.Programs = append(mytype.Programs, mytype.Program)
			}
		} else {
			if mytype.PMemory != 0 {
				mytype.PMemory--
			} else {
				log.Println("unable to decrease before starting index")
			}
			mytype.Insptr++
		}
	case '.':
		if mytype.Shouldomit {
			if mytype.Memory[mytype.PMemory] != 0 {
				mytype.Programs = append(mytype.Programs, mytype.Program)
			}
		} else {
			_, err := mytype.Output.Write([]byte{mytype.Memory[mytype.PMemory]})
			if err != nil {
				log.Fatalln("unable to write:", err)
				return err
			}
			mytype.Insptr++
		}
	case ',':
		if mytype.Shouldomit {
			if mytype.Memory[mytype.PMemory] != 0 {
				mytype.Programs = append(mytype.Programs, mytype.Program)
			}
		} else {
			buf := make([]byte, 1)
			_, err := mytype.Input.Read(buf)

			switch err {
			case nil:
				mytype.Memory[mytype.PMemory] = buf[0]
			case io.EOF:
				mytype.Memory[mytype.PMemory] = 0
			default:
				log.Fatalln("unable to read from Input:", err)
			}
			mytype.Insptr++
		}
	case '[':
		mytype.Idx++
		mytype.Programs = append(mytype.Programs, mytype.Program)
		mytype.Shouldomit = true
		return fmt.Errorf("unbalanced loop end at index %d", mytype.Insptr)
	case ']':
		mytype.Programs = append(mytype.Programs, mytype.Program)
		mytype.Idx--
		if mytype.Idx <= 0 && len(mytype.Programs) != 0 {
			mytype.Shouldomit = false
			// Execute the commands inside the loop
			err := RunLoop(mytype)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ScanLoops creates a dictionary which maps start positions of loops to their respective end positions
func scanLoops(mytype *mydata) (map[int]int, error) {
	loopPositions := make(map[int]int)
	s := NewStack()
	for i, instr := range mytype.Programs {
		switch instr {
		case '[':
			s.Push(i)
		case ']':
			start, err := s.Pop()
			if err != nil {
				return nil, fmt.Errorf("unbalanced loop end at index %d", i+mytype.Insptr)
			}

			loopPositions[start] = i
		}
	}

	start, err := s.Pop()

	if err == nil {
		return nil, fmt.Errorf("unbalanced loop start at index %d", start+mytype.Insptr)
	}

	return loopPositions, nil
}
