# Shipping Rate CLI

## Purpose

A Go CLI application that calculates parcel shipping costs in EUR from the parcel's weight and destination zone.

## How It Works

1. The CLI reads the required weight, length, width, height and zone flags into a `Parcel`, converting the zone name to lowercase.
2. `ParcelService` validates the parcel and looks up the destination zone.
3. The shipping rate calculator selects a weight-based rate and adds the zone surcharge.
4. The CLI prints the total, or returns an error if the input is invalid.

Base rates are EUR 5 for weight up to 1 kg, EUR 10 for weight up to 5 kg, and EUR 20 above 5 kg. Zone surcharges are EUR 0 for `domestic`, EUR 8 for `eu`, and EUR 15 for `international`.

Weight is in kilograms; dimensions are in centimetres. All measurements must be finite and greater than zero. Dimensions are validated but do not affect the price.

## Packages

| Package | Purpose |
| --- | --- |
| `spf13/cobra` | CLI commands and flags |
| `juju/errors` | Error creation and annotations |
| `stretchr/testify` | Test suites, assertions and mock expectations |
| `vektra/mockery/v3` | Generate mocks from interfaces via `make mock-gen` |

## Architectural Patterns

- **Simplified DDD:** The domain contains parcel and zone entities, the `Money` value object and currency enums. Application services coordinate validation, zone lookup and pricing.
- **Onion Architecture:** Domain models sit at the centre, surrounded by application services, infrastructure repositories and presentation commands. Services use interfaces for their dependencies. The implementation is pragmatic: the `NewServices` wiring helper references the infrastructure repository container.

The CLI entry point is `cmd/cli`. A future web API can use `cmd/web-api` with separate presentation handlers and reuse the application services and domain models.

## Code Patterns

- **Dependency Injection:** Constructors receive dependencies through interfaces; `NewRepositories` and `NewServices` assemble them.
- **Repository:** `IZoneRepository` separates zone lookup from its in-memory implementation.
- **Strategy:** `IShippingRateCalculator` allows the pricing algorithm to be replaced without changing `ParcelService`. Currently, there is one implementation.
- **Value Object:** `Money` groups the amount and currency and checks currency compatibility when adding amounts.

## Development with Codex and Claude Code

`AGENTS.md` contains the project context and development conventions used by Codex. `CLAUDE.md` is a relative symlink to `AGENTS.md`, making the same instructions available to Claude Code.

Update `AGENTS.md` to keep both tools aligned with the project.

## Makefile Commands

On macOS, install GNU Make with Homebrew:

```sh
brew install make
```

Homebrew installs it as `gmake`. Use `gmake build`, `gmake test` and `gmake mock-gen`, or add it to your PATH to use the `make` commands below:

```sh
export PATH="$(brew --prefix make)/libexec/gnubin:$PATH"
```

| Command | Purpose |
| --- | --- |
| `make build` | Build the CLI executable at `bin/shipping-rate` |
| `make test` | Run all tests with race detection and coverage |
| `make mock-gen` | Regenerate interface mocks using the pinned Mockery version |

## Quickstart

Requires Go 1.26 or later and Make for the Makefile commands.

Run directly from the project directory:

```sh
go run ./cmd/cli calculate --weight 3 --length 30 --width 20 --height 10 --zone EU
```

Or build and run the executable:

```sh
make build
./bin/shipping-rate calculate --weight 3 --length 30 --width 20 --height 10 --zone eu
```

Both commands print:

```text
EUR 18.00
```

For available flags, run `go run ./cmd/cli calculate --help`.
