/*
Package container provides 'dependency injection' functionality.

A large part of this is done through the 'dependency injection container', which
is reified through the 'container.Container' interface, and created by calling the
'container.New()' func. As in:

	Container := container.New()

You do 3 things with a 'dependency injection container'.

(Although it is possible that in practice you might only do 2 of these 3 things.
But if you end up doing all 3 of these things, that's OK too.)

#1: You register things with the container, giving the thing you registed a
(string) name when you register it.

For example:

	// Register a service.
	if err := Container.Register("logger", myLogger); nil != err {
		//@TODO: Handle an error better than this!
		panic(err)
	}

Also, for example:

	// Register a config value.
	if err := Container.Register("pool-size", 20); nil != err {
		//@TODO: Handle an error better than this!
		panic(err)
	}

#2: You manually get things out of the container.

(Doing this is a bit involved, but #3 provides a better alternative to this method.)

For example:

	serviceName := "logger"
	service, err := Container.Get(serviceName)
	if nil != err {
		//@TODO: Handle an error better than this!
		panic(err)
	}
	
	var ok bool
	var logger *log.Logger
	logger, ok = service.(*log.Logger)
	if !ok {
		//@TODO: Handle an error better than this!
		panic("Got service, but it was the wrong type!")
	}

Also, for example:

	configValueName := "pool-size"
	configValue, err := Container.Get(configValueName)
	if nil != err {
		//@TODO: Handle an error better than this!
		panic(err)
	}
	
	var ok bool
	var poolSize int
	logger, ok = configValue.(int)
	if !ok {
		//@TODO: Handle an error better than this!
		panic("Got service, but it was the wrong type!")
	}

#3: The container can inject dependencies into something.

For example:

	type workerPoolDependencies struct {
		Logger   *log.Logger `inject:"logger"`
		PoolSize  int        `inject:"pool-size"`
	}
	
	type WorkerPool struct {
		dependencies workerPoolDependencies
		Name string
	}
	
	func (wp *WorkerPool) Dependencies() interface{} {
		// Notice that we are returning a POINTER to
		// the dependencies here. That's important!
		return &wp.dependencies
	}
	
	
	Container.Inject(commandHandler)

*/
package container
