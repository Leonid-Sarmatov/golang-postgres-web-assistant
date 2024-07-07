package generator

import (
	//"fmt"
	//"log"
	//"unicode"

	//"time"
	//"reflect"
	//"database/sql"
	//"os"
	//"os/exec"
	//"plugin"
	//"text/template"

	_ "github.com/lib/pq"
)

var Port = ":8080"

var StructTeg = "`json:\"rows\"`"

var StartCode = 
`package main

import (
	"net/http"
	"database/sql"
	_ "github.com/lib/pq"
)

type Message struct {
	Rows sql.Rows `+StructTeg+`
}

type TableStruct struct {
`

var EndCode = `
}

func TableStruct() *TableStruct {
	return &TableStruct{}
}

func main() {
	http.HandleFunc("/parser", func(w http.ResponseWriter, r *http.Request) {
        var message Message
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&message)
		if err != nil {
			http.Error(w, "ERROR", http.StatusBadRequest)
		}

		jsonResponse, err := json.Marshal(message)
		if err != nil {
			http.Error(w, "ERROR", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
    })

    http.ListenAndServe("`+Port+`", nil)
}
`

/*func CreateAndStartService(table_name string, port string, structMap map[string]string) error {
	Port = port

	for fieldName, fieldType := range structMap {

	}

	// Парсинг кода из строки
	t := template.Must(template.New("ParceStruct").Parse(ParceStruct))
	f, err := os.Create("./cmd/serv_"+table_name+"/generated.go")
	if err != nil {
		return err
	}
	defer f.Close()
	err = t.Execute(f, nil)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "run", "./cmd/app2/generated.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error: %v, output: %s", err, output)
		return err
	}

	return nil

	/*cmd := exec.Command("go", "build", "-buildmode=plugin", "-o",
		"./cmd/app/generated.so", "./cmd/app/generated.go",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("failed to compile plugin: %v, output: %s", err, output)
		return nil, err
	}

	plug, err := plugin.Open("./cmd/app/generated.so")
	if err != nil {
		return nil, err
	}

	sym, err := plug.Lookup("NewResult")
	if err != nil {
		return nil, err
	}

	newResultFunc := sym.(func() interface{})
	resultVal := reflect.ValueOf(newResultFunc()).Elem()

	fmt.Println(resultVal)
}*/