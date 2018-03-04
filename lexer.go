// Package lexer provides a simple scanner for handling lexical analysis.
// the implementation is based on Rob Pike's talk.
// http://www.youtube.com/watch?v=HxaD_trXwRE
package lexer

type ItemType int

// ItemType identifies the type of lex items.
const (
	ItemError ItemType = iota
	ItemIdent          // identifier
	ItemQuote          // quote "
	ItemText           // plain text
)

// Item represents a token returned from the scanner.
type Item struct {
	Type  ItemType // Type, such as ItemNumber.
	Value string   // Value, such as "23.2".
}

// StateFn represents the state of the scanner
// as a function that returns the next state.
type StateFn func(*Lexer) StateFn

// Lexer holds the state of the scanner.
type Lexer struct {
	name  string    // used only for error reports.
	input string    // the string being scanned.
	start int       // start position of this item.
	pos   int       // current position in the input.
	width int       // width of last rune read from input.
	items chan Item //  channel of scanned items.
}

// Run lexes the input by executing state functions
// until the state is nil
func (l *Lexer) Run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items) // No more tokens will be delivered.
}

// Emit passes an item back to the client.
func (l *Lexer) Emit(t ItemType) {
	l.items <- Item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}
