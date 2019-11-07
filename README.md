# Pipeline
To install:
```
go get github.com/coreyog/pipeline
```

## Usage

See the [Examples](https://github.com/coreyog/pipeline/blob/master/examples/main.go):

```
pipe := pipeline.New()
pipe.PushFunc(os.Open)
pipe.PushFunc(ioutil.ReadAll)
results, _ := pipe.Call("example.txt")
fmt.Println(string(results[0].([]byte)))
```