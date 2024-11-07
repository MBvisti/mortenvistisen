My approach to working with data in Golang has changed a bit over the years. I remember the early days where it was a lot of string concatenation while raw dogging sql, trying not to introduce injection vulnerabilities.

Sqlx was a great improvement but you still had an awkward situation of actually writing the sql. 

I've seen it done directly in the code like:
```go
func GetArticle(id uuid.UUID, db *sql.DB) (Article, error) {
    stmt := `select * from articles where id=$1`

    var article Article
    row := db.Query(stmt, id)

    if err := row.Scan(&article.ID, &article.Content); err != nil {
        return Article{}, err
    }

    return article, nil
}
```
To writing the queries in an actual sql file, loading in said file and then execution the query.

It wasn't pretty but it worked, sort of.

ORMs has long been villainized in the Go community even though we've big projects like GORM. My first encounter early in my career was with Rails ActiveRecord ORM, it never really did it for me. I much prefered to write pure sql.

And this is where Sqlc comes in. I first stumbed across this library during a constract I had with a web3 company that used it, and it just fits the bill so perfectly.

You get to write all of your query logic in sql and have it compiled to type-safe Go code. You eliminate so much boilerplate you'd other have to write. 

You pretty much install the CLI, create a `sqlc.yaml` file with some configuration and then just start creating queries:
```sql
-- name: QueryArticleByID :one
select * from articles where id=$1;
```

Running `sqlc generate` will generate a method on a Queries struct that has a 'QueryArticleByID' method, and works just like a typical method.

This will probably be enough for 8/10 queries you have. You get to be 100% in charge of how you handle data access and can much easier optimize your queries if you need to.

Since Sqlc validates all queries against your migrations you can even chuck large views into the equation, giving you efficiency and readability with little effort. 

If you're building a JSON based data API, which is like 99% of people working with Go, Sqlc is the way to go. The removal of boilerplate makes the developer experience and speed go through the roof. Whatever gains your ORM gave you will be short lived.

Easy 'queries' are now just as quick to make compared to an ORM like GORM while the complex ones now are out there in the open, for you to tackle if you need and not hidden away in the GORM execution layer somewhere.

But, are there any downsides? Sqlc doesn't handle dynamic queries well. You can 'hack' your way out of some scenarios. But it's just, for now at least, that good at it.

A tool like squrriel are what you want here. And since Sqlc compiles to Go code, you can just create a method on the Queries struct for your dynamic queries and use squrriel there.
