# Your 8-Month Roadmap to Infrastructure Engineering

## Month 1-2: Container Mastery

**Build:** Microservices app (3+ services + database + cache)

- Frontend, API gateway, 2 backend services, Postgres, Redis
- Use Docker Compose to orchestrate everything

**Learn:**

- Stateless vs stateful workloads (this is CRITICAL for interviews)
- Persistent volumes (how data survives container restarts)
- Service discovery (how containers find each other)
- Health checks and auto-restart
- Network isolation between services
- Resource limits (CPU/memory)

**Deliverable:** Blog post “How I Built a Production-Ready Microservices App with Docker”

-----

## Month 3-4: Kubernetes Basics

**Build:** Deploy same app to local Kubernetes (minikube)

- Translate Docker Compose → K8s manifests
- Set up Deployments, Services, StatefulSets
- Configure persistent storage for database
- Add health/readiness probes

**Learn:**

- How orchestrators schedule workloads
- Difference between Deployments (stateless) and StatefulSets (stateful)
- Service discovery in K8s
- ConfigMaps and Secrets
- Scaling replicas
- Rolling updates

**Experiment:** Kill pods, drain nodes, watch K8s recover automatically

**Deliverable:** GitHub repo with full K8s setup + README explaining choices

-----

## Month 5-6: Golang + Basic Systems

**Build Part 1:** Rewrite one backend service in Go

- REST API with database
- Graceful shutdown
- Structured logging
- Health endpoints

**Build Part 2:** Simple process manager in Go

- Start/stop processes
- Monitor health
- Restart on failure
- Set resource limits (cgroups)

**Learn:**

- Go fundamentals (goroutines, channels, error handling)
- Working with databases in Go
- Process management basics
- Linux cgroups for resource control

**Resources:** [tour.golang.org](http://tour.golang.org), “Let’s Go” book

**Deliverable:** Two GitHub repos showing Go proficiency

-----

## Month 7-8: Mini Orchestrator (Interview Project)

**Build:** System that manages stateless and stateful workloads

**Components:**

1. **API Server** - accepts workload requests (gRPC or REST)
1. **Scheduler** - decides where to run workloads based on:

- Available resources
- Stateless: can run anywhere
- Stateful: needs persistent storage

1. **Node Agent** - runs on each “node”, manages local workloads:

- Starts processes in isolated namespaces
- Attaches volumes for stateful apps
- Monitors health

1. **Storage Manager** - handles persistent volumes
1. **Network Manager** - sets up service discovery

**Learn:**

- Linux namespaces for isolation
- Scheduling algorithms
- State management in distributed systems
- gRPC for service communication
- How to answer Railway’s interview question

**Deliverable:** Working prototype + architecture doc explaining design decisions

-----

## Continuous (Whole Journey)

**Follow & Learn:**

- [Fly.io](http://Fly.io) engineering blog
- Railway blog
- PostHog engineering posts
- Kubernetes documentation
- “Container From Scratch” talk by Liz Rice (YouTube)

**Practice:**

- [Fly.io](http://Fly.io)’s public hiring challenges on GitHub (fly-hiring org)
- Contribute to open source (Docker, K8s, or related tools)

-----

## The Interview Prep

**For Railway/Fly.io/PostHog interviews:**

**You’ll discuss:**

- How you’d design a runtime for stateless vs stateful workloads
- Scheduling decisions (where to place workloads)
- Storage strategies (local volumes vs distributed storage)
- Networking (service discovery, load balancing)
- Failure handling (health checks, restarts, data recovery)

**Your Answer Framework:**

1. **Stateless:** “Like web servers - can restart anywhere, no data persistence needed, easy horizontal scaling”
1. **Stateful:** “Like databases - need persistent storage, careful scheduling, graceful shutdown critical”
1. **Reference your projects:** “In my mini orchestrator, I handled this by…”

-----

## Timeline Checkpoints

**Month 2:** Docker expert, understand containers deeply
**Month 4:** Can deploy complex apps on Kubernetes
**Month 6:** Comfortable writing Go, understand Linux process management
**Month 8:** Built mini orchestrator, ready to interview

-----

## Where to Apply

**Month 7-8:**

1. **Fly.io** (best fit - they love discovering talent)
1. **PostHog** (multiple infrastructure roles)
1. **Railway** (original goal)

**Backup:**

- DigitalOcean (App Platform)
- Cloudflare (Workers platform)
- HashiCorp (infrastructure tools)

-----

## Success Metrics

✅ Can explain stateless vs stateful from first principles
✅ Built something involving containers, orchestration, and Go
✅ Understand Linux namespaces, cgroups, and process isolation
✅ Can discuss trade-offs in distributed systems
✅ Have public GitHub repos showing your work
✅ Wrote blog posts teaching others (shows deep understanding)

-----

## The Key Insight

**These companies don’t care about your resume.** They care if you can:

- Solve real infrastructure problems
- Learn quickly
- Ship working code
- Think about systems holistically

Your projects + blog posts + GitHub = your resume.

-----

**Start Date:** This weekend
**End Date:** 8 months from now
**Outcome:** Interview-ready for infrastructure engineering roles at top startups

Go build something awesome. 🚀​​​​​​​​​​​​​​​​
