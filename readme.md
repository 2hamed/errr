# Errr (The Rich Structured Error Package missing from Go)

The errr package was created to fill the gap between error handling and error reporting. There are countless times that we need more details about why and how an error is happening, and since the default error doesn't allow us to include structured data we end up with error logs littered everywhere throughout our code base.  
You might suggest that we could use fmt.Errorf() to wrap the errors with more details. But that doesn't quite cut it. Using fmt.Errorf produces an unstructured output which is not suitable for querying and searching in for example GCP Logs.  
And why not just use logs? Because as I said it will lead to a several layers of redundant logs. If you check now, we have several lines of logs for just a single error.  
And the benfit of this package is that it's just error, so it can be easily returned as a regular error and the higher up functions could just wrap or return them as is.

Some example usage:

```go
func main() {
    err := layer1()
    if err != nil{
        panic(err)
    }
}

func layer1() error {
    err := layer2()
    
    return errr.WithValue("caller", "layer1").Wrap(err)
}

func layer2() error {
    err := some_error

    return errr.WithValue("user_id", user_id).Wrap(err)
}
```