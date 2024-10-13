package read

import "bufio"

func Readln(reader *bufio.Reader) (string, bool, error) {
	line, isPrefix, err := reader.ReadLine()

	return string(line), isPrefix, err
}
