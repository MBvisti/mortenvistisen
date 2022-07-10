<img src="https://res.cloudinary.com/practicaldev/image/fetch/s--P1NWBtsb--/c_imagga_scale,f_auto,fl_progressive,h_900,q_auto,w_1600/https://thepracticaldev.s3.amazonaws.com/i/rlyibpr58qk49ci8y1rk.png" alt="gopher" />

Picture this: you've just left Node for the promised land of Go. You've learned
about composition, interfaces, simplicity, domain-driven design (and understood
nothing) and unit tests. You feel all grown-up and ready to take on concurrency
and ride off into the promised land. You've built a sweet CRUD application that
will revolutionise how we do X. Since this is a sure-fire success; you've of
course already contacted a real estate agent to view that penthouse apartment
you've always wanted and tweeted at multiple VC funds and investors. You're
basically ready to go live but before you release your code into the wild, you
whip out your favourite terminal and type in go run main.go. It runs, glory
days!

You set up multiple demos and investor meetings. You add "CEO and Founder" to
your Linkedin profile (and tinder, of course) and tell your mom that you're
going to vacate the basement any day now. You show up to your first meeting, do
a couple of power poses and start demoing that
really-crucial-endpoint-that-is-the-backbone-of-the-app. With your right thumb
on the ENTER key (and your left thumb mentally ready to pop the champagne), you
go ahead. Aaaaaand, it completely crashes. What the fuck. Not only do you get
laughed at your super important meeting, but you also have to tell the real
estate agent that you probably can't afford that penthouse apartment, you have
also to stay in your mom's basement (and continue trying to convince your tinder
dates that it's definitely a temporary solution).

---

Hopefully, you haven't experienced the above scenario and are just here because
you typed some words into google. I'm going to show you how I set up and run
integration tests in Go, using docker and some nifty Makefile cmds.

## Aim of the article
A lot has been written on this subject (check out the resources list at the end
of the article). This is not an attempt to improve on what is currently
available. It's an attempt to show how you can set up and configure integrations
tests, both locally and externally (CI env), that is extendable and portable.
I'm sure there are ways to improve upon what I'm going to show you here,
however, the main goal is to provide you with the tools to build out a great
test suite. Once you've understood the concepts and configurations it's much
easier to customise them to your need. Furthermore, you get a better idea of
what to google (unknown unknowns are a tricky subject), so you can solve your
own issues that might come up later on.

### Requirements
Alright, first thing first. You'll need some tools, so please install and setup the 
following:
- Go 
- Docker

# There are more than one way to skin a cat
As with everything in software, there are a lot of ways of doing the same thing
and even more opinions about how to do that thing. With integration testing, we
might discuss how to set up/configure the test suite, what library (if any) to
use and even when something classifies as an integration vs. unit test, E2E
test, and so forth. This can make concepts more complex than they need to be
especially when you're just starting out. I think one of the best ways to
approach such situations is to agree on some sort of guiding principles/first
principles, that becomes the basis for what we do. I've found one of the best
principles for this is:
> test behaviour, not implementation details

To expand upon this, we want to have our tests cover cases/scenarios of how a
user might actually use our software. I use the terms user and software very
broadly here as it's very context-specific. A user might be a co-worker who
calls some method in the codebase that we have created or it might also be
someone consuming our API. Our tests shouldn't really worry about how we
actually implemented the code, only that given X input we see Y response. We are
in a sense testing that the "contract" we have made as developers by exposing
our code/functions/methods to the world, is actually kept.

If we want a more concrete definition of an integration test we can dust off the
big ol' google book and see how they define it:
>...medium-scoped tests (commonly called integration tests) are designed to verify 
interactions between a small number of components; for example, between a server and a 
database

<div className="w-full flex justify-center"><p className="italic">Source: Software 
Engineering at Google, chap. 11</p></div>

If we borrow some terminology from [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) 
we can think of it as when we include code from the infrastructure layer in our
tests, we're in the integration testing territory.

