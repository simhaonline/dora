// Package dora TODO: package docs
package dora

import (
	"fmt"
	"github.com/bradford-hamilton/dora/pkg/ast"
	"github.com/bradford-hamilton/dora/pkg/lexer"
	"github.com/bradford-hamilton/dora/pkg/parser"
	"github.com/spf13/cast"
)

// Client represents a dora client. The client holds things like a copy of the input, the tree (the
// parsed AST representation built with Go types), the user's query & parsed version of the query, and
// a query result. Client exposes public methods which access this underlying data.
type Client struct {
	input       []rune
	tree        *ast.RootNode
	query       []rune
	parsedQuery []queryToken
	result      string
}

// NewFromString takes a string, creates a lexer, creates a parser from the lexer,
// and parses the json into an AST. Methods on the Client give access to private
// data like the AST held inside.
func NewFromString(jsonStr string) (*Client, error) {
	l := lexer.New(jsonStr)
	p := parser.New(l)
	tree, err := p.ParseJSON()
	if err != nil {
		return nil, err
	}
	return &Client{tree: &tree, input: l.Input}, nil
}

// NewFromBytes takes a slice of bytes, converts it to a string, then returns `NewFromString`, passing in the JSON string.
func NewFromBytes(bytes []byte) (*Client, error) {
	return NewFromString(string(bytes))
}

func (c *Client) preflight(query string) error {
	if err := c.prepareQuery(query, c.tree.Type); err != nil {
		return err
	}
	if err := c.executeQuery(); err != nil {
		return err
	}

	return nil
}

// Get takes a dora query, prepares and validates it, executes the query, and returns the result or an error.
func (c *Client) Get(query string) (string, error) {
	if err := c.preflight(query); err != nil {
		return "", err
	}
	return c.result, nil
}

func (c *Client) Set(cursor string, val string) error {
	if err := c.preflight(cursor); err != nil {
		return fmt.Errorf("error not able to walk path provided by cursor: %w", err)
	}

	if c.result == val {
		return nil
	}

	c.result = val

	return nil
}

// inspired by viper.Get<T>() implementation
func (c *Client) GetString(cursor string) string {
	result, err := c.Get(cursor)
	if err != nil {
		return ""
	}

	return cast.ToString(result)
}

// TODO: implement GetBool()
// TODO: implement GetFloat64(), JSON's only number type
// TODO: implement GetArray() or maybe Slice{} ?
// TODO: implement GetObject() or maybe Struct{} ?
// TODO: implement GetNull() ? for completeness ?
