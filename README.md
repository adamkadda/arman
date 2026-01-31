# arman

Arman is short for **ar**tist **man**agement software, its a content management system with an admin dashboard.

> Arman is a deliberate rewrite of an earlier artist management CMS, undertaken to improve architecture, correctness, and long-term maintainability.  
This repository reflects the current, vetted implementation; features from earlier iterations are being re-introduced incrementally as they are redesigned and validated against the new goals.  
[The original implementation can be found here](https://github.com/adamkadda/ntumiwa/tree/rewrite).

At the moment, here is what Arman is primarily built with:
- [net/http](https://pkg.go.dev/net/http) for HTTP routing and handling
- [jackc/pgx](https://pkg.go.dev/github.com/jackc/pgx/v5) for its PostgreSQL driver and toolkit for Go
- [google/uuid](https://pkg.go.dev/github.com/google/uuid) for logging request IDs
- [caarlos0/env](https://pkg.go.dev/github.com/caarlos0/env/v11) for loading env variables with struct tags
- [stretchr/testify](https://pkg.go.dev/github.com/stretchr/testify) for its testing tools

I'm sure that list will steadily grow, but I'm aiming to keep things light.

Here are some features arman aims to have:
- Streamlined management of classical repertoire (composers & pieces)
- Simplified programme creation
- Practical event management (drafting, publishing, and archiving)

The feature set will evolve, but the current focus is on building a solid, maintainable foundation.