## What are we doing
When I started out writing integration tests in Go; what I struggled the most
with was how to configure and set it up for real-world usage. This, again, will
also differ based upon the developer/use-case/reality as this approach might not
work for well a big multinational company. My criteria for my integration tests
are that they should run across different operating systems, be easy to run and
play nicely with CI/CD workflows (basically, play nice within a dockerized
environment).

I'm a strong believer in YAGNI (you're not going to need it), so I'm going to
show you two ways of setting up integration tests:
- the vanilla only-using-the-standard-library approach
- using a library

This should hopefully illustrate how you can start out relatively simple (we
could skip the Docker part, but honestly, that would make it a little trickier
to set up our CI/CD flow) and then add on as needed.

## What we are testing
I'm going to re-use some of the code from my article on how to structure Go apps
(which can be found [here](https://blog.mortenvistisen.com/posts/a-practical-approach-to-structuring-golang-applications). If you haven't read it, it basically builds a small
app that lets your track your weight gain during the lockdown. The article is
due for an update (I would suggest checking this [repo](https://github.com/bnkamalesh/goapp) for a great example on how
to structure Go apps using DDD and clean architecture), so we're just going to
focus on adding integration tests based on the code (a good example of testing
behaviour vs. implementation details). We want to make sure that calls to our
services behave as expected.

You can find the repo [here](https://github.com/MBvisti/integration-test-in-go).

# "Infrastructure" Setup
Much of modern web development uses Docker and this tutorial will be no
exception. This is not a tutorial on Docker so I won't be going into much detail
about the setup, but provide some foundations on how to get started. There are
ways to improve upon this setup by extending the docker-compose file, using tags
inside of our Dockerfile etc. But I often find that to give me more headaches
than gains. We're violating DRY, to some extent, but it does give us the ability
to have completely separate dev and test environment. And after reading this,
you could yourself try to shorten the docker setup.

### Docker test environment setup
```docker
FROM golang:1-alpine

RUN apk add --no-cache git gcc musl-dev

RUN go get -u github.com/rakyll/gotest

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app
```

And the ```docker-compose.yaml``` file:
```yaml
version: "3.8"
services:
    database:
        image: postgres:13
        container_name: test_weight_tracker_db_psql
        environment:
            - POSTGRES_PASSWORD=password
            - POSTGRES_USER=admin
            - POSTGRES_DB=test_weight_tracker_database
        ports:
            - "5436:5432" # we are mapping the port 5436 on the local machine
              # to the image running inside the container
        volumes:
            - test-pgdata:/var/lib/postgresql/data

    app:
        container_name: test_weight_tracker_app
        build:
            context: .
            dockerfile: Dockerfile.test
        ports:
            - "8080:8080"
        working_dir: /app
        volumes:
            - ./:/app
            - test-go-modules:/go
        environment:
            - DB_NAME=test_weight_tracker_database
            - DB_PASSWORD=${DB_PASSWORD}
            - DB_HOST=test_weight_tracker_db_psql
            - DB_PORT=${DB_PORT}
            - DB_USERNAME=${DB_USERNAME}
            - ENVIRONMENT=test
        depends_on:
            - database

volumes:
    test-pgdata:
    test-go-modules:
```

With that in place we can now easily run our integration tests, which can be done with:
```bash
make run-integration-tests
```

# Approach 1: Vanilla setup 
*Side note: if you want to see the code separated from the rest, check out the
branch `vanilla-approach/running-integration-tests-using-std-library`*

There is a tendency in the Go community to lean more towards the standard
library than to pull in external libraries. And for a good reason, you can do a
lot with just the standard library. I'm by no means a purist in this regard and
we are also using external libraries for the routing, migration etc. here, but I
think it gives a great understanding of what is happening by starting with the
standard library and then using other libraries as you go.
 
In the earlier segment, we got our infrastructure up and running so we have an
active database, a running application, or in this case, the possibility to
trigger a test run against our app.

To maximize our confidence in our code we want our integration tests to mimic
our production environment as much as possible. This means that we need to have
some setup and teardown functions that can run migrations, populate the database
with seed data and tear everything down after the test so we have a clean
environment each time. The last part is important as we want to have a reliable
test environment so we don't want previously run tests to affect the current one
running.

It's worth noting here that this also means that our integrations tests have to
run sequentially, and not in parallel, as we most likely would have our unit
test do which means it would take longer running the test suite. I wouldn't
worry too much about this at the beginning and simply focus on having some
integration tests with decent coverage. Once you start to reach unbearable long
integration test runs then it's time to start looking into other
setups/improving the current one. For example, we could create a new database
for each test and use that, but it would add a more complex setup. Or we could
use an SQLite database for our integration tests, but again, cross that
bridge once you get to it.

It would be a great exercise to try and change this code to try out different
strategies to speed up the test runs.

## Migrations
I'm a big fan of the golang-migrate library so that is what we are going to use
to write our migrations. In short, it generates up/down migration pairs and
treats each migration as a new version, so you can roll back to the last working
version if you need to. 

I'm not going to touch upon migration strategies here, so our tests will be
written with the assumption that we have a database with all the latest
migrations. To achieve this in our tests we will run the up version before each
test is run. For that we need the following function:

```golang
func RunUpMigrations(cfg config.Config) error {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../migrations")
	migrationDir := filepath.Join("file://" + basePath)
	db, err := sql.Open("postgres", cfg.GetDatabaseConnString())
	if err != nil {
		return errors.WithStack(err)
	}
	defer db.Close()
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.WithStack(err)
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(migrationDir, "postgres", driver)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, ErrNoNewMigrations) {
			return errors.WithStack(err)
		}
	}
	m.Close()
	return nil
}
```

After the test is down, we would like our environment to be clean so we also
need this function:

```golang
func RunDownMigrations(cfg config.Config) error {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../migrations")
	migrationDir := filepath.Join("file://" + basePath)
	db, err := sql.Open("postgres", cfg.GetDatabaseConnString())
	if err != nil {
		return errors.WithStack(err)
	}
	defer db.Close()
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.WithStack(err)
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(migrationDir, "postgres", driver)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := m.Down(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
```

Basically, we get a new connection to the database, create a new migrate
instance where we pass a path to the migrations folder, the database we use and
the driver. We then run the migrations and close the connection again, pretty
straightforward. 

Next up, we need some data in our database to run our tests against. I'm a
pretty big fan of just keeping it in a SQL file and then having a helper
function run the SQL script against the database. To do that we just need a
function similar to the two above:

```golang
func LoadFixtures(cfg config.Config) error {
	pathToFile := "/app/fixtures.sql"
	q, err := os.ReadFile(pathToFile)
	if err != nil {
		return errors.WithStack(err)
	}

	db, err := sql.Open("postgres", cfg.GetDatabaseConnString())
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = db.Exec(string(q))
	if err != nil {
		return errors.WithStack(err)
	}
	err = db.Close()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
```

With that setup we are ready to write our first tests. 

## Our first test
Technically, we could test the code interacting with the database by providing a
mocked version of the methods in the database/SQL package. But that doesn't
really give us much as it would be tricky to mock a situation where you, for
example, miss a variable in your .Scan method or have some syntax issue.
Therefore, I tend to write integration tests for all my database functionality.
Let's add a test for the CreateUser function. We need the following:

```golang

// testing the happy path only - to improve upon these tests, we could consider
// using a table test
func TestIntegration_CreateUser(t *testing.T) {
	// create a NewStorage instance and run migrations
	cfg := config.NewConfig()
	storage := psql.NewStorage()

	err := psql.RunUpMigrations(*cfg)
	if err != nil {
		t.Errorf("test setup failed for: CreateUser, with err: %v", err)
		return
	}

	// run the test
	t.Run("should create a new user", func(t *testing.T) {
		newUser, err := entity.NewUser(
			"Jon Snow", "male", "90", "thewhitewolf@stark.com", 16, 182, 1)
		if err != nil {
			t.Errorf("failed to run CreateUser with error: %v", err)
			return
		}

		// to ensure consistency we could consider adding in a static date
		// i.e. time.Date(insert-fixed-date-here)
		// creationTime := time.Now()
		err = storage.CreateUser(*newUser)
		// assert there is no err
		if err != nil {
			t.Errorf("failed to create new user with err: %v", err)
			return
		}

		// now lets verify that the user is actually created using a
		// separate connection to the DB and pure sql
		db, err := sql.Open("postgres", cfg.GetDatabaseConnString())
		if err != nil {
			t.Errorf("failed to connect to database with err: %v", err)
			return
		}
		queryResult := entity.User{}
		err = db.QueryRow("SELECT id, name, email FROM users WHERE email=$1",
			"thewhitewolf@stark.com").Scan(
			&queryResult.ID, &queryResult.Name, &queryResult.Email,
		)
		if err != nil {
			t.Errorf("this was query err: %v", err)
			return
		}

		if queryResult.Name != newUser.Name {
			t.Error(`failed 'should create a new user' wanted name did not match 
				returned value`)
			return
		}
		if queryResult.Email != newUser.Email {
			t.Error(`failed 'should create a new user' wanted email did not match 
				returned value`)
			return
		}
		if int64(queryResult.ID) != int64(1) {
			t.Error(`failed 'should create a new user' wanted id did not match 
				returned value`)
			return
		}

	})

	// // run some clean up, i.e. clean the database so we have a clean env
	// // when we run the next test
	t.Cleanup(func() {
		err := psql.RunDownMigrations(*cfg)
		if err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				return
			}
			t.Errorf("test cleanup failed for: CreateUser, with err: %v", err)
		}
	})
}
```

We start by creating a new instance of config and storage (just like we would in
main.go when running the entire application) and then run the up migrations
function. If nothing goes wrong, we should have something similar to what we
would have in production.

We then use the storage instance that we just set up to create a new user, open
a new connection to query for the user we just created and verify that said user
is created with the values we expect. After, use the Cleanup function provided
by the testing package to call the down migrations. This basically clears the
database.

One more thing you might notice is that we have a psql_test.go file. Open it,
and you will find the following function:

```golang

// TestMain gets run before running any other _test.go files in each package
// here, we use it to make sure we start from a clean slate
func TestMain(m *testing.M) {
	cfg := config.NewConfig()
	// make sure we start from a clean slate
	err := psql.DropEverythingInDatabase(*cfg)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
```

TestMain is a special function that gets called before all other tests in the
package it's located in. Here, we're being (justifiable, I would say) paranoid
and call a function that drops everything in the database so we are sure we are
starting from a clean slate. You can find the function in repository/psql.go if
you want to take a closer look.

And that's basically it for running integration tests against our database
functions. We could use table-driven tests here and probably should, but this
will do for illustration purposes. See here for an explanation of table-driven
tests if you don't know them. Next up, let's do a "proper" integration tests and
ensure that our endpoints are working as expected!

## Testing our endpoints
Now we get into the meat of the stuff. We've arrived at (or closer to, at least)
the definition found in the ol' google book. We're testing that multiple parts
of our codes work together, mostly in the infrastructure layer, works together
in the way we expect. Here, we want to ensure that whenever our API receives a
request that is fulfilling the contract we as developers put out, it does what
we want. That is, we want to test the happy path. Ideally, we would also want to
test the sad path (not sure if that's the word for it, but this is my article,
so now it is) but integration tests are more "expensive" so it's a delicate
balance. You could choose to mock out database responses and test the sad path
in a more unit-test kind of way, or you could add integration tests until the
time it takes to run the test suite becomes unbearable. I would probably err on
adding 1 integration test to many, and deal with the "cost" when it becomes too
big.

Alright, enough rambling. Let's get started. A side note here; I'm using gofiber
which is inspired by express, the Node web framework. The way I'm setting up the
POST request sort of depends on how gofiber does things. I say sort of because
the underlying thing when sending a post request from Go is using Marshaling. I
will point it out when we get to it, but just be aware that if you like gorilla
or gin you might have to google a bit.

### A quick rundown of router setup
We're not going to be spending much time on this, as you can find the code in
the repo. Basically, we have this:

```golang
type serverConfig interface {
	GetServerReadTimeOut() time.Duration
	GetServerWriteTimeOut() time.Duration
	GetServerPort() int64
}

type Http struct {
	router        *fiber.App
	serverPort    int64
	userHandler   userHandler
	weightHandler weightHandler
}

func NewHttp(
	cfg serverConfig, userHandler userHandler, weightHandler weightHandler) *Http {
	r := fiber.New(fiber.Config{
		ReadTimeout:  cfg.GetServerReadTimeOut(),
		WriteTimeout: cfg.GetServerWriteTimeOut(),
		AppName:      "Weight Tracking App",
	})
	return &Http{
		router:        r,
		serverPort:    cfg.GetServerPort(),
		userHandler:   userHandler,
		weightHandler: weightHandler,
	}
}
```

We set up an HTTP struct that has some dependencies to get our server up and
running with a router. On that struct, we define some server-specific methods.
It's pretty straightforward.

## Testing our endpoint to create a new user
Our endpoint is pretty simple. There are no middleware and authentication,
everybody can just spam our server with requests and create a ton of users.
That's not ideal, but also not really what we care about right now. We just want
to make sure our API does what it's supposed to do.

```golang

func TestIntegration_UserHandler_New(t *testing.T) {
	cfg := config.NewConfig()
	storage := psql.NewStorage()

	err := psql.RunUpMigrations(*cfg)
	if err != nil {
		t.Errorf("test setup failed for: CreateUser, with err: %v", err)
		return
	}

	userService := service.NewUser(storage)
	weightService := service.NewWeight(storage)

	userHandler := http.NewUserHandler(userService)
	weightHandler := http.NewWeightHandler(weightService)

	srv := http.NewHttp(cfg, *userHandler, *weightHandler)

	srv.SetupRoutes()
	r := srv.GetRouter()

	req := http.NewUserRequest{
		Name:          "Test user",
		Sex:           "male",
		WeightGoal:    "80",
		Email:         "test@gmail.com",
		Age:           99,
		Height:        185,
		ActivityLevel: 1,
	}
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(req)
	if err != nil {
		log.Fatal(err)
	}
	rq, err := h.NewRequest(h.MethodPost, "/api/user", &buf)
	if err != nil {
		t.Error(err)
	}
	rq.Header.Add("Content-Type", "application/json")

	res, err := r.Test(rq, -1)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Error(errors.New("create user endpoint did not return 200"))
	}

	// query the database to verify that a user was created based on the request
	// we sent
	newUser, err := storage.GetUserFromEmail(req.Email)
	if err != nil {
		t.Error(err)
	}

	if newUser.Height != req.Height {
		t.Error(errors.New("create user endpoint did not create user with correct details"))
	}

	t.Cleanup(func() {
		err := psql.RunDownMigrations(*cfg)
		if err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				return
			}
			t.Errorf("test cleanup failed for: CreateUser endpoint, with err: %v", err)
		}
	})
}
```

Most of this looks similar to what we had in the repository tests. We set up the
database, the services and lastly, the server. We create a request, encode it,
send it to our endpoint and check the response. An important thing to note here
is that we don't really know what is happening under the hood of this beast. We
just know that we sent a request with some data and that returns OK and create a
user in the database with the expected data. This is also known as black-box
testing. We don't care how this is done, we care that the expected behaviour
occurs.

One thing about the above code is that there is quite some repetitiveness in how
we set up the tests and tear them down after each run. It would be nice if we
didn't have to copy-paste all of this and take a long hot bath after each test
run because we violated DRY. We could do this ourselves of course, or we could
use *Approach 2 - using test suites with Testify*.

# Approach 2 - using Testify to run our integration tests
For this, we are going to use the testify package which I have used for quite
some time now. The main thing this does for us is save some lines on
configuration and ensuring consistency in our tests suites. It's easy enough to
have the entire codebase in your head when it's only this size, but as it grows,
having the setup and configuration done in one place makes things so much
easier. Let's see how the setup is done for our handler integration tests:

```golang
type HttpTestSuite struct {
	suite.Suite
	TestStorage *psql.Storage
	TestDb      *sql.DB
	TestRouter  *fiber.App
	Cfg         *config.Config
}

func (s *HttpTestSuite) SetupSuite() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cfg := config.NewConfig()

	db, err := sql.Open("postgres", cfg.GetDatabaseConnString())
	if err != nil {
		panic(errors.WithStack(err))
	}

	err = db.Ping()
	if err != nil {
		panic(errors.WithStack(err))
	}
	storage := psql.NewStorage()

	userService := service.NewUser(storage)
	weightService := service.NewWeight(storage)

	userHandler := http.NewUserHandler(userService)
	weightHandler := http.NewWeightHandler(weightService)

	srv := http.NewHttp(cfg, *userHandler, *weightHandler)

	srv.SetupRoutes()
	r := srv.GetRouter()

	s.Cfg = cfg
	s.TestDb = db
	s.TestStorage = storage
	s.TestRouter = r
}
```

We basically take the entire setup step and automate it for each test suite. If
we check the documentation for the SetupSuite method we see that it's basically
a method that runs before the test in a suite is run. So the whole setup we did
with the standard library like here:

```golang

func TestIntegration_UserHandler_CreateUser(t *testing.T) {
	cfg := config.NewConfig()
	storage := psql.NewStorage()

    ..... irrelevant code removed

	userService := service.NewUser(storage)
	weightService := service.NewWeight(storage)

	userHandler := http.NewUserHandler(userService)
	weightHandler := http.NewWeightHandler(weightService)

	srv := http.NewHttp(cfg, *userHandler, *weightHandler)

	srv.SetupRoutes()
	r := srv.GetRouter()
    
    ..... irrelevant code removed

}
```

is automated for us, nice! Now, we also did have some other requirements, namely
that we had a "fresh" environment for each test run. This means that we need to
run up/down migrations to ensure our database is clean. This was done in the
setup and teardown portion before in each test, but with testify, we can just
define beforeTest and afterTest where we can run the same methods as we did
before, without having to copy-paste them for each test.

One thing you will notice if you check out the repo is that we have almost the
entirely same code in the repository as we do here. Only except for the
TestRouter in the struct. I don't really mind the duplication here as my needs
for the endpoints test could change in the future, and keeping my dependencies
as few as possible is desirable. So you could, if you wanted, make one large
integration test suite. I just prefer to split things up, each to their
own.

# In conclusion
Will the above steps prevent you from the disaster we went through at the
beginning of the article? Maybe, it depends (every senior developer's favourite
reply). But, it will definitely increase the amount of confidence you can have
in your code. How many integration tests to have is always a balance since they
do take some time longer to run, but there are ways at fixing that, so testing
the happy path until things becomes unbearably slow is a good rule of thumb.
Cross that bridge when you get to it.

As mentioned in the introduction, this is not an attempt to add something new
and revolutionise the way we do integration tests in the Go community. But to
give you another perspective and some boilerplate code to get going with your
integration testing adventures. Should you want to continue and learn from
people way smarter than me (you definitely should), check out the resource
section.

## Resources
[Learn go with tests](https://quii.gitbook.io/learn-go-with-tests/) - basically read this entire thing. Chris does an amazing
job showing you how to get started with Go and test-driven development.
Definitely worth the read.

[HTML forms, databases, integration tests](https://www.lpalmieri.com/posts/2020-08-31-zero-to-production-3-5-html-forms-databases-integration-tests/) - though it's not in Go, but in Rust,
Luca does a great job explaining integration testing. Always try to look for
what concepts transcend programming languages and what doesn't. It's always
beneficial to have a nuanced view.
