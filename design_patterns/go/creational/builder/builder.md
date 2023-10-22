Background:
* Some objects/structs are simple and can be created in a simple constuctor call or factory function calls or initialising partial fields.
* If an object has 10 fields, then you would force user to make lot of decisions on user end within a single expression.
* With Builder, we make construction process, a multi-stage process; piecewise.

* Builder pattern provides an API for object construction step by step.
e.g. In Go, string builder

**Problem Statement:**
*Write a web server that serves HTML* (It also serves JS, but let's keep it out of scope)
* So we have to build HTML out of ordinary text.
```go
func main() {
	hello := "hello"
	sb := strings.Builder{}
	sb.WriteString("<p>")
	sb.WriteString(hello)
	sb.WriteString("</p>")
	fmt.Println(sb.String())
	
	words := []string{"hello", "world"}
	sb.Reset()
	// step 1
	sb.WriteString("<ul>")
	for _, v := range words {
		// step 2
		sb.WriteString("<li>")
		sb.WriteString(v)
		sb.WriteString("</li>")
    }
	// step 3
    sb.WriteString("</ul>")
	fmt.Println(sb.String())
}
```
* The above code is a bit too much work - we know where we want li and /li
* So lets put it into structures and methods which are more handy to use.




