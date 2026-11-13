# Changelog

Release notes for every version are published with its
[GitHub Release](https://github.com/LFDT-Iguazu/iguazu-types/releases).

## v0.4.0

### Changed

- Requires Go 1.26 (`go 1.26.0`, toolchain `go1.26.11`); modules depending on `iguazu-types` must use Go 1.26 or later.
- Dependencies updated to versions with known security fixes: `golang.org/x/crypto` v0.59.0, `golang.org/x/sys` v0.50.0.
- `github.com/supranational/blst` updated to v0.3.17.

### Removed

- Unused dependencies `github.com/stretchr/testify`, `github.com/pkg/errors` and `gotest.tools`.

## v0.3.0

### Added

- `ValidatorSync` gets the optional `StakeAmount` field, the part of `Amount` taken from the stake. A nil value means it is not set and is not the same as zero. `Copy()` deep-copies it, JSON encodes it as `stakeAmount` (omitted when nil) and `Print()` shows it. JSON payloads without the field are unchanged.

## v0.2.0

### Changed

- Reimplemented the `common` types (`Hash`, `Address`, `BlsPubKey`, `BlsSignature`, hex helpers), the `common/hexutil` package and the `crypto/bls_sig` wrapper. Behavior and wire formats are unchanged.

### Removed

- Unused APIs: `MixedcaseAddress`, `UnprefixedHash`, `UnprefixedAddress`, `UnprefixedBlsPubKey`, `UnprefixedBlsSig`, `StorageSize`, and the GraphQL and `testing/quick` helpers of the fixed-size types.

### Deprecated

- `v0.1.0` is retracted; use `v0.2.0` or later.

### Added

- Slot duration fork: `SlotInfo` gets the `ForkSlotUpSecondsPerSlot` and `UpSecondsPerSlot` fields and the `IsForkUpSecondsPerSlot`, `SecondsPerSlotBySlot`, `SlotByTimestamp` and `MatchCoordinatorParams` methods; `StartSlotTime` and `CurrentSlot` follow the fork. `iface.ChainConfig` gets the same two fields. With the fork disabled, timing and the JSON encoding of `SlotInfo` are unchanged.
- GitHub Actions: CodeQL static analysis (security and quality queries).
- GitHub Actions: CI (build, vet, test), OpenSSF Scorecard and signed releases with SLSA provenance (see [RELEASE.md](RELEASE.md)); Dependabot.
- Repository files required by LF Decentralized Trust: `RELEASE.md`, `NOTICE`, `ADOPTERS.md`, issue and pull request templates.

## v0.1.0

- Initial release of `iguazu-types` in the LF Decentralized Trust Iguazu project.
