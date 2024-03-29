package mapreduce

import (
	"encoding/json"
	"os"
	"sort"
)

// ByKey implements sort.Interface for []KeyValue based on
// the Key field.
type ByKey []KeyValue

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTask int, // which reduce task this is
	outFile string, // write the output here
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	//
	// doReduce manages one reduce task: it should read the intermediate
	// files for the task, sort the intermediate key/value pairs by key,
	// call the user-defined reduce function (reduceF) for each key, and
	// write reduceF's output to disk.
	//
	// You'll need to read one intermediate file from each map task;
	// reduceName(jobName, m, reduceTask) yields the file
	// name from map task m.
	//
	// Your doMap() encoded the key/value pairs in the intermediate
	// files, so you will need to decode them. If you used JSON, you can
	// read and decode by creating a decoder and repeatedly calling
	// .Decode(&kv) on it until it returns an error.
	//
	// You may find the first example in the golang sort package
	// documentation useful.
	//
	// reduceF() is the application's reduce function. You should
	// call it once per distinct key, with a slice of all the values
	// for that key. reduceF() returns the reduced value for that key.
	//
	// You should write the reduce output as JSON encoded KeyValue
	// objects to the file named outFile. We require you to use JSON
	// because that is what the merger than combines the output
	// from all the reduce tasks expects. There is nothing special about
	// JSON -- it is just the marshalling format we chose to use. Your
	// output code will look something like this:
	//
	// enc := json.NewEncoder(file)
	// for key := ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()
	//
	// Your code here (Part I).
	//
	var kv KeyValue
	values := make([]KeyValue, 0)
	for i := 0; i < nMap; i++ {
		rName := reduceName(jobName, i, reduceTask)
		f, err := os.Open(rName)
		if err != nil {
			panic(err)
		}
		dec := json.NewDecoder(f)
		for err := dec.Decode(&kv); err == nil; err = dec.Decode(&kv) {
			values = append(values, kv)
		}
		f.Close()
		output, err := os.Create(outFile)
		if err != nil {
			panic(err)
		}
		enc := json.NewEncoder(output)
		sort.Sort(ByKey(values))
		kvs := make(map[string][]string)
		for _, curr := range values {
			kvs[curr.Key] = append(kvs[curr.Key], curr.Value)
		}
		for key := range kvs {
			enc.Encode(KeyValue{key, reduceF(key, kvs[key])})
		}
		output.Close()
	}

}
