package db


var createSearchTable = `
CREATE TABLE IF NOT EXISTS search (
   id bool PRIMARY KEY DEFAULT TRUE, 
   token text, CONSTRAINT search_id_uni CHECK (id)
);
`
type Search struct {
	Id bool `db:"id"`
	Token  string `db:"token"`
}
