package klondike

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestKlondike_Draw(t *testing.T) {
	k := NewKlondike()
	k.Init()

	assert.Equal(t, len(k.table[stock]), 24)
	assert.Equal(t, len(k.table[waste]), 0)
	err := k.Draw()
	assert.Nil(t, err)
	assert.Equal(t, len(k.table[stock]), 23)
	assert.Equal(t, len(k.table[waste]), 1)
	assert.Equal(t, k.table[waste][0].Open, true)

	k.table[stock] = k.table[stock][:0]
	assert.Equal(t, len(k.table[stock]), 0)
	err = k.Draw()
	assert.Equal(t, err, StockIsEmpty)
}

func TestKlondike_WasteToStock(t *testing.T) {
	k := NewKlondike()
	k.Init()

	err := k.WasteToStock()
	assert.Equal(t, err, StockIsNotEmpty)

	k.Draw()
	k.table[stock] = Cards{}

	assert.Equal(t, len(k.table[stock]), 0)
	assert.Equal(t, len(k.table[waste]), 1)
	err = k.WasteToStock()
	assert.Nil(t, err)
	assert.Equal(t, len(k.table[stock]), 1)
	assert.Equal(t, len(k.table[waste]), 0)
	assert.Equal(t, k.table[stock][0].Open, false)
}

func TestKlondike_MoveToFoundation(t *testing.T) {
	k := NewKlondike()
	k.Init()

	var err error

	err = k.MoveToFoundation(&Position{stock, 0}, &Position{foundation1, 0})
	assert.Equal(t, err, InvalidMovement)

	err = k.MoveToFoundation(&Position{waste, 0}, &Position{column1, 0})
	assert.Equal(t, err, InvalidMovement)

	err = k.MoveToFoundation(&Position{column2, 0}, &Position{foundation1, 0})
	assert.Equal(t, err, CardIsNotLastCol)

	k.table[column2][1].Num = 2
	err = k.MoveToFoundation(&Position{column2, 1}, &Position{foundation1, 0})
	assert.Equal(t, err, CanNotPutInFoundation)

	k.table[column2][1].Num = 1
	assert.Equal(t, len(k.table[column2]), 2)
	assert.Equal(t, len(k.table[foundation1]), 0)
	err = k.MoveToFoundation(&Position{column2, 1}, &Position{foundation1, 0})
	assert.Equal(t, len(k.table[column2]), 1, "場札から組み札へ移動できた")
	assert.Equal(t, len(k.table[foundation1]), 1, "場札から組み札へ移動できた")
	assert.Equal(t, k.table[column2][0].Open, true, "移動後に一番上のカードをめくる")

	k.Draw()
	k.table[waste][0].Num = 2
	k.table[waste][0].Suit = k.table[foundation1][0].Suit
	log.Println(k)

	err = k.MoveToFoundation(&Position{waste, k.LastCol(waste)}, &Position{foundation1, 0})
	assert.Nil(t, err)
	assert.Equal(t, len(k.table[waste]), 0, "捨て札から組み札へ移動できた")
	assert.Equal(t, len(k.table[foundation1]), 2, "捨て札から組み札へ移動できた")
}

func TestKlondike_MoveToColumn(t *testing.T) {
	k := NewKlondike()
	k.Init()

	var err error

	err = k.MoveToColumn(&Position{stock, 0}, &Position{column1, 0})
	assert.Equal(t, err, InvalidMovement)

	err = k.MoveToColumn(&Position{waste, 0}, &Position{foundation1, 0})
	assert.Equal(t, err, InvalidMovement)

	err = k.MoveToColumn(&Position{column2, 0}, &Position{column1, 0})
	assert.Equal(t, err, CardIsNotOpen)

	k.Draw()
	k.table[waste][0].Num = k.table[column1][0].Num
	err = k.MoveToColumn(&Position{waste, 0}, &Position{column1, 0})
	assert.Equal(t, err, CanNotPutInTableau)

	k.table[waste][0].Num = k.table[column1][0].Num + 1
	k.table[waste][0].Suit = k.table[column1][0].Suit
	err = k.MoveToColumn(&Position{waste, 0}, &Position{column1, 0})
	assert.Equal(t, err, CanNotPutInTableau)

	k.table[waste][0].Suit = Hearts
	k.table[column1][0].Suit = Clubs
	k.table[waste][0].Num = 5
	k.table[column1][0].Num = 6
	assert.Equal(t, len(k.table[waste]), 1)
	assert.Equal(t, len(k.table[column1]), 1)
	err = k.MoveToColumn(&Position{waste, 0}, &Position{column1, 0})
	assert.Nil(t, err, "捨て札から場札へ移動できた")
	assert.Equal(t, len(k.table[waste]), 0)
	assert.Equal(t, len(k.table[column1]), 2)

	k.table[column2][1].Num = 7
	k.table[column2][1].Suit = Diamonds
	err = k.MoveToColumn(&Position{column1, 0}, &Position{column2, 1})
	assert.Nil(t, err, "場札から場札へ複数枚移動できた")
	assert.Equal(t, len(k.table[column1]), 0)
	assert.Equal(t, len(k.table[column2]), 4)

	log.Println(k)
}
