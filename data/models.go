package data

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/mattn/go-sqlite3"
	db2 "github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
	"github.com/upper/db/v4/adapter/postgresql"
	"github.com/upper/db/v4/adapter/sqlite"
)

var db *sql.DB
var upper db2.Session

type Models struct {
	// any models inserted here (and in the New function)
	// are easily accessible throughout the entire application

	//RememberToken RememberToken
	//Tokens        Token
	//Users         User
}

func New(databasePool *sql.DB) Models {
	db = databasePool

	switch os.Getenv("DATABASE_TYPE") {
	case "mysql", "mariadb":
		upper, _ = mysql.New(databasePool)

	case "postgres", "postgresql":
		upper, _ = postgresql.New(databasePool)

	case "sqlite", "sqlite3":
		upper, _ = CreateSqliteUpperConnection(os.Getenv("DATABASE_NAME"))
		// Does not work because the databasePool is not
		// a proper one in the case of sqlite3
		// ie sqlite3 returns a driver.Conn type and
		// not the driver.Connector type as in the case
		// for postgres.
		// THerefore use the open statement as in the statement
		// above???
		//upper, _ = sqlite.New(databasePool)

	default:
		// do nothing
	}

	return Models{

		//RememberToken: RememberToken{},
		//Users:         User{},
		//Tokens:        Token{},
	}
}

func CreateSqliteUpperConnection(path string) (db2.Session, error) {
	var settings = sqlite.ConnectionURL{
		Database: path,
	}
	db, err := sqlite.Open(settings)

	return db, err
}

/* version for integer ID's
func getInsertID(i db2.ID) int {
	idType := fmt.Sprintf("%T", i) // returns the type of i
	if idType == "int64" {
		return int(i.(int64))
	}

	return i.(int)
}
*/

// UUID ID version - return string
func getInsertID(i db2.ID) string {
	idType := fmt.Sprintf("%T", i) // returns the type of i
	if idType != "string" {
		fmt.Println(">>> Invalid ID type - expecting string/uuid")
		return ""
	}

	return i.(string)
}

// take a struct and convert it into a []byte
func EncodeToBytes(p interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Compress compresses the []byte into another []byte
func Compress(s []byte) []byte {
	zipbuf := bytes.Buffer{}
	zipped := gzip.NewWriter(&zipbuf)
	zipped.Write(s)
	zipped.Close()
	return zipbuf.Bytes()
}

// Decompress uncompresses the []byte into another []byte
func Decompress(s []byte) ([]byte, error) {
	rdr, _ := gzip.NewReader(bytes.NewReader(s))
	data, err := ioutil.ReadAll(rdr)
	defer rdr.Close()
	if err != nil {
		return data, err
	}
	return data, nil
}

// itob returns an 8-byte big endian representation of v.
func Itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
