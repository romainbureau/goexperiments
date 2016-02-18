package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func explodeLine(line string, separator string) []string {
	explodedLine := strings.Split(line, separator)

	return explodedLine
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		explodedLine := explodeLine(line, ";")
		/*
			var buffer bytes.Buffer
			for _, value := range explodedLine {
				buffer.WriteString(value)
			}
			fmt.Println(buffer.String())
		*/
		newLine := strings.Join(explodedLine, ",")
		fmt.Println(line + " -> " + newLine)
	}
}
