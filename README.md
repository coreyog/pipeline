# Pipeline
To install:
```
go get github.com/coreyog/pipeline
```

## Usage

See the [Examples](github.com/coreyog/pipeline/examples/test.go):

```
pipe := pipeline.New()
pipe.PushFunc(os.Open)
pipe.PushFunc(ioutil.ReadAll)
results, _ := pipe.Call("example.txt")
fmt.Println(string(results[0].([]byte)))
```