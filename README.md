scrappy

Usage: scrappy <url> <element>

scrappy will read all the text and save it to a local sqlite db in ~/.scrappy/scrappy_notes.db

the table has the schema

notes (
        url TEXT PRIMARY KEY,
        content TEXT
)
