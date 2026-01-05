// The content packages contains the domain definitions of each resource. It is not
// intended to change and should not import any other project packages. Each type has
// it's own validation rules, exposed as a method of that type.
//
// The content package is only concerned with the business logic of resources. It
// does not care about whether other resources exist, and validation functions
// should reflect this.
package content
