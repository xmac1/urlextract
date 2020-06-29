package main

import (
	"flag"
	"fmt"

	"github.com/globalsign/mgo/bson"

	"github.com/tealeg/xlsx"

	"github.com/globalsign/mgo"
)

var (
	mongo     string
	test      bool
	inputFile string
)

func main() {
	flag.StringVar(&mongo, "mongo", "", "mongo addr")
	flag.StringVar(&inputFile, "i", "", "input file name")
	flag.BoolVar(&test, "T", false, "test mode")
	flag.Parse()

	session, err := mgo.Dial(mongo)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	file, err := xlsx.OpenFile(inputFile)
	if err != nil {
		panic(err)
	}

	for _, sheet := range file.Sheets {
		for i, row := range sheet.Rows {
			if i == 0 {
				continue
			}
			postIDCell := row.Cells[1]
			postID := postIDCell.String()
			if postID == "" {
				continue
			}

			var result struct {
				UID       string `json:"uid" bson:"uid"`
				PostID    string `json:"post_id" bson:"post_id"`
				Content   string `json:"content" bson:"content"`
				Extention string `json:"extention" bson:"extention"`
			}

			col := session.DB("basepost").C("postmedia")
			err := col.Find(bson.M{
				"post_id": postID,
			}).One(&result)
			if err != nil {
				fmt.Println("get failed", err)
				continue
			}
			ac := row.AddCell()
			b, _ := bson.Marshal(result)
			ac.SetString(string(b))
		}
	}
}
