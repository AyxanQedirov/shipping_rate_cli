# Project context

We are developing a Go CLI application that calculates parcel shipping rates in
EUR. The CLI uses Cobra. The Go module is `shipping_rate_cli`, and `go.mod`
requires Go 1.26.0. Read `README.md` for usage and input assumptions.

## Structure

```text
cmd/cli/main.go                          Entry point
internal/domain/entities/                Parcel, Zone
internal/domain/value-objects/           Money value object
internal/domain/enums/                   Currency enum
internal/application/abstractions/       Service and repository interfaces
internal/application/services/           Services and pricing implementation
internal/infrastructure/repositories/    In-memory tariff repositories
internal/presentation/commands/          Cobra commands
```

Keep each entity and each interface in its own file. Do not consolidate them
into `models.go` or `interfaces.go`.

The current executable entry point is `cmd/cli/main.go`. Entry points under
`cmd/` are named by presentation type; a future web API would use `cmd/web-api/`
and separate presentation handlers while reusing the shared services and models.
The CLI binary remains named `shipping-rate`.

## Initialization and execution flow

1. `main.go` calls `repositories.NewRepositories()`.
2. `NewRepositories()` creates the zone repository and exposes it through an
   interface-typed field on `Repositories`.
3. `services.NewServices(repos)` creates the shipping rate calculator and parcel
   service and exposes them through interface-typed fields on `Services`.
4. `commands.NewCommand(appServices.ParcelService)` accepts `IParcelService`.
5. The root command registers `NewCalculateCommand(service)`. Its `calculate`
   subcommand binds the input flags to an `entities.Parcel` and calls the service.
6. `ParcelService` validates the parcel and resolves its zone using
   `IZoneRepository.FindByName()` before invoking the injected calculator.
7. `ShippingRateCalculator.calculateBaseRate()` selects the weight tier with
   explicit if checks. `Calculate()` adds the zone surcharge with `Money.Add()`
   and propagates currency mismatch errors.
8. The command prints the result, for example `EUR 18.00`.

Interfaces live in `application/abstractions`: `IParcelService`,
`IShippingRateCalculator`, and `IZoneRepository`.
Implementations include compile-time interface assertions.

`IParcelService.Calculate()` returns `(*valueobjects.Money, error)`: a non-nil
result on success and nil on error. The calculator also returns `*Money`.

`IZoneRepository.FindByName()` returns `*entities.Zone`, or nil when the zone
does not exist. `ParcelService` handles nil and returns the unknown-zone error.

The service constructors accept interfaces for their individual dependencies.
The `NewServices()` assembly constructor accepts the concrete
`*repositories.Repositories` container. Preserve this existing wiring when
extending the app.

## Current rules and behavior

- Weight is in kilograms; length, width and height are in centimetres.
- All five flags on the `calculate` subcommand are required: `weight`, `length`, `width`, `height`, `zone`.
- Weight and dimensions must be finite numbers greater than zero.
- Dimensions are validated but do not affect current shipping rates.
- CLI zone input is lowercased when assigned to `Parcel.ZoneName`, so `EU`
  and `eu` resolve to the same zone. Repository names are `domestic`, `eu`,
  `international`, and repository lookup remains case-sensitive.
- Base rates: weight <= 1 kg costs EUR 5; 1 < weight <= 5 kg costs EUR 10;
  weight > 5 kg costs EUR 20.
- Zone surcharges: domestic EUR 0; eu EUR 8; international EUR 15.
- `Money` is a struct with `Amount float64` in major currency units and
  `Currency enums.Currency`. `Currency` is a string-based enum in `domain/enums` with the `EUR` constant.
  Services return `Money`; the command prints its currency and amount with two
  decimal places. Current tariffs use EUR. `Money.Add()` mutates the amount
  only when currencies match; otherwise it returns an error.
- The zone repository stores zones in memory and returns copies so callers
  cannot mutate stored surcharges.
- Errors are returned through the layers. The entry point prints errors to
  stderr and exits with status 1; successful rates go to stdout.

## Development conventions

- Keep entities in `domain/entities`, value objects in `domain/value-objects`
  (Go package `valueobjects`), and enums in `domain/enums`.
- Keep CLI parsing and output in `presentation/commands`, pricing and service
  orchestration in `application/services`, and tariff storage in
  `infrastructure/repositories`.
- Add new interfaces to separate files under `application/abstractions`, using
  the existing `I` naming prefix.
- Wire new repositories in `NewRepositories()` and new services in
  `NewServices()`.
- Keep tariff changes explicit and update README behavior descriptions when
  changing input or pricing rules.
- Name every method receiver `this`, including suite and mock receivers.
- Use `github.com/juju/errors` for error creation and `errors.Annotate()` for
  context when propagating errors between layers. Preserve the original cause.
  Use gerund-based messages ending in `failed`, such as
  `validating parcel failed`; leaf errors include the reason after a colon.
- Tests use `testify/suite` and `testify/mock`. Each suite method must follow
  `Test_{Method}_Returs{Result}_InCaseOf{Case}` (keep the `Returs` spelling).
  Test runners use `suite.Run()`. Mockery generates `*_mock.go` files next to
  their implementations, in the same package, using `.mockery.yml`. This makes
  the mocks importable by other packages and adds `testify/mock` to ordinary
  service and repository builds. Run `make mock-gen` after interface changes; do not edit
  generated mock files manually. The generator version is pinned in Makefile
  and the local template preserves `this` receiver names.
- Put each chained method call on its own line in tests. Separate setup, mock
  expectations, execution and assertions with blank lines. Label test sections
  with exactly one `//Arrange`, `//Act`, and `//Assert` per test method.
  Name table-driven case slices `testCases`, give each case a `name`, and use
  `this.Run(testCase.name, ...)` for named subtests.
- Keep tests focused on core pricing, validation, zone lookups and basic flag
  binding. Avoid output/log assertions, custom failing writers and tests of
  third-party CLI help behavior.
- Run `make build`, `make test`, and `go vet ./...` after relevant changes.

## Useful commands

```sh
gofmt -w cmd internal
make build
make mock-gen
make test
go vet ./...
go run ./cmd/cli --help
go run ./cmd/cli calculate --help
go run ./cmd/cli calculate --weight 3 --length 30 --width 20 --height 10 --zone eu
```

The example parcel costs `EUR 18.00`. For pricing changes, also manually check
weights at 1 kg and 5 kg and immediately above these bounds. For input changes,
check missing flags, an unknown zone, and non-positive or non-finite numbers.

Build output under `bin/` is ignored by Git. Stage only files relevant to the
requested change and preserve unrelated local files and edits.
