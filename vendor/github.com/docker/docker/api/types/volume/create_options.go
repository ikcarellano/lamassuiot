package volume

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// CreateOptions VolumeConfig
//
// Volume configuration
// swagger:model CreateOptions
type CreateOptions struct {

	// cluster volume spec
	ClusterVolumeSpec *ClusterVolumeSpec `json:"ClusterVolumeSpec,omitempty"`

	// Name of the volume driver to use.
	Driver string `json:"Driver,omitempty"`

	// A mapping of driver options and values. These options are
	// passed directly to the driver and are driver specific.
	//
	DriverOpts map[string]string `json:"DriverOpts,omitempty"`

	// User-defined key/value metadata.
	Labels map[string]string `json:"Labels,omitempty"`

	// The new volume's name. If not specified, Docker generates a name.
	//
	Name string `json:"Name,omitempty"`
}
