# Contact Manager

A small command-line contact manager in Go. It stores contacts in memory and exposes
add / list / find / delete through an interactive menu.

The program is deliberately simple, but it is layered the way a larger Go service would
be. This README explains the concepts the code demonstrates, since that is the point of
the project.

## Running it

```bash
go run .
```

Requires Go 1.25.2 or newer (see `go.mod`). There are no third-party dependencies — only
the standard library.

```
===== CONTACT MANAGER =====
1. Add Contact
2. List Contacts
3. Find Contact
4. Delete Contact
5. Exit
Enter choice:
```

Contacts live in memory only, so they are gone when the process exits.

## Layout

```
main.go                            CLI: menu loop and input handling
model/contact.go                   the Contact struct
repository/contact_repository.go   in-memory storage
service/contact_service.go         validation and business rules
storage/contact_store.go           interface for future file persistence (unimplemented)
go.mod                             module declaration
```

## The core idea: dependencies point one way

```
main  ──▶  service  ──▶  repository  ──▶  model
                                            ▲
                            storage ────────┘
```

Each layer knows only the one beneath it, and no arrow ever points back up. `model`
imports nothing. `repository` imports `model`. `service` imports `repository`. `main`
imports `service`.

The payoff: swapping in-memory storage for a file or a database touches `repository`
only. Nothing above it needs to change, because nothing above it knows how storage
works. That is the single most useful thing to take from this codebase.

## Concepts, layer by layer

### `go.mod` — the module path is the import prefix

```go
module contact-manager
```

Every internal import is built from this line: `contact-manager/model`,
`contact-manager/service`. It is **not** derived from the directory name, which here is
`contact-manager-golang`. Mismatching the two is a common source of
"package ... is not in std" errors.

### `model` — plain data

`Contact` is a struct with no methods and no imports: `ID`, `Name`, `Phone`. `ID` is
assigned by the repository, not by the caller, so code that creates a `Contact` leaves it
zero.

### `repository` — encapsulation and slice mechanics

**Unexported fields.** `contacts` and `nextID` are lowercase, so nothing outside the
package can reach them. Callers must go through the methods. This is Go's encapsulation,
and it is the entire reason this layer exists as a separate package.

**The constructor idiom.** Go has no constructors. The convention is a `NewX` function
returning `*X`:

```go
func NewContactRepository() *ContactRepository {
	return &ContactRepository{
		contacts: make([]model.Contact, 0),
		nextID: 1,
	}
}
```

**Pointer receiver vs value parameter.** `Add` shows both in seven lines:

```go
func (r *ContactRepository) Add(contact model.Contact) model.Contact {
	contact.ID = r.nextID
	r.nextID++
	r.contacts = append(r.contacts, contact)
	return contact
}
```

`r` is a *pointer* receiver, so `r.nextID++` mutates the real repository. `contact` is a
*value* parameter, so `contact.ID = ...` mutates a local copy — the caller's struct is
untouched. The modified copy is what gets appended and returned, which is why callers
read the new ID off the return value rather than off their own variable.

**The comma-ok pattern.** `GetByID` returns `(*model.Contact, bool)`, the Go idiom for
"here is the thing, and here is whether it was found", rather than returning nil and
making the caller test for it.

**Removing from a slice.** `Delete` uses the standard idiom:

```go
r.contacts = append(r.contacts[:i], r.contacts[i+1:]...)
```

Everything before index `i`, then everything after it, spread back in. Note that `[i+1:]`
is a *slice*; `[i+1]` would be a single element and would not compile with `...`.

**IDs are never reused.** `nextID` only increments, so deleting contact 2 does not free
that ID. Contacts are identified stably for the life of the process.

### `service` — rules and translation

**Dependency injection by hand.** `ContactService` holds a
`*repository.ContactRepository`, passed in through `NewContactService`. The service never
creates its own repository, so the caller decides what it talks to.

**Validation lives here, not in the repository.** `AddContact` rejects an empty name or
phone. The repository's job is to store what it is given; deciding what is *valid* is a
business rule, and business rules belong in one place.

**Translating between layers.** `GetContact` is the clearest example of what this layer is
actually for:

```go
contact, found := s.repository.GetByID(id)
if !found {
	return model.Contact{}, errors.New("Contact not found")
}
return *contact, nil
```

It converts the repository's storage-shaped answer (`pointer, bool`) into a caller-shaped
one (`value, error`). The `*contact` dereference copies the struct, which also insulates
callers from the repository's internal slice.

### `main` — the CLI

An infinite `for` loop prints the menu, reads a line, and dispatches on a `switch`.
Returning from `main` on choice 5 exits the program.

**Reading input.** One `bufio.Reader` is created over `os.Stdin` and passed to every
handler. Each read is `reader.ReadString` on a newline delimiter followed by
`strings.TrimSpace` — the delimiter is included in the returned string, so trimming is
required, not optional. Numeric input then goes through `strconv.Atoi`, whose error is how
invalid input is detected.

**Handler shape.** All four handlers follow one pattern: prompt, read, trim, convert if
numeric, call the service, print the error or the result. The CLI never touches the
repository and never validates anything itself; it only collects strings and displays what
comes back.

**A wrinkle worth noticing.** In `main`:

```go
repository := repository.NewContactRepository()
```

The variable shadows the package name. It compiles and works, but after this line
`repository` refers to the variable, so the package is unreachable for the rest of the
function. Worth understanding; better not to imitate.

### `storage` — an interface as a plan

```go
type ContactStore interface {
	Save(contact model.Contact) error
	Load() ([]model.Contact, error)
}
```

Nothing implements or imports this yet. It is a declaration of intent for file-backed
persistence. In Go, interfaces are satisfied implicitly — any type with these two methods
is a `ContactStore`, with no `implements` keyword and no change to the interface.

Note that the method names are exported (capital `S`, capital `L`). An interface with a
lowercase method could only ever be satisfied by a type declared inside `package storage`,
which would defeat the purpose.

## Known limitations

These are real, and worth knowing before extending the code.

- **No persistence.** Contacts exist only in memory; `ContactStore` is unimplemented.
- **No update operation.** You can add, read, and delete, but not edit a contact.
- **`GetByID` returns a pointer into the backing array.** A later `Add` can reallocate the
  slice and a later `Delete` shifts elements, either of which leaves such a pointer
  referring to the wrong contact. It is safe as currently used, because `GetContact`
  dereferences immediately, but it is a trap for any new caller.
- **`Delete` does not clear the vacated tail slot**, so the removed contact's strings stay
  reachable and cannot be garbage collected until the slot is overwritten.
- **`ContactStore` is asymmetric.** `Save` takes one contact while `Load` returns all of
  them, which does not fit a file that has to be rewritten wholesale.
- **Not concurrency-safe.** The repository has no mutex. Fine for a single CLI goroutine,
  not fine if this ever becomes an HTTP handler.
- **Error strings are capitalized**, against Go convention (`staticcheck` ST1005), since
  error text is usually wrapped into longer sentences.
- **No tests.**

## Exercises

1. Implement `ContactStore` with a JSON file. You will have to decide what to do about the
   asymmetric `Save`/`Load` signatures — working that out is the point of the exercise.
2. Add `UpdateContact` end to end, through all four layers.
3. Change `GetByID` to return `(model.Contact, bool)` by value and watch the pointer
   aliasing problem disappear.
