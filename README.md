# 🏋️ Coding Dojo: Deep Modules & Tell Don't Ask

A hands-on coding dojo exploring key ideas from **A Philosophy of Software Design** by John Ousterhout. You'll refactor a deliberately "shallow" Go module through three rounds, each introducing a core design principle.

[A Philosophy of Software Design](https://github.com/lamb/a-philosophy-of-software-design)
[Tell, Don't Ask](https://martinfowler.com/bliki/TellDontAsk.html)

## 📖 Background: A Philosophy of Software Design

The book's central thesis is simple: **the primary job of a software designer is managing complexity.** The best way to do that is through good abstractions, clean interfaces, and a strategic mindset.

Here are the core principles:

### Complexity is the root problem

The central enemy of good software design is complexity — anything that makes a system hard to understand or modify. The goal of design is to minimize the complexity that developers must deal with at any given time.

### Deep modules over shallow ones

The best modules provide powerful functionality behind a simple interface. A "deep" module hides a lot of complexity; a "shallow" module has an interface nearly as complex as its implementation, offering little abstraction benefit.

```
SHALLOW                      DEEP
┌──────────────────────┐     ┌───────┐
│      interface       │     │ intf  │
├──────────────────────┤     ├───────┤
│    implementation    │     │       │
└──────────────────────┘     │       │
                             │ impl  │
                             │       │
                             │       │
                             └───────┘
```

### Information hiding

Each module should encapsulate design decisions and knowledge, exposing as little as possible. Information leakage — where the same knowledge is spread across multiple modules — is a major source of complexity.

### Define errors out of existence

Rather than throwing exceptions for edge cases, design APIs so that those cases simply can't occur or are handled naturally. Every exception site adds complexity for callers.

### Tactical vs. strategic programming

Tactical programmers just get features working as fast as possible, accumulating technical debt. Strategic programmers invest a little extra time in good design with each change, which pays off over time.

### Complexity is incremental

No single decision ruins a codebase. Complexity creeps in through hundreds of small choices, so every decision matters — even seemingly minor ones.

### Separate general-purpose and special-purpose code

General-purpose interfaces tend to be simpler and more reusable. Pushing special-case logic upward (or out) keeps lower layers clean.

### Pull complexity downward

It's better for a module's developer to suffer a bit of internal complexity than to push that burden onto every user of the module. Make life easy for the caller.

### Write comments that describe things not obvious from the code

Good comments capture *why* and high-level *what*, not low-level *how*. They should describe abstractions, invariants, and design rationale.

### Design it twice

Before committing to an approach, consider at least one alternative. Comparing designs — even briefly — often reveals a clearly better option.

---

## 🧩 The Connection: Deep Modules & Tell Don't Ask

This dojo focuses on the intersection of two ideas:

**Deep Modules** (Ousterhout): a module should provide powerful functionality behind a narrow interface.

**Tell, Don't Ask** (pragmatic OO principle): instead of *asking* an object for its state and making decisions for it, *tell* the object what you want and let it figure out how.

These reinforce each other naturally. When you stop asking and start telling, behavior moves inside the module. The interface shrinks, the implementation grows richer, and the module gets **deeper**.

**Ask pattern (shallow):**

```go
if !account.GetIsFrozen() {
    if account.GetBalance() - amount >= -account.GetOverdraftLimit() {
        account.SetBalance(account.GetBalance() - amount)
    }
}
```

**Tell pattern (deep):**

```go
err := account.Withdraw(amount)
```

---

## 🎯 The Exercise: Bank Account

You start with a `BankAccount` that is deliberately designed as a **shallow module**. All fields are exposed through getters and setters. All logic lives in free functions outside the struct. The caller must interrogate the account, make every decision, and mutate its state directly.

Your job is to refactor it through three rounds.

### Round 1 — Tell, Don't Ask (~20 min)

Move `Withdraw`, `Deposit`, `Transfer`, `Freeze`, and `Unfreeze` into `BankAccount` as methods.

The caller should go from this:

```go
err := Withdraw(account, 100.0)
err := Transfer(from, to, 50.0)
```

To this:

```go
err := account.Withdraw(100.0)
err := from.Transfer(to, 50.0)
```

Questions to discuss:

- Which getters and setters can you eliminate now?
- How much does the caller still need to know about overdraft logic?
- Where is the validation boilerplate repeated, and can you consolidate it?

### Round 2 — Deep Module (~20 min)

A new requirement arrives: every operation must record a `Transaction` with a timestamp, type, amount, resulting balance, and description.

```go
type Transaction struct {
    Timestamp   time.Time
    Type        string   // "withdrawal", "deposit", "transfer_in", "transfer_out"
    Amount      float64
    Balance     float64
    Description string
}
```

Add a `Transactions()` method that returns the account's history. Then uncomment the **Round 2 tests** and make them pass.

Questions to discuss:

- In your Round 1 design, where does the logging naturally go?
- In the original code, how many free functions would you need to modify?
- Which design absorbs this new requirement more easily, and why?

This is the deep module payoff: a narrow interface that hides growing internal complexity.

---

## 🚀 Getting Started

```bash
git clone <repo-url>/go-deep-modules-kata.git
cd go-deep-modules-kata
go test ./...
```

All tests pass against the "before" code. Refactor round by round, updating the test call syntax as you go. Uncomment the Round 2 and Round 3 tests when you reach those stages.

### Project Structure

```
.
├── README.md
├── go.mod
└── pkg/
    ├── account/
    │   ├── account.go       # The shallow module — your starting point
    │   └── account_test.go  # Tests (includes commented sections for Rounds 2 & 3)
    └── order/
        ├── order.go         # Bonus: shallow order pipeline (Cart, Inventory, Pricer, Order)
        └── order_test.go    # Tests (includes commented sections for Rounds 1–3)
```

---

## ⏱️ Suggested Dojo Timeline (~55 min)

| Time  | Activity |
|-------|----------|
| 0:00  | **Intro** — Walk through the shallow vs deep diagram. Show the before code. Ask: *"What smells do you see?"* |
| 0:15  | **Round 1** — Pairs refactor: move logic into methods, eliminate getters/setters |
| 0:35  | **Debrief** — Compare solutions. Key insight: the caller no longer knows about overdraft logic |
| 0:40  | **Round 2** — New requirement: transaction history. Pairs add it |
| 0:55  | **Final retro** — Key insight: logging goes in ONE place now. What smells will you watch for going forward? What heuristics do you take back to your codebase? |

---

## 🚀 Bonus: Order Processing Pipeline

Finished early? This stretch challenge applies the same principles across **multiple interacting modules**, where information leakage *between* them becomes the central problem.

### The domain

```
┌────────────┐     ┌────────────┐     ┌────────────┐
│    Cart     │────▶│  Checkout  │────▶│ Fulfillment│
│ (items,    │     │ (pricing,  │     │ (shipping, │
│  quantities)│     │  payment)  │     │  tracking) │
└────────────┘     └────────────┘     └────────────┘
```

### Shallow starting point

One massive orchestration function reaches into every module:

```go
func PlaceOrder(cart *Cart, inventory *Inventory, pricer *Pricer, payment *Payment, shipper *Shipper) error {
    for _, item := range cart.GetItems() {
        if inventory.GetStock(item.GetSKU()) < item.GetQuantity() {
            return errors.New("out of stock")
        }
        inventory.SetStock(item.GetSKU(), inventory.GetStock(item.GetSKU()) - item.GetQuantity())
    }
    total := 0.0
    for _, item := range cart.GetItems() {
        total += pricer.GetPrice(item.GetSKU()) * float64(item.GetQuantity())
    }
    if err := payment.Charge(cart.GetCustomerID(), total); err != nil {
        return err
    }
    shipper.Ship(cart.GetCustomerID(), cart.GetItems())
    return nil
}
```

Adding a discount rule means touching `PlaceOrder`, the cart, *and* the pricer. The function knows the internal structure of every module.

### Apply the same two rounds

**Round 1 — Tell Don't Ask:** `cart.Checkout(inventory)` returns an `Order`. The cart validates stock internally. The caller no longer knows about inventory levels.

**Round 2 — Deep Modules:** Add an order state machine (`pending → paid → shipped → delivered`) with an event log. Each transition records what happened. The caller just says `order.Pay(payment)` and `order.Ship(address)`.

### The key insight

In the bank kata, one module got deeper. Here, *three* modules get deeper, and the real win is that information stops leaking *between* them — changing pricing logic no longer requires touching fulfillment code.

---

## 🌍 Deep Modules in the Real World

These principles aren't academic — they're behind some of the most successful abstractions in computing.

### Unix file descriptors

The file descriptor API is Ousterhout's canonical example of a deep module. The interface is just five calls: `open`, `close`, `read`, `write`, `lseek`. Behind that tiny surface the kernel hides disk drivers, buffer caches, file systems, permissions, block allocation, and device multiplexing. Callers don't know — or care — whether they're talking to an SSD, a network socket, or a pipe.

### Go's `io.Reader` / `io.Writer`

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

One method. Behind it: files, HTTP bodies, gzip streams, TLS connections, in-memory buffers, and anything else that produces bytes. The entire Go I/O ecosystem composes through this single interface — `json.NewDecoder(r)` works the same whether `r` is an `os.File` or an `http.Response.Body`. This is depth through abstraction.

### TCP/IP sockets

Application code calls `send` and `recv`. The stack hides: packet segmentation, retransmission, flow control, congestion avoidance, routing, checksums, and reassembly. Entire RFCs worth of complexity behind a stream-oriented interface.

### `database/sql` in Go

```go
db.QueryRow("SELECT name FROM users WHERE id = ?", 42).Scan(&name)
```

One package, one interface. Behind it: connection pooling, prepared statement caching, driver negotiation, transaction isolation, and the ability to swap Postgres for MySQL by changing one import. The caller never manages connections directly.

### `net/http` in Go

```go
http.HandleFunc("/hello", handler)
http.ListenAndServe(":8080", nil)
```

Two lines to start an HTTP server. Hidden: TLS handshakes, keep-alive management, chunked transfer encoding, header parsing, graceful shutdown, HTTP/2 negotiation. Adding HTTPS is a one-line change (`ListenAndServeTLS`), not a rewrite.

### `crypto/tls`

Callers call `tls.Dial` or wrap a `net.Conn` with `tls.Client`. Behind that: certificate validation chains, cipher suite negotiation, key exchange, session resumption, OCSP stapling, and protocol version handling. The caller gets a `net.Conn` back — the same interface as a plain TCP connection.

### `context.Context`

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key any) any
}
```

Four methods. Carries deadlines, cancellation signals, and request-scoped values through the entire call chain. The HTTP handler creates it, middleware enriches it, the database driver respects it, and the gRPC client propagates it — all without any of those layers knowing about each other.

### `encoding/json`

```go
data, err := json.Marshal(myStruct)
err := json.Unmarshal(data, &myStruct)
```

Two functions for the common case. Behind them: reflection over struct fields, struct tag parsing (`json:"name,omitempty"`), type switching, streaming tokenizer, Unicode escaping, custom `Marshaler`/`Unmarshaler` interface dispatch, and handling of nested/recursive types. The caller just passes a value and gets bytes back.

### `sort.Interface`

```go
type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}
```

Three methods to sort *anything*. The implementation hides: quicksort, heapsort, insertion sort, pivot selection heuristics, and the decision of when to switch between them (`pdqsort` since Go 1.19). The caller defines ordering; the package handles all algorithmic complexity.

### `sync.WaitGroup`

```go
var wg sync.WaitGroup
wg.Add(1)
go func() { defer wg.Done(); work() }()
wg.Wait()
```

Three methods: `Add`, `Done`, `Wait`. Behind them: atomic counters, OS-level semaphores, memory barriers, and goroutine scheduling coordination. The caller never thinks about futexes or memory ordering.

### `bufio.Scanner`

```go
scanner := bufio.Scanner(reader)
for scanner.Scan() {
    line := scanner.Text()
}
```

One loop pattern to consume any stream line by line. Hidden: internal buffering, dynamic buffer growth, configurable split functions (lines, words, bytes, runes), EOF handling, and error accumulation. Replacing line-based scanning with word-based scanning is a one-line change (`scanner.Split(bufio.ScanWords)`).

### Notable Go projects

#### Docker (`github.com/docker/docker`)

```go
cli.ContainerCreate(ctx, config, hostConfig, networkConfig, nil, "my-app")
cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
```

Two calls to run a container. Hidden: Linux namespaces (PID, network, mount), cgroups for resource limits, overlay filesystem layering, image pulling and caching, port mapping with iptables, DNS resolution, volume mounts, and seccomp profiles. The caller says "run this image" — the engine handles everything from kernel primitives to networking.

#### Kubernetes client-go (`k8s.io/client-go`)

```go
clientset.CoreV1().Pods("production").Get(ctx, "my-pod", metav1.GetOptions{})
```

One method chain. Behind it: REST transport, bearer token or certificate authentication, API version negotiation, JSON/protobuf serialization, retry with exponential backoff, rate limiting, and watch/informer caching. Switching from in-cluster auth to kubeconfig is a config change, not a code change.

#### gRPC-Go (`google.golang.org/grpc`)

```go
conn, _ := grpc.Dial("server:443", grpc.WithTransportCredentials(creds))
client := pb.NewGreeterClient(conn)
resp, _ := client.SayHello(ctx, &pb.HelloRequest{Name: "world"})
```

Looks like a local function call. Hidden: HTTP/2 framing, protobuf serialization, TLS handshake, connection pooling, client-side load balancing, interceptor chains, deadline propagation, retry policies, and health checking. Adding a unary interceptor for logging or metrics is one option — the RPC interface doesn't change.

#### GORM (`gorm.io/gorm`)

```go
db.Where("age > ?", 18).Preload("Orders").Find(&users)
```

A fluent query. Hidden: SQL dialect differences (Postgres vs MySQL vs SQLite), connection pooling, prepared statement caching, eager loading with separate queries, struct-to-column mapping via reflection, soft deletes, callbacks/hooks, and transaction management. Switching databases is a driver import change.

#### Zap (`go.uber.org/zap`)

```go
logger.Info("user signed in",
    zap.String("user", userID),
    zap.Duration("latency", elapsed),
)
```

One call per log line. Hidden: zero-allocation encoding, field pooling, log level sampling (emit only 1 in N for high-frequency events), multiple output sinks, JSON or console formatting, atomic level changes at runtime, and buffered async writes. All of that behind `logger.Info`.

#### go-redis (`github.com/redis/go-redis`)

```go
rdb.Set(ctx, "key", "value", 10*time.Minute)
rdb.Get(ctx, "key")
```

Simple key-value calls. Hidden: connection pooling with health checks, automatic pipelining, cluster-aware slot routing and redirection (`MOVED`/`ASK`), Sentinel failover, pub/sub multiplexing, Lua script caching, and context-based timeouts. Switching from standalone Redis to a cluster is a config change.

### The pattern

Every example shares the same shape: **a narrow interface hiding enormous implementation complexity.** That's the deep module idea in practice. When you refactor your `Account` from getters/setters into `account.Withdraw(amount)`, you're applying the exact same principle that makes `io.Reader`, Docker's client API, and Kubernetes client-go so powerful.

---

## 🔍 Smells to Watch For

These are signs you're looking at a shallow module with Tell Don't Ask violations:

- **Getter chains followed by logic followed by setters** — the object is a data bag, not a module
- **Repeated validation boilerplate** — every caller re-checks the same preconditions
- **Functions that know the internals of multiple objects** — like `Transfer` reaching into two accounts
- **Adding a feature requires touching many call sites** — the logic isn't centralized behind an interface
- **The interface is as wide as the implementation** — public getters/setters for every field, hiding nothing

---

## 📚 Further Reading

- **A Philosophy of Software Design** — John Ousterhout
- **Tell, Don't Ask** — Martin Fowler ([martinfowler.com/bliki/TellDontAsk.html](https://martinfowler.com/bliki/TellDontAsk.html))
- **The Pragmatic Programmer** — Hunt & Thomas (for more on pragmatic design principles)
