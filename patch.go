package uipack

type BundleMetadataPatch struct {
	Version             uint32
	DeprecatedVariables []DeprecatedVariable
}
type AddedMode struct {
	Name string
}

type AddedModeVariant struct {
	Name string
}

type AddedVariable struct {
	Name string
}

type DeprecatedVariable struct {
	Identifier uint64
}

type BundlePatch struct {
}
