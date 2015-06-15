package container


// Depender provides the Dependencies method.
//
// The Depender interface comes into play with the Container's
// Inject method.
//
// Normally the Container's Inject method can only do injections
// into a structs public fields. (I.e., field name that begin with
// a capitl letter.)
//
// However, to keep thing clean, a struct may want to "hide" its
// dependencies.
//
// To make such a thing possible, if what is passed (as an argument)
// to Container's Inject method fits a Depender (i.e., it has a
// Dependencies method, of the proper signature) then the Container's
// Inject method will try to (also) inject whatever is returned
// from the Dependencies methods.
//
// The motivation for this is that a struct could have its dependencies
// in a non-public location.
//
// Consider the following example:
//
//	type cherryDependencies struct {
//		Logger  *log.Logger    `inject:"logger"`
//		Emailer *email.Emailer `inject:"emailer"`
//	}
//	
//	type Cherry struct {
//		dependencies cherryDependencies
//		Weight float64
//	}
//	
//	func NewCherry(weight float64) *Cherry {
//		cherry := new(Cherry)
//
//		cherry.Weight = weight
//
//		return cherry
//	}
//	
//	func (c *Cherry) Dependencies() interface{} {
//		// Notice that we returned a pointer to
//		// the dependencies!
//		return &c.dependencies
//	}
//
// Then if we have an instance of Cherry, as in:
//
//	// Note that the NewCherry() func returns a *Cherry
//	// and not just Cherry.
//	cherry := NewCherry(10.2)
//
// Then if we call the following:
//
//	Container.Inject(cherry)
//
// Notice that the Cherry struct's dependencies were all under
// the dependencies field.
//
type Depender interface {
	Dependencies() interface{}
}
