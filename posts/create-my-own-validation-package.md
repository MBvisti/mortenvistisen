In my on-going effort to make [Grafto](https://github.com/mbv-labs/grafto) simple yet able to handle whatever usecase a SaaS business might have, I ended up writing my own validation package.

Initially I used the `go-playground/validator` pacakge and it worked _great_ until I had to customize it. At the time of writing I can't for the life of me remember what I had to customize, but something had to be customized!

After trying to stay with the library and get the customization I needed, I decided to look under the hood. It had already been borrowing a bit of how they construct their error messages. For the unfamiliar, when a field fails some validation(s) you get an error message like:

```go
var baseErrMsg = "Field: '%s' with Value: '%v' has Error(s): validation failed due to '%v'"
```

Essentially show you all the validation rules this field failed making error handling in the UI a breeze.

Also, I wanted to get something that looked similar to how Laravel does validation which closely resembles how you would do a map in Go. 

```go
validation := map[string][]ValidationRules{
    "ID": NotNilType
}
```

The ability to write my own custom types and validation rules that adhere to a simple interface is very freeing.
