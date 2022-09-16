<strong>2022 update</strong>. At the time of writing this update, it's been over a year since I 
published this article. At the end of this article I promise a part two, which I must 
admit, will probably never come for a couple of reasons. I no longer use this 
structure (I still think this article gives people some help to get going, but would 
like to update it at one point. I suggest you checkout this [repo](https://github.com/bnkamalesh/goapp). 
Secondly, in regards to integration testing, I've just relased a full-length article 
about it which you can find [here](https://mortenvistisen.com/posts/how-to-do-integration-test-with-go).

Go is an amazing language. It’s simple, easy to reason about and gives you much tooling right out of the box. However, when I started with Go I struggled to figure out how to structure my applications in a way that didn’t use “enterprisey” approaches. This is my take on structuring Golang applications that start out simple, but with the flexibility to grow, and what I would have wanted when I started with Go.

## Disclaimer

This article is based on the following people (link to their Twitter):
- [Mat Ryer](https://twitter.com/matryer)
- [Kat Zien](https://twitter.com/kasiazien)
- [Jon Calhoun](https://twitter.com/joncalhoun)
- [Robert C. Martin](https://twitter.com/unclebobmartin)

All of the above-mentioned people have articles and tutorials that are amazing on their own, which you should definitely checkout.

If you just want to see the code for this article’s project, check out the [Github repository](https://github.com/MBvisti/weight-tracker-article).

## Project & tools

It is always frustrating to learn a new language and having to create the same to-do app over and over again. So we are going to make an application that allows people to create an account, where they can track their weight and calculate their nutritional requirements (inspired by weight increase during lockdown). Feel free to add to this application with features you see fit.

Now, I’m assuming some familiarity with Golang and programming in general, as this is a practical example of how to structure a Go application. We will be using the following:

- Gin Web Framework
- Go version 1.15+
- Postgresql
- Golang-migrate
- Your favorite IDE (could be Goland or VS code — I highly suggest Goland, it’s amazing)

## Elements of a robust and scalable structure

Spend some time working in software development and you quickly learn how programs sometimes can be brittle and other times rigid. Therefore, I look for these three things in my application structure:

- Testability
- Readability
- Adaptability

Writing programs in a way that also makes it easy to write tests dramatically affect developer productivity. Along with readability, it not only makes it easier for you to maintain and improve old code. You also make it easier for new people to add to the project as they can feel somewhat safe when changing or adding code as you already have test in place. Lastly, having adaptability means your project can adapt to changing needs and complexity.

Starting out as simple as possible have the advantage of you being able to quickly iterate on new ideas getting to the actual problem solving right from the get-go. Most developers don’t need to think about domain-driven development when starting new projects, and it would probably be a waste of time in most cases (unless, of course, you are working in Big Co. but then why would you be reading an article for beginner Go developers).

The tutorial is split into three parts:

- The application structure
- User Service implementation
- TDD implementation of Weight Service

This should hopefully give you an overview of the overall structure, a step-by-step guide on how to add a service and finally how easy it is to write test and develop with TDD.

## The application structure

Alright, enough talk, let’s go going. Create a new folder wherever you like to store your projects called *weight-tracker* and run the following command inside the folder:


```go
go mod init weight-tracker
```

The application is going to follow a structure that looks something like this:


```go
weight-tracker
- cmd
  - server
    main.go
- pkg
  - api
    user.go
    weight.go
  - app
    server.go
    handlers.go
    routes.go
  - repository
    storage.go
```

Go ahead and create the above structure in your designated folder.

Most go projects seem to be following a convention of having a *cmd* and a *pkg* directory. The *cmd* will be the entry point into the program and gives you the flexibility to interact with the program in several ways. The *pkg* directory will contain everything else:
- routes
- database interactions
- services etc.

We will be working with four packages:
- main
- api
- app
- repository

The *main* package should be self-explanatory if you’ve ever done any golang programming before. All of our services, i.e. *user* and *weight* service, are going into the *api* package along with a *definitions* file that will contain all of our structs (we will create this later). In the *app* package, we will have our *server*, *handlers* and *routes*. Lastly, we have *repository* that will contain all of our code related to database operations.

Open up *main.go* and add the following:

```go
package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
	"weight-tracker/pkg/api"
	"weight-tracker/pkg/app"
	"weight-tracker/pkg/repository"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "this is the startup error: %s\\n", err)
		os.Exit(1)
	}
}

// func run will be responsible for setting up db connections, routers etc
func run() error {
	// I'm used to working with postgres, but feel free to use any db you like. You just have to change the driver
	// I'm not going to cover how to create a database here but create a database
	// and call it something along the lines of "weight tracker"
	connectionString := "postgres://postgres:postgres@localhost/**NAME-OF-YOUR-DATABASE-HERE**?sslmode=disable"

	// setup database connection
	db, err := setupDatabase(connectionString)

	if err != nil {
		return err
	}

	// create storage dependency
	storage := repository.NewStorage(db)

	// create router dependecy
	router := gin.Default()
	router.Use(cors.Default())

	// create user service
	userService := api.NewUserService(storage)

	// create weight service
	weightService := api.NewWeightService(storage)

	server := app.NewServer(router, userService, weightService)

	// start the server
	err = server.Run()

	if err != nil {
		return err
	}

	return nil
}

func setupDatabase(connString string) (*sql.DB, error) {
	// change "postgres" for whatever supported database you want to use
	db, err := sql.Open("postgres", connString)

	if err != nil {
		return nil, err
	}

	// ping the DB to ensure that it is connected
	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}
```

You should have quite a few errors showing in your IDE now, we are fixing that in a moment.
But first off, let me explain what is going on here. In our *func main()* we call a method named *run*, that returns an error (or nil, if there is no error). If *run* should return an error our programs exits and gives us an error message. Setting up our *main* function in this way allows us to test it and thereby following the element of a robust service structure, testability. This way of setting up the main function comes from Mat Ryer, who talks more about it in his [blog post](https://pace.dev/blog/2020/02/12/why-you-shouldnt-use-func-main-in-golang-by-mat-ryer.html).

The next two topics we need to discuss are **Clean Architecture** and **Dependency Injection**. You might have heard about Clean Architecture under a different name, as a few people wrote about it around the same time. The majority of my inspiration comes from [Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) so this is who I will be referencing.

We are not going to be religiously following Clean Architecture but mainly adopt the idea of The Dependency Rule. Take a look at the below picture:

![https://cdn-images-1.medium.com/max/1600/1*_5WeMzRt5aCVxXNWLlxAJw.png](https://cdn-images-1.medium.com/max/1600/1*_5WeMzRt5aCVxXNWLlxAJw.png)

Source: Robert C. Martin’s blog — Clean Coder Blog

A quote from Robert Martin’s article helps define this rule:

> This rule says that source code dependencies can only point inwards. Nothing in an inner circle can know anything at all about something in an outer circle.

In rough terms, the inner layers should not know about the outer layers. This allows us to practically change the database we use. Say we are changing from PostgreSQL to MySQL or exposing our API through gRPC instead of HTTP. You would most likely never do this but it does speak for the adaptability of this concept.

To achieve this, we are going to use dependency injection (DI for short). Basically, DI is the idea that services should receive their dependencies when they are created. It allows us to decouple the creation of a service’s dependencies from the creation of the service itself. This will be helpful when we get to testing our code. If you want to read more about DI, I suggest this [article](https://blog.drewolson.org/dependency-injection-in-go).

I learn best by looking at actual code, so let’s start by adding some of the missing code that will make the code in `main.go` make more sense. Open up `storage.go` and add the following:

```go
package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage interface{}

type storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) Storage {
	return &storage{
		db: db,
	}
}
```

This allows us to create a new storage component where and whenever we want, as long as it receives a valid db argument of type `*sql.DB`.

As you may have noticed, we have a lowercase and an uppercase version of `storage`, one is a struct and one is an interface. We will define any methods (like CRUD operations) on the lowercase version of storage and define the methods in the uppercase version. By doing so, we can now easily mock out `storage` for unit testing. Also, we now get some nice suggestions from our IDE when we add methods to the storage interface when we haven’t implemented the methods.

Now, lets set up the base structure of one of our services, `user.go`. This should give you a feel for how the API package is going to be structuring services. You will have to repeat this for the `weight.go` service as well. Just copy-paste the content of `user.go`and change the naming. Open up `user.go` and add the following:


```go
package api

// UserService contains the methods of the user service
type UserService interface{}

// UserRepository is what lets our service do db operations without knowing anything about the implementation
type UserRepository interface{}

type userService struct {
	storage UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{
		storage: userRepo,
	}
}
```

Notice that our user service have a repository dependency. We will define repository methods (i.e. database operations) here later on, but the important part is that we only define the methods we need. Since we are not interested in the implementation of these methods, only the behaviour, we can write mock functions that suit a given test case.

This might seem a bit fluffy right now but bear with me for now, it will make sense later on. Let’s get a running application that we can add actual logic to the services.

Open up `server.go` and add the following:


```go
package app

import (
	"github.com/gin-gonic/gin"
	"log"
	"weight-tracker/pkg/api"
)

type Server struct {
	router        *gin.Engine
	userService   api.UserService
	weightService api.WeightService
}

func NewServer(router *gin.Engine, userService api.UserService, weightService api.WeightService) *Server {
	return &Server{
		router:        router,
		userService:   userService,
		weightService: weightService,
	}
}

func (s *Server) Run() error {
	// run function that initializes the routes
	r := s.Routes()

	// run the server through the router
	err := r.Run()

	if err != nil {
		log.Printf("Server - there was an error calling Run on router: %v", err)
		return err
	}

	return nil
}
```

Next, open `routes.go` and add the following:


```go
package app

import "github.com/gin-gonic/gin"

func (s *Server) Routes() *gin.Engine {
	router := s.router

	// group all routes under /v1/api
	v1 := router.Group("/v1/api")
	{
		v1.GET("/status", s.ApiStatus())
	}

	return router
}
```

We are going to take advantage of gin’s group functionality so we can easily group our endpoints after the resources they intend to serve. Next up, let’s add a handler so we can actually make calls to the status endpoint and verify that our application is running. Open up `handlers.go`:


```go
package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) ApiStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		response := map[string]string{
			"status": "success",
			"data":   "weight tracker API running smoothly",
		}

		c.JSON(http.StatusOK, response)
	}
}
```

At this point, we only need to sync a few dependencies: Gin Web Framework and a driver for PostgreSQL. Go ahead and type the following into your terminal: 
```sh 
  go get github.com/gin-contrib/cors
  go get github.com/gin-gonic/gin
  go get github.com/lib/pg
```

Everything should be ready now, so go into your terminal and write: ```go go run
cmd/server/main.go```, visit ```sh http://localhost:8080/v1/api/status```, and you should receive a message along the lines of:

```json
{
	"data": "weight tracker API running smoothly",
	"status": "success"
}
```

If you don’t get the above message, please go back to the previous steps and see if you missed something or checkout the accompanying [GitHub repository](https://github.com/MBvisti/weight-tracker-article).

### User Service implementation

At this point, we are ready to start building the meat of the application. To do so, we probably need to set up a database with tables and relations. To do this, we are going to use the library `[golang-migrate](https://github.com/golang-migrate/migrate/)`. I personally like this library as it has a CLI tool to add migrations which enables me to do things like adding Makefile commands to create migrations. I encourage you to give the library’s documentation a look, it is an amazing project.

I won’t cover how to set this up as the article would be too long. For now, go into your terminal, make sure you are in the `repository` folder and run:

```sh
git clone https://gist.github.com/ed090d782dc6ebb35e344ff82aafdddf.git
```

This clones the migrations needed for the project. You should also now have a folder named `ed090d782dc6ebb35e344ff82aafdddf`, ****lets change that to migrations by running:


```sh
mv ed090d782dc6ebb35e344ff82aafdddf migrations
```

The last thing we need now is to add the `RunMigrations` method to `storage.go`:


```go
// update your imports to look like this:
import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"path/filepath"
	"runtime"
)

// add the RunMigrations method to our interface
type Storage interface {
	RunMigrations(connectionString string) error
}

// add this below NewStorage
func (s *storage) RunMigrations(connectionString string) error {
	if connectionString == "" {
		return errors.New("repository: the connString was empty")
	}
	// get base path
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../..")

	migrationsPath := filepath.Join("file://", basePath, "/pkg/repository/migrations/")

	m, err := migrate.New(migrationsPath, connectionString)

	if err != nil {
		return err
	}

	err = m.Up()

	switch err {
	case errors.New("no change"):
		return nil
	}

	return nil
}
```

To run the migrations, open up `main.go` and add the following:

```go
// everything stays the same, so add this below
// storage := repository.NewStorage(db)
// run migrations

// note that we are passing the connectionString again here. This is so
// we can easily run migrations against another database, say a test version,
// for our integration- and end-to-end tests
err = storage.RunMigrations(connectionString)

if err != nil {
    return err
}
```

You should now have two tables in your database: user and weight. Let’s get started on writing some actual business logic.

We want to let people create an account through our API. So let’s started by defining what a user request is, create a file called under the API folder `definitions.go` and add the following:


```go
package api

type NewUserRequest struct {
	Name          string `json:"name"`
	Age           int    `json:"age"`
	Height        int    `json:"height"`
	Sex           string `json:"sex"`
	ActivityLevel int    `json:"activity_level"`
	WeightGoal    string `json:"weight_goal"`
	Email         string `json:"email"`
}
```

We are defining what a new user request should look like. Note that this might be different from what a user struct would look like, we define our struct to only include the data we need. Next up, open `user.go` and the following to `UserService` and `UserRepository`:


```go
// user.go
type UserService interface {
    New(user NewUserRequest) error
}

// User repository is what lets our service do db operations without knowing anything about the implementation
type UserRepository interface {
    CreateUser(NewUserRequest) error
}
```

Here we define a method on our `UserService`, called `New` and a method on the UserRepository called `CreateUser`. Remember we talked about The Dependency Rule earlier? This is what’s going on with the `CreateUser` method on `UserRepository`, our service does not know about the actual implementation of the method, what type of database it is etc. Just that there is a method called `CreateUser` and takes in a `NewUserRequest` and returns an error. The benefit of this is twofold: we get some indication from our IDE that a method is missing (open up `main.go` and check `api.NewUserService`) and what it needs, and it allows us to easily write unit tests. You should also see an error from `NewUserService`in `user.go`, telling us that we are missing a method. Let’s fixed that, add the following:


```go
// add these imports after the package declaration
import (
	"errors"
	"strings"
)

// add this after NewUserService
func(u * userService) New(user NewUserRequest) error {
    // do some basic validations
    if user.Email == "" {
        return errors.New("user service - email required")
    }

    if user.Name == "" {
        return errors.New("user service - name required")
    }

    if user.WeightGoal == "" {
        return errors.New("user service - weight goal required")
    }

    // do some basic normalisation
    user.Name = strings.ToLower(user.Name)
    user.Email = strings.TrimSpace(user.Email)

    err := u.storage.CreateUser(user)

    if err != nil {
        return err
    }

    return nil
}
```

We do some basic validations and normalisation, but this method could definitely be improved upon. We still need to add the `CreateUser` method, so open up `storage.go` and add the following to the `CreateUser` method to the `Storage` interface:


```go
// storage.go
type Storage interface {
    RunMigrations(connectionString string) error
    CreateUser(request api.NewUserRequest) error
}
```

Notice that this makes the error in `main.go` go away, but results in a new error on the `NewStorage` function. We need to implement the method, just like we did with the `UserService`. Add this below `RunMigrations`:


```go
// add "log" to imports to your imports look like this:
import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"path/filepath"
	"runtime"
	"weight-tracker/pkg/api"
)

// add this below RunMigrations
func (s *storage) CreateUser(request api.NewUserRequest) error {
	newUserStatement := `
		INSERT INTO "user" (name, age, height, sex, activity_level, email, weight_goal)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
		`

	err := s.db.QueryRow(newUserStatement, request.Name, request.Age, request.Height, request.Sex, request.ActivityLevel, request.Email, request.WeightGoal).Err()

	if err != nil {
		log.Printf("this was the error: %v", err.Error())
		return err
	}

	return nil
}
```

Again, this code could probably do with some improvements, but I’m going to keep it short.

Now, all that is left to do is to expose this through HTTP so that people can actually start using our API and create their accounts. Open handlers.go and add the following:


```go
// update your imports to look like this
import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"weight-tracker/pkg/api"
)

// add this below APIStatus method
func (s *Server) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var newUser api.NewUserRequest

		err := c.ShouldBindJSON(&newUser)

		if err != nil {
			log.Printf("handler error: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		err = s.userService.New(newUser)

		if err != nil {
			log.Printf("service error: %v", err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		response := map[string]string{
			"status": "success",
			"data":   "new user created",
		}

		c.JSON(http.StatusOK, response)
	}
}
```

We are going to accept a request with a JSON payload and using gin’s ShouldBindJSON method to extract the data. Go ahead and try it out!

It took some time to get here, so let me show you one of the benefits that I keep talking about: **testability**.

Create a file called `user_test.go` under the API folder and add the following:


```go
package api_test

import (
	"errors"
	"reflect"
	"testing"
	"weight-tracker/pkg/api"
)

type mockUserRepo struct{}

func (m mockUserRepo) CreateUser(request api.NewUserRequest) error {
	if request.Name == "test user already created" {
		return errors.New("repository - user already exists in database")
	}

	return nil
}

func TestCreateNewUser(t *testing.T) {
	mockRepo := mockUserRepo{}
	mockUserService := api.NewUserService(&mockRepo)

	tests := []struct {
		name    string
		request api.NewUserRequest
		want    error
	}{
		{
			name: "should create a new user successfully",
			request: api.NewUserRequest{
				Name:          "test user",
				WeightGoal:    "maintain",
				Age:           20,
				Height:        180,
				Sex:           "female",
				ActivityLevel: 5,
				Email:         "test_user@gmail.com",
			},
			want: nil,
		}, {
			name: "should return an error because of missing email",
			request: api.NewUserRequest{
				Name:          "test user",
				Age:           20,
				WeightGoal:    "maintain",
				Height:        180,
				Sex:           "female",
				ActivityLevel: 5,
				Email:         "",
			},
			want: errors.New("user service - email required"),
		}, {
			name: "should return an error because of missing name",
			request: api.NewUserRequest{
				Name:          "",
				Age:           20,
				WeightGoal:    "maintain",
				Height:        180,
				Sex:           "female",
				ActivityLevel: 5,
				Email:         "test_user@gmail.com",
			},
			want: errors.New("user service - name required"),
		}, {
			name: "should return error from database because user already exists",
			request: api.NewUserRequest{
				Name:          "test user already created",
				Age:           20,
				Height:        180,
				WeightGoal:    "maintain",
				Sex:           "female",
				ActivityLevel: 5,
				Email:         "test_user@gmail.com",
			},
			want: errors.New("repository - user already exists in database"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := mockUserService.New(test.request)

			if !reflect.DeepEqual(err, test.want) {
				t.Errorf("test: %v failed. got: %v, wanted: %v", test.name, err, test.want)
			}
		})
	}
}
```

A lot of things are happening here so let us go through them step-by-step.

We start by creating a struct called `mockUserRepo`. Our `UserService`knows that it needs a `UserRepository` with the following method: `CreateUser`. However, it does care about the actual implementation of said method, which allows us to mock the behaviour any way we want. In this case, we say that when the Name of the request is equal to “test user already created” (a bad name I know, but hopefully you get the point), return an error, and if not, just return nil. This allows us to mock the behaviour of our database to make it fit different situations, and test if our logic handles it in a way we expect it too.

Next, we create a new variable called `mockRepo` of type `mockUserRepo` We then create a `mockUserService` and pass the `mockRepo` to it, and we are ready to go! The actual test is what is known as a table-driven test. I won’t go into details about it as it is beyond the scope of this article, but if you want to know more, check out Dave Chaney’s article about it [here](https://dave.cheney.net/2013/06/09/writing-table-driven-tests-in-go).

Now, whenever we add a method to the `UserRepository` we would also have to add it here, in our `mockUserRepo`. We would of course also like to have some integration test, but what I really wanted to show with this article was how to do unit tests as these are cheap and easy to write.

Run `go test ./...` from the root directory and all tests should pass.

### TDD implementation of Weight Service

Our application is not worth much right now. We can only create a user but not track our weight or calculate the amount of calories we need. Let’s change that!

Our application is not worth much right now. We can only create a user but not track our weight or calculate the number of calories we need. Let’s change that!I’m a big fan of test-driven development (TDD) and the way this application is structured makes it really easy to use. We are going to need three methods on our weight service: New, `CalculateBMR` and `DailyIntake` and two on our repository: `CreateWeightEntry` and GetUser. Open `weight.go` and add the following to `WeightService` and `WeightRepository`:


```go
// weight.go

// WeightService contains the methods of the user service
type WeightService interface {
	New(request NewWeightRequest) error
	CalculateBMR(height, age, weight int, sex string) (int, error)
	DailyIntake(BMR, activityLevel int, weightGoal string) (int, error)
}

type WeightRepository interface {
	CreateWeightEntry(w Weight) error
	GetUser(userID int) (User, error)
}
```

Next, we need to add three structs to our definitions file: `´NewWeightRequest`, `Weight` and `User`. Open up `weight.go` and add the following:


```go
// add this after the package declaration
import "time"

// add this below NewUserRequest
type User struct {
	ID            int       `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Name          string    `json:"name"`
	Age           int       `json:"age"`
	Height        int       `json:"height"`
	Sex           string    `json:"sex"`
	ActivityLevel int       `json:"activity_level"`
	WeightGoal    string    `json:"weight_goal"`
	Email         string    `json:"email"`
}

type Weight struct {
	Weight             int `json:"weight"`
	UserID             int `json:"user_id"`
	BMR                int `json:"bmr"`
	DailyCaloricIntake int `json:"daily_caloric_intake"`
}

type NewWeightRequest struct {
	Weight int `json:"weight"`
	UserID int `json:"user_id"`
}
```

Now, you will see an error on `NewWeightService` about missing methods. We don’t want to write the actual implementation yet, as we are doing TDD, so for now, just add this below `NewWeightService`:


```go
// add these below NewWeightService
func (w *weightService) CalculateBMR(height, age, weight int, sex string) (int, error) {
	panic("implement me")
}

func (w *weightService) DailyIntake(BMR, activityLevel int, weightGoal string) (int, error) {
	panic("implement me")
}

func (w *weightService) New(request NewWeightRequest) error {
	panic("implement me")
}
```

Open up `main.go` and you will see that we also have some missing methods on storage passed to `api.NewWeightService`. Open up `storage.go` and add these:


```go
// add this to the Storage interface
type Storage interface {
	CreateWeightEntry(request api.Weight) error
	GetUser(userID int) (api.User, error)
}

// add this below the CreateUser method
func (s *storage) CreateWeightEntry(request api.Weight) error {
	panic("implement me!")
}

func (s *storage) GetUser(userID int) (api.User, error) {
	panic("implement me!")
}
```

Lets now add the tests that we will need, create a file called`weight_test.go` in the API folder and add the following:


```go
package api_test

import (
	"errors"
	"reflect"
	"testing"
	"weight-tracker/pkg/api"
)

type mockWeightRepo struct{}

func (m mockWeightRepo) CreateWeightEntry(w api.Weight) error {
	return nil
}

func (m mockWeightRepo) GetUser(userID int) (api.User, error) {
	if userID != 1 {
		return api.User{}, errors.New("storage - user doesn't exists")
	}

	return api.User{
		ID:            userID,
		Name:          "Test user",
		Age:           20,
		Height:        185,
		WeightGoal:    "maintain",
		Sex:           "female",
		ActivityLevel: 5,
		Email:         "test@mail.com",
	}, nil
}

func TestCreateWeightEntry(t *testing.T) {
	mockRepo := mockWeightRepo{}
	mockUserService := api.NewWeightService(&mockRepo)

	tests := []struct {
		name    string
		request api.NewWeightRequest
		want    error
	}{
		{
			name: "should create a new user successfully",
			request: api.NewWeightRequest{
				Weight: 70,
				UserID: 1,
			},
			want: nil,
		}, {
			name: "should return a error because user already exists",
			request: api.NewWeightRequest{
				Weight: 70,
				UserID: 2,
			},
			want: errors.New("storage - user doesn't exists"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := mockUserService.New(test.request)
			if !reflect.DeepEqual(err, test.want) {
				t.Errorf("test: %v failed. got: %v, wanted: %v", test.name, err, test.want)
			}
		})
	}
}

func TestCalculateBMR(t *testing.T) {
	mockRepo := mockWeightRepo{}
	mockUserService := api.NewWeightService(&mockRepo)

	tests := []struct {
		name   string
		Height int
		Age    int
		Weight int
		Sex    string
		want   int
		err    error
	}{
		{
			name:   "should calculate BMR for a female",
			Height: 170,
			Age:    22,
			Weight: 65,
			Sex:    "female",
			want:   1441,
			err:    nil,
		}, {
			name:   "should calculate BMR for a male",
			Height: 170,
			Age:    22,
			Weight: 65,
			Sex:    "male",
			want:   1607,
			err:    nil,
		}, {
			name:   "should return error because sex wasn't properly specified",
			Height: 170,
			Age:    22,
			Weight: 65,
			Sex:    "",
			want:   0,
			err:    errors.New("invalid variable sex provided to CalculateBMR. needs to be either male or female"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			BMR, err := mockUserService.CalculateBMR(test.Height, test.Age, test.Weight, test.Sex)
			if !reflect.DeepEqual(err, test.err) {
				t.Errorf("test: %v failed. got: %v, wanted: %v", test.name, err, test.want)
			}

			if !reflect.DeepEqual(BMR, test.want) {
				t.Errorf("test: %v failed. got: %v, wanted: %v", test.name, BMR, test.want)
			}
		})
	}
}

func TestDailyIntake(t *testing.T) {
	mockRepo := mockWeightRepo{}
	mockUserService := api.NewWeightService(&mockRepo)

	tests := []struct {
		name          string
		BMR           int
		ActivityLevel int
		weightGoal    string
		want          int
		err           error
	}{
		{
			name:          "should calculate daily intake for activity level 1 with weight loss as goal",
			weightGoal:    "loose",
			BMR:           1441,
			ActivityLevel: 1,
			want:          1229,
			err:           nil,
		}, {
			name:          "should calculate daily intake for activity level 2 with weight loss as goal",
			weightGoal:    "loose",
			BMR:           1441,
			ActivityLevel: 2,
			want:          1481,
			err:           nil,
		}, {
			name:          "should calculate daily intake for activity level 3 with weight loss as goal",
			weightGoal:    "loose",
			BMR:           1441,
			ActivityLevel: 3,
			want:          1733,
			err:           nil,
		}, {
			name:          "should calculate daily intake for activity level 4 with weight increase as goal",
			weightGoal:    "gain",
			BMR:           1441,
			ActivityLevel: 4,
			want:          2985,
			err:           nil,
		}, {
			name:          "should calculate daily intake for activity level 5 with weight maintenance as goal",
			weightGoal:    "maintain",
			BMR:           1441,
			ActivityLevel: 5,
			want:          2737,
			err:           nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			BMR, err := mockUserService.DailyIntake(test.BMR, test.ActivityLevel, test.weightGoal)
			if !reflect.DeepEqual(err, test.err) {
				t.Errorf("test: %v failed. got: %v, wanted: %v", test.name, err, test.want)
			}

			if !reflect.DeepEqual(BMR, test.want) {
				t.Errorf("test: %v failed. got: %v, wanted: %v", test.name, BMR, test.want)
			}
		})
	}
}
```

The actual content of the tests is not so important as the fact that we can quickly write tests and mock external dependencies such as the methods to interact with the database.

Run `go test ./...` from the root directory and you should receive a lot of failed tests. We are going to fix this by implementing the logic of these methods, starting with the three in our `WeightService`. Add this below `NewWeightService`:


```go
// add this after the package declaration
import "errors"

// we define some activity multipliers here - these are just part of the
// formula I choose to calculate your daily caloric intake.
const (
	// Very Low Intensity activity multiplier - 1
	veryLowActivity = 1.2
	// Light exercise activity multiplier - 3-4x per week - 2
	lightActivity = 1.375
	// Moderate exercise activity multiplier - 3-5x per week 30-60 mins/session - 3
	moderateActivity = 1.55
	// High exercise activity multiplier - (6-7x per week for 45-60 mins) - 4
	highActivity = 1.725
	// Very high exercise activity multiplier - for people who is an Athlete - 5
	veryHighActivity = 1.9
)

func (w *weightService) New(request NewWeightRequest) error {
	if request.UserID == 0 {
		return errors.New("weight service - user ID cannot be 0")
	}

	user, err := w.storage.GetUser(request.UserID)

	if err != nil {
		return err
	}

	bmr, err := w.CalculateBMR(user.Height, user.Age, request.Weight, user.Sex)

	if err != nil {
		return err
	}

	dailyIntake, err := w.DailyIntake(bmr, user.ActivityLevel, user.WeightGoal)

	if err != nil {
		return err
	}

	newWeight := Weight{
		Weight:             request.Weight,
		UserID:             user.ID,
		BMR:                bmr,
		DailyCaloricIntake: dailyIntake,
	}

	err = w.storage.CreateWeightEntry(newWeight)

	if err != nil {
		return err
	}

	return nil
}

func (w *weightService) CalculateBMR(height, age, weight int, sex string) (int, error) {
	var sexModifier int

	switch sex {
	case "male":
		sexModifier = -5
	case "female":
		sexModifier = 161
	default:
		return 0, errors.New("invalid variable sex provided to CalculateBMR. needs to be either male or female")
	}

	return (10 * weight) + int(float64(height)*6.25) - (5 * age) - sexModifier, nil
}

func (w *weightService) DailyIntake(BMR, activityLevel int, weightGoal string) (int, error) {
	var maintenanceCalories int

	switch activityLevel {
	case 1:
		maintenanceCalories = int(float64(BMR) * veryLowActivity)
	case 2:
		maintenanceCalories = int(float64(BMR) * lightActivity)
	case 3:
		maintenanceCalories = int(float64(BMR) * moderateActivity)
	case 4:
		maintenanceCalories = int(float64(BMR) * highActivity)
	case 5:
		maintenanceCalories = int(float64(BMR) * veryHighActivity)
	default:
		return 0, errors.New("invalid variable activityLevel - needs to be 1, 2, 3, 4 or 5")
	}

	var dailyCaloricIntake int

	switch weightGoal {
	case "gain":
		dailyCaloricIntake = maintenanceCalories + 500
	case "loose":
		dailyCaloricIntake = maintenanceCalories - 500
	case "maintain":
		dailyCaloricIntake = maintenanceCalories
	default:
		return 0, errors.New("invalid weight goal provided - must be gain, loose or maintain")
	}

	return dailyCaloricIntake, nil
}
```

With these methods added all of the tests should pass.

The last thing we need is to add the methods needed to interact with the database. Open `storage.go` and add the following:


```go
// replace the methods: CreateWeightEntry and GetUser with the code below
func (s *storage) CreateWeightEntry(request api.Weight) error {
	newWeightStatement := `
		INSERT INTO weight (weight, user_id, bmr, daily_caloric_intake)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
		`

	var ID int
	err := s.db.QueryRow(newWeightStatement, request.Weight, request.UserID, request.BMR, request.DailyCaloricIntake).Scan(&ID)

	if err != nil {
		log.Printf("this was the error: %v", err.Error())
		return err
	}

	return nil
}

func (s *storage) GetUser(userID int) (api.User, error) {
	getUserStatement := `
		SELECT id, name, age, height, sex, activity_level, email, weight_goal FROM "user"
		where id=$1;
		`

	var user api.User
	err := s.db.QueryRow(getUserStatement, userID).Scan(&user.ID, &user.Name, &user.Age, &user.Height, &user.Sex, &user.ActivityLevel, &user.Email, &user.WeightGoal)

	if err != nil {
		log.Printf("this was the error: %v", err.Error())
		return api.User{}, err
	}

	return user, nil
}
```

We didn’t do this in true TDD style, however, it should again paint a picture of how we can structure Go applications and develop them using TDD. Feel free to re-implement this section and do it in true TDD style: create one test, make it pass, create the next one, and so on.

We have almost everything we need now. The last thing we need is to add routes and a handler to create a weight entry for a given user. I’m going to leave this up to the reader to implement, as the building blocks should be in place. If you are feeling lazy you can just check out the accompanying Github repository.

## Conclusion

Hopefully, this tutorial gave you some insights on how to structure your Golang applications. It was a long one I know, but hopefully, you didn’t waste your time. In a part two, I will show how to add Makefile commands, integration tests and easily deploy the application. If you have any questions or criticisms, please do not hesitate to reach out.

My [twitter](https://twitter.com/mbvisti) and [github](https://github.com/mbvisti).

### Resources

- [Project GitHub](https://github.com/MBvisti/weight-tracker-article)
- [The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Golang-migrate](https://github.com/golang-migrate/migrate/)
- [Table-driven tests](https://dave.cheney.net/2013/06/09/writing-table-driven-tests-in-go)
