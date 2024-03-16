I recently got into a discussion on twitter; I suggested that root level `internal` directories doesn't make much sense when you're creating applications that are __not__ a library. It's conventional wisdom in the Go community, to suggest putting any code you do not want to expose to other applications inside the `internal` directory. And it's true, if you're creating something that is meant to important into other peoples code. The truth is, however, that most people are likely not working on other libraries but rather applications (apis, web servers, etc). Take this layout, for example:

```
- app
    - cmd
	- app
	    - main.go
    - internal
    - services
    - repository
    .git
    go.mod
    go.sum
```

A standard repository on github, gitlab etc. 

What does `internal` protect against here? Literally any other package in this application will be on the same level as the `internal` directory, all will be able to import from it. The only situation, where this would not hold true is if you were to take this entire code and copy it into another repository. And, honestly, when was the last time you did that?

So, yes, it's technically correct (best kind) that using root-level `internal` directories provides encapsulation (which is good). But, in practice, you don't really get any of the benefits of `internal`. If we were to add a view package to our application, and wanted to start grouping different views together in a sub-package, an `internal` directory suddenly starts to make sense. Let me illustrate:

```
- app
    ...omitted
    - views
	layouts.go
	- home
	    home.go
	- about
	    about.go
```

## When does internal make sense?

If we had some layout components (base, dashboard etc) in the `layouts.go` file that aren't exported, `home.go` and `about.go` wouldn't be able to import it as both are no longer part of the view package. If we export it, __any__ package in the application can import it, obviously not desirable. But, we can add an `internal`directory, export the layout components and encapsulate them from the rest of the application, like so:

```
- app
    ...omitted
    - views
	- internal
	    - layouts
		base.go
		dashboard.go
	- home
	    home.go
	- about
	    about.go
```

This is kind of a pet peeve of mine; it isn't really that important now that we have gotten modules (pre modules was one of the reasons for internal which I will get into in a bit) and as long as you stay consistent about what goes into `internal`, it's all good. I personally prefer to use `pkg`, as that to me at least better communicate what the directory is for. And I've seen quite a few Go applications in my time now, they either use `internal` or `pkg` and most often contain the same kind of code. Auxiliary packages that are needed for the application as a whole but does not contain core business logic. Code to interact with `aws`, wrappers around a library to send mails, a [queue](https://github.com/mbv-labs/grafto/tree/master/pkg/queue) etc. To me, for non-library code, this is the most important thing to get right, what do you expect to find in the directory not as much the same (again, you don't get any of the benefits from internal if you place it in the root).

## Why do we have internal?

To answer that question, we have to go quite a big back. All the way back to the time before modules, where the creators of Go thought it a great idea to force everyone into their (Google's) preferred way of structuring code. You'd be forced into something like this:

```
- $GOPATH
	- src
	- github.com
	    - user
		- repo
		    - cmd
		    - internal
		    - pkg
```

Before Go 1.4, you could only have local and global components which could lead to frustrating situations where you'd like to have code that is only used within the module its defined, and not part of the public repository. And in Go 1.4, the solution to this was introduced with the introduction of the `internal` directory. The Go compiler will verify that the package doing the import is within the same module as the package being imported.

With the above structure, using `internal` as the default makes a ton of sense as it's quite easy to accidentally expose something that otherwise should've been kept private. But with the introduction of modules, this doesn't really apply. Reading `internal` makes me thing of something that's internal to a __package__ not the application as a hole.


## Does this really matter

No. 

I tend to pick this "fight" just because I think its a bit funny how some people will go out of their way to argue that `internal` is the only correct way to do it.

As argued above, if you're building something that is not meant to be imported into other peoples code, this is good advice. If you're building an application, consider using `pkg`. But most importantly, be consistent. We spend too much time in software focusing on right and wrong as black and white, but it all depends. If the code that I'm arguing should go into a `pkg` directory feels more natural for __you and your team__ to be in an `internal` directory, then by all means, do that.

Just be consistent.
