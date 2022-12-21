package main

import (
	"bufio"
	"os"
	"testing"
)
func loadTestFile(path string) [][]byte {
	file, err := os.Open(path)
	if err != nil {
		panic("Couldn't open " + path)
	}
	defer file.Close()

	var words [][]byte
	reader := bufio.NewReader(file)
	for {
		if line, err := reader.ReadBytes(byte('\n')); err != nil {
			break
		} else {
			if len(line) > 0 {
				words = append(words, line[:len(line)-1])
			}
		}
	}
	return words
}



func BenchmarkInsertOrderedWords(b *testing.B) {
	// create a new ART tree
	// generate a large number of test keys to insert into the tree
	// keys := generateTestKeys()
	words := loadTestFile("dict/words")
  
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tree := NewTree()

		for i := 0; i < 500; i++ {
			word := words[i]
			tree.Insert(word, word)
			// do something with the first 100 words
		}

	}


}


func BenchmarkInsertRandomWords(b *testing.B) {
	// create a new ART tree
	// generate a large number of test keys to insert into the tree
	// keys := generateTestKeys()
	words := loadTestFile("dict/uuids")
  
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tree := NewTree()

		for i := 0; i < 500; i++ {
			// killer dando erro
			word := words[i]
			tree.Insert(word, word)
			// do something with the first 100 words
		}

	}

}