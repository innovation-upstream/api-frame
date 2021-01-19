package field

// FieldPurpose describes a fields purpose
type FieldPurpose uint32

const (
	// PurposeReferenceOwner is for indicating a field references an internal identity
	PurposeReferenceOwner FieldPurpose = 0
	// PurposeReferenceInternal is for indicating a field references internal data
	PurposeReferenceInternal FieldPurpose = 1
	// PurposeReferenceExternal is for indicating a field references external data
	PurposeReferenceExternal FieldPurpose = 2
	// PurposeReferenceNone is for indicating a field doesn't reference known data
	PurposeReferenceNone FieldPurpose = 3
)
