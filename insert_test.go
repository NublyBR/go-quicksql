package quicksql

import (
	"bytes"
	"testing"
)

func TestInsertSimple(t *testing.T) {
	var (
		buf    = bytes.NewBuffer(nil)
		expect = "REPLACE INTO `sample` (`id`, `name`, `bytes`) VALUES\n\t(1, \"Demo\", 0x48656c6c6f2c20576f726c6421),\n\t(2, \"Test\", 0x54657374);\n\n"
	)

	type sample struct {
		ID    int
		Name  string
		Bytes []byte
	}

	ins := NewInsert(buf, "sample", sample{}).
		Replace().
		Add(sample{
			ID:    1,
			Name:  "Demo",
			Bytes: []byte("Hello, World!"),
		}).
		Add(sample{
			ID:    2,
			Name:  "Test",
			Bytes: []byte("Test"),
		}).
		Flush()

	if err := ins.Err(); err != nil {
		t.Fatal(err)
	}

	if expect != buf.String() {
		t.Errorf("expected buffer to equal %q, got %q", expect, buf.String())
	}

}

func TestInsertSplit(t *testing.T) {
	var (
		buf    = bytes.NewBuffer(nil)
		expect = "INSERT INTO `split` (`number`) VALUES\n\t(0),\n\t(1),\n\t(2),\n\t(3);\n\nINSERT INTO `split` (`number`) VALUES\n\t(4),\n\t(5),\n\t(6),\n\t(7);\n\nINSERT INTO `split` (`number`) VALUES\n\t(8),\n\t(9);\n\n"
	)

	ins := NewInsert(buf, "split", "number").Every(4)

	for i := 0; i < 10; i++ {
		ins.Add(i)
	}

	ins.Flush()

	if err := ins.Err(); err != nil {
		t.Fatal(err)
	}

	if expect != buf.String() {
		t.Errorf("expected buffer to equal %q, got %q", expect, buf.String())
	}
}
