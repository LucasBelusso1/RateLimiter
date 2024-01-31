# Rate limiter written in GO

### How is it works?

The Rate limiter uses a instance of Redis to control the request limit for a specified endpoint. The ideia of the
middleware is to store a key (ip or token) inside Redis and controls the limit by the key. For example, if a request is
made with a token and the token can make 10 requests per minute, on the first request a key will be created on redis
using the token as the index and the request count as the value. The key that was created has an expires defined by the
token expires configuration on `.env`.

To create a new rate limiter middleware, it's necessary to provide the IP request limit and the time for the IP key to
be expired.

### Strategy

The middleware uses a design pattern called `strategy`, that enables selecting an algorithm at runtime, it means that
it's possible to change between Redis and another In memory database solution (or any DB solution) by just creating a
new strategy.
The main ideia is to have a "context" that recieves a strategy interface. From the context, the strategy is
called and the function of the strategy is executed.

### Limiters

It's possible to create **N** rules for **N** tokens, you just need to instanciate `Limiters` specifing the `context`,
the value of the field to be validated, the expires of the redis key and the request limit:

```GO
limiters := []limiter.Limiter{
	{
		DbContext:    context,
		Field:        config.TokenAName,
		TimeLimit:    config.TimeForTokenA,
		RequestLimit: config.TokenALimit,
	},
}

middleware.NewMiddleware(context, limiters, {{IpLimit}}, {{TimeForIp}})
```

### How can i test it?

From the repository, with `docker` and `docker compose` installed, you just need to enter inside `/.docker` and run:

```SHELL
docker compose up -d
```

It will run a docker container with redis, now you need do go to `/cmd` and run:

```SHELL
go run main.go
```

To execute the tests, from `/` run:

```SHELL
go test ./...
```