package model

// ColValue represents an XML element that holds a column's value.
// It may carry an isNull attribute or a CDATA text value.
type ColValue struct {
	IsNull string `xml:"isNull,attr"`
	Value  string `xml:",chardata"`
}

// Column represents a single <column> element inside an <operation>.
type Column struct {
	Name        string    `xml:"name,attr"`
	Index       int       `xml:"index,attr"`
	BeforeValue *ColValue `xml:"before-value"`
	AfterValue  *ColValue `xml:"after-value"`
}

// Operation represents a single <operation> element from the CDC XML feed.
type Operation struct {
	Table    string   `xml:"table,attr"`
	Type     string   `xml:"type,attr"`
	OpType   string   `xml:"opType,attr"`
	TxInd    string   `xml:"txInd,attr"`
	Ts       string   `xml:"ts,attr"`
	NumCols  int      `xml:"numCols,attr"`
	Position string   `xml:"position,attr"`
	Columns  []Column `xml:"column"`
}
