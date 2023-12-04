package rules

import (
	"fmt"

	"github.com/fbiville/markdown-table-formatter/pkg/markdown"
	"golang.org/x/sync/errgroup"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/internal/utils"
)

type (
	ServiceClass        string
	BillingModel        string
	DisponibilityClass  string
	StorageProfileClass string

	VCPUInMhz       map[BillingModel]RuleValues
	CPUAllocated    map[BillingModel]RuleValues
	StorageProfiles map[StorageProfileClass]RuleValues

	DisponibilityClasses   []DisponibilityClass
	BillingModels          []BillingModel
	ServiceClasses         []ServiceClass
	StorageBillingModels   []BillingModel
	StorageProfilesClasses []StorageProfileClass

	RuleValues struct {
		Editable bool `json:"editable"`
		Min      *int `json:"min"`
		Max      *int `json:"max"`
		Equal    *int `json:"equal"`
	}

	Rule struct {
		// System
		VCPUInMhz       VCPUInMhz    `json:"cpu_in_mhz"`
		CPUAllocated    CPUAllocated `json:"cpu_allocated"`
		MemoryAllocated RuleValues   `json:"memory_allocated"`

		// Contract
		BillingModels        BillingModels        `json:"billing_model"`
		DisponibilityClasses DisponibilityClasses `json:"disponibility_class"`
		StorageProfiles      StorageProfiles      `json:"storage_profiles"`
		StorageBillingModel  StorageBillingModels `json:"storage_billing_model"`
	}

	Rules map[ServiceClass]Rule
)

const (
	BillingModelReserved BillingModel = "RESERVED"
	BillingModelPayg     BillingModel = "PAYG"
	BillingModelDraas    BillingModel = "DRAAS"

	DisponibilityClassOneRoom    DisponibilityClass = "ONE-ROOM"
	DisponibilityClassHaDualRoom DisponibilityClass = "HA-DUAL-ROOM"
	DisponibilityClassDualRoom   DisponibilityClass = "DUAL-ROOM"

	ServiceClassEco  ServiceClass = "ECO"
	ServiceClassStd  ServiceClass = "STD"
	ServiceClassHp   ServiceClass = "HP"
	ServiceClassVoip ServiceClass = "VOIP"

	StorageProfileClassSilver       StorageProfileClass = "silver"
	StorageProfileClassSilverR1     StorageProfileClass = "silver_r1"
	StorageProfileClassSilverR2     StorageProfileClass = "silver_r2"
	StorageProfileClassGold         StorageProfileClass = "gold"
	StorageProfileClassGoldR1       StorageProfileClass = "gold_r1"
	StorageProfileClassGoldR2       StorageProfileClass = "gold_r2"
	StorageProfileClassGoldHm       StorageProfileClass = "gold_hm"
	StorageProfileClassPlatinum3k   StorageProfileClass = "platinum3k"
	StorageProfileClassPlatinum3kR1 StorageProfileClass = "platinum3k_r1"
	StorageProfileClassPlatinum3kR2 StorageProfileClass = "platinum3k_r2"
	StorageProfileClassPlatinum3kHm StorageProfileClass = "platinum3k_hm"
	StorageProfileClassPlatinum7k   StorageProfileClass = "platinum7k"
	StorageProfileClassPlatinum7kR1 StorageProfileClass = "platinum7k_r1"
	StorageProfileClassPlatinum7kR2 StorageProfileClass = "platinum7k_r2"
	StorageProfileClassPlatinum7kHm StorageProfileClass = "platinum7k_hm"

	storageProfileMinMemoryGib = 500
	storageProfileMaxMemoryGib = 50000

	memoryAllocatedMinGib = 1
	memoryAllocatedMaxGib = 5120
)

var (
	ALLBillingModels = BillingModels{
		BillingModelReserved,
		BillingModelPayg,
		BillingModelDraas,
	}

	ALLDisponibilityClasses = DisponibilityClasses{
		DisponibilityClassOneRoom,
		DisponibilityClassHaDualRoom,
		DisponibilityClassDualRoom,
	}

	ALLServiceClasses = ServiceClasses{
		ServiceClassEco,
		ServiceClassStd,
		ServiceClassHp,
		ServiceClassVoip,
	}

	ALLStorageProfilesClass = StorageProfilesClasses{
		StorageProfileClassSilver,
		StorageProfileClassSilverR1,
		StorageProfileClassSilverR2,
		StorageProfileClassGold,
		StorageProfileClassGoldR1,
		StorageProfileClassGoldR2,
		StorageProfileClassGoldHm,
		StorageProfileClassPlatinum3k,
		StorageProfileClassPlatinum3kR1,
		StorageProfileClassPlatinum3kR2,
		StorageProfileClassPlatinum3kHm,
		StorageProfileClassPlatinum7k,
		StorageProfileClassPlatinum7kR1,
		StorageProfileClassPlatinum7kR2,
		StorageProfileClassPlatinum7kHm,
	}

	ALLStorageBillingModels = StorageBillingModels{
		BillingModelPayg,
		BillingModelReserved,
	}
)

var (
	defaultStorageProfileRule = RuleValues{
		Editable: true,
		Min:      utils.ToPTR(storageProfileMinMemoryGib),
		Max:      utils.ToPTR(storageProfileMaxMemoryGib),
	}

	defaultStoragesProfiles = func() StorageProfiles {
		sp := StorageProfiles{}
		for _, c := range ALLStorageProfilesClass {
			sp[c] = defaultStorageProfileRule
		}
		return sp
	}

	defaultMemoryAllocatedRule = RuleValues{
		Editable: true,
		Min:      utils.ToPTR(memoryAllocatedMinGib),
		Max:      utils.ToPTR(memoryAllocatedMaxGib),
	}
)

var (
	ErrServiceClassNotFound = fmt.Errorf("service class not found")

	ErrBillingModelNotAvailable = fmt.Errorf("billing model is not available")

	ErrDisponibilityClassNotFound = fmt.Errorf("disponibility class not found")

	ErrStorageProfileClassNotFound = fmt.Errorf("storage profile class not found")

	ErrStorageBillingModelNotFound = fmt.Errorf("storage billing model not found")

	ErrVCPUInMhzInvalid               = fmt.Errorf("VCPU in mhz is not valid")
	ErrCPUAllocatedInvalid            = fmt.Errorf("CPU allocated is not valid")
	ErrMemoryAllocatedInvalid         = fmt.Errorf("memory allocated is not valid")
	ErrStorageProfileDefault          = fmt.Errorf("only one storage profile can be default")
	ErrStorageProfileLimitInvalid     = fmt.Errorf("storage profile limit is not valid")
	ErrStorageProfileLimitNotIntegrer = fmt.Errorf("storage profile limit is not integrer")
)

var vdcRules = Rules{
	ServiceClassEco: {
		BillingModels: []BillingModel{
			BillingModelReserved,
			BillingModelPayg,
			BillingModelDraas,
		},
		DisponibilityClasses: []DisponibilityClass{
			DisponibilityClassOneRoom,
			DisponibilityClassDualRoom,
		},
		VCPUInMhz: VCPUInMhz{
			BillingModelReserved: {
				Editable: true,
				Min:      utils.ToPTR(1200),
				Max:      utils.ToPTR(2200),
			},
			BillingModelPayg: {
				Editable: false,
				Equal:    utils.ToPTR(2200),
			},
			BillingModelDraas: {
				Editable: false,
				Equal:    utils.ToPTR(2200),
			},
		},
		CPUAllocated: CPUAllocated{
			BillingModelPayg: {
				Editable: true,
				Min:      utils.ToPTR(5 * 2200),
				Max:      utils.ToPTR(200 * 2200),
			},
			BillingModelDraas: {
				Editable: true,
				Min:      utils.ToPTR(5 * 2200),
				Max:      utils.ToPTR(200 * 2200),
			},
			BillingModelReserved: {
				Editable: true,
				Min:      utils.ToPTR(3000),
				Max:      utils.ToPTR(2500000),
			},
		},
		MemoryAllocated:     defaultMemoryAllocatedRule,
		StorageProfiles:     defaultStoragesProfiles(),
		StorageBillingModel: ALLStorageBillingModels,
	},
	ServiceClassStd: {
		BillingModels: []BillingModel{
			BillingModelReserved,
			BillingModelPayg,
			BillingModelDraas,
		},
		DisponibilityClasses: []DisponibilityClass{
			DisponibilityClassOneRoom,
			DisponibilityClassHaDualRoom,
			DisponibilityClassDualRoom,
		},
		VCPUInMhz: VCPUInMhz{
			BillingModelReserved: {
				Editable: true,
				Min:      utils.ToPTR(1200),
				Max:      utils.ToPTR(2200),
			},
			BillingModelPayg: {
				Equal: utils.ToPTR(2200),
			},
			BillingModelDraas: {
				Equal: utils.ToPTR(2200),
			},
		},
		CPUAllocated: CPUAllocated{
			BillingModelPayg: {
				Editable: true,
				Min:      utils.ToPTR(5 * 2200),
				Max:      utils.ToPTR(200 * 2200),
			},
			BillingModelDraas: {
				Editable: true,
				Min:      utils.ToPTR(5 * 2200),
				Max:      utils.ToPTR(200 * 2200),
			},
			BillingModelReserved: {
				Editable: true,
				Min:      utils.ToPTR(3000),
				Max:      utils.ToPTR(2500000),
			},
		},
		MemoryAllocated:     defaultMemoryAllocatedRule,
		StorageProfiles:     defaultStoragesProfiles(),
		StorageBillingModel: ALLStorageBillingModels,
	},
	ServiceClassHp: {
		BillingModels: []BillingModel{
			BillingModelReserved,
			BillingModelPayg,
		},
		DisponibilityClasses: []DisponibilityClass{
			DisponibilityClassOneRoom,
			DisponibilityClassHaDualRoom,
			DisponibilityClassDualRoom,
		},
		VCPUInMhz: VCPUInMhz{
			BillingModelReserved: {
				Equal: utils.ToPTR(2200),
			},
			BillingModelPayg: {
				Equal: utils.ToPTR(2200),
			},
		},
		CPUAllocated: CPUAllocated{
			BillingModelPayg: {
				Editable: true,
				Min:      utils.ToPTR(5 * 2200),
				Max:      utils.ToPTR(200 * 2200),
			},
			BillingModelReserved: {
				Editable: true,
				Min:      utils.ToPTR(3000),
				Max:      utils.ToPTR(2500000),
			},
		},
		MemoryAllocated:     defaultMemoryAllocatedRule,
		StorageProfiles:     defaultStoragesProfiles(),
		StorageBillingModel: ALLStorageBillingModels,
	},
	ServiceClassVoip: {
		BillingModels: []BillingModel{
			BillingModelReserved,
		},
		DisponibilityClasses: []DisponibilityClass{
			DisponibilityClassOneRoom,
			DisponibilityClassHaDualRoom,
			DisponibilityClassDualRoom,
		},
		VCPUInMhz: VCPUInMhz{
			BillingModelReserved: {
				Equal: utils.ToPTR(3000),
			},
		},
		CPUAllocated: CPUAllocated{
			BillingModelReserved: {
				Editable: true,
				Min:      utils.ToPTR(3000),
				Max:      utils.ToPTR(2500000),
			},
		},
		MemoryAllocated:     defaultMemoryAllocatedRule,
		StorageProfiles:     defaultStoragesProfiles(),
		StorageBillingModel: ALLStorageBillingModels,
	},
}

// String returns the string representation of the DisponibilityClass.
func (dc DisponibilityClasses) String() string {
	s := ""
	for _, c := range dc {
		s += string(c) + ", "
	}
	return s[:len(s)-2]
}

// String returns the string representation of the BillingModel.
func (bm BillingModels) String() string {
	s := ""
	for _, m := range bm {
		s += string(m) + ", "
	}
	return s[:len(s)-2]
}

// String returns the string representation of the StorageBillingModel.
func (sbm StorageBillingModels) String() string {
	s := ""
	for _, m := range sbm {
		s += string(m) + ", "
	}
	return s[:len(s)-2]
}

// String returns the string representation of the StorageProfileClass.
func (sp StorageProfiles) String() string {
	s := ""
	for _, c := range ALLStorageProfilesClass {
		s += fmt.Sprintf("%s: %s, ", c, sp[c])
	}
	return s[:len(s)-2]
}

// String returns the string representation of the ServiceClass.
func (sc ServiceClasses) String() string {
	s := ""
	for _, c := range sc {
		s += string(c) + ", "
	}

	return s[:len(s)-2]
}

// String returns the string representation of the RuleValues.
func (rv RuleValues) String() string {
	s := ""
	if rv.Editable {
		s += "** "
	}
	if rv.Min != nil {
		s += fmt.Sprintf("min: %d, ", *rv.Min)
	}
	if rv.Max != nil {
		s += fmt.Sprintf("max: %d, ", *rv.Max)
	}
	if rv.Equal != nil {
		s += fmt.Sprintf("equal: %d, ", *rv.Equal)
	}
	return s[:len(s)-2]
}

// ParseServiceClass returns the ServiceClass from the given string.
func ParseServiceClass(s string) (ServiceClass, error) {
	for _, sc := range ALLServiceClasses {
		if sc == ServiceClass(s) {
			return sc, nil
		}
	}
	return "", fmt.Errorf("%w (Allowed values: %v)", ErrServiceClassNotFound, ALLServiceClasses)
}

// ParseStorageBillingModel returns the StorageBillingModel from the given string.
func ParseStorageBillingModel(s string) (BillingModel, error) {
	for _, bm := range ALLStorageBillingModels {
		if bm == BillingModel(s) {
			return bm, nil
		}
	}
	return "", fmt.Errorf("%w (Allowed values: %v)", ErrStorageBillingModelNotFound, ALLStorageBillingModels)
}

// GetRuleByServiceClass returns the Rule for the given ServiceClass.
func GetRuleByServiceClass(sc ServiceClass) (Rule, error) {
	r, ok := vdcRules[sc]
	if !ok {
		return Rule{}, ErrServiceClassNotFound
	}
	return r, nil
}

// billingModelIsValid returns true if the given BillingModel is valid for the given ServiceClass.
func (r Rule) billingModelIsValid(bm BillingModel) bool {
	for _, m := range r.BillingModels {
		if m == bm {
			return true
		}
	}
	return false
}

// disponibilityClassIsValid returns true if the given DisponibilityClass is valid for the given ServiceClass.
func (r Rule) disponibilityClassIsValid(dc DisponibilityClass) bool {
	for _, c := range r.DisponibilityClasses {
		if c == dc {
			return true
		}
	}
	return false
}

// storageProfileClassIsValid returns true if the given StorageProfileClass is valid for the given ServiceClass.
func (r Rule) storageProfileClassIsValid(sp StorageProfileClass) bool {
	for _, c := range ALLStorageProfilesClass {
		if c == sp {
			return true
		}
	}
	return false
}

// storageBillingModelIsValid returns true if the given BillingModel is valid for the given ServiceClass.
func (r Rule) storageBillingModelIsValid(bm BillingModel) bool {
	for _, m := range r.StorageBillingModel {
		if m == bm {
			return true
		}
	}
	return false
}

// isValid RuleValues returns true if the given RuleValues is valid for the given ServiceClass.
func (rv RuleValues) isValid(value int) bool {
	if rv.Min != nil && value < *rv.Min {
		return false
	}

	if rv.Max != nil && value > *rv.Max {
		return false
	}

	if rv.Equal != nil && value != *rv.Equal {
		return false
	}

	return true
}

// vCPUInMhzIsValid returns true if the given vCPUInMhz is valid for the given ServiceClass.
func (r Rule) vCPUInMhzIsValid(bm BillingModel, vcpuInMhz int) bool {
	m, ok := r.VCPUInMhz[bm]
	if !ok {
		return false
	}

	return m.isValid(vcpuInMhz)
}

// cpuAllocatedIsValid returns true if the given CPUAllocated is valid for the given ServiceClass.
func (r Rule) cpuAllocatedIsValid(bm BillingModel, cpuAllocated int) bool {
	m, ok := r.CPUAllocated[bm]
	if !ok {
		return false
	}

	return m.isValid(cpuAllocated)
}

// memoryAllocatedIsValid returns true if the given MemoryAllocated is valid for the given ServiceClass.
func (r Rule) memoryAllocatedIsValid(memoryAllocated int) bool {
	return r.MemoryAllocated.isValid(memoryAllocated)
}

type ValidateData struct {
	ServiceClass       ServiceClass
	BillingModel       BillingModel
	DisponibilityClass DisponibilityClass
	StorageProfiles    map[StorageProfileClass]struct {
		Limit   int
		Default bool
	}
	StorageBillingModel BillingModel
	VCPUInMhz           int
	CPUAllocated        int
	MemoryAllocated     int
}

// Validate returns true if the given BillingModel, DisponibilityClass and StorageProfileClass are valid for the given ServiceClass.
func Validate(data ValidateData, isUpdate bool) error {
	r, err := GetRuleByServiceClass(data.ServiceClass)
	if err != nil {
		return err
	}

	// TODO check if value is editable

	var wg errgroup.Group

	// goroutine to validate the BillingModel
	wg.Go(func() error {
		if !r.billingModelIsValid(data.BillingModel) {
			return fmt.Errorf("if service class is %s the %w: %s (Allowed values: %v)", data.ServiceClass, ErrBillingModelNotAvailable, data.BillingModel, r.BillingModels)
		}
		return nil
	})

	// goroutine to validate the VCPUInMhz
	wg.Go(func() error {
		if !r.vCPUInMhzIsValid(data.BillingModel, data.VCPUInMhz) {
			return fmt.Errorf("if service class is %s and the billing model is %s the value of %w: %d (Allowed values: %v)", data.ServiceClass, data.BillingModel, ErrVCPUInMhzInvalid, data.VCPUInMhz, r.VCPUInMhz[data.BillingModel])
		}
		return nil
	})

	// goroutine to validate the CPUAllocated
	wg.Go(func() error {
		if !r.cpuAllocatedIsValid(data.BillingModel, data.CPUAllocated) {
			return fmt.Errorf("if service class is %s and the billing model is %s the value of %w: %d (Allowed values: %v)", data.ServiceClass, data.BillingModel, ErrCPUAllocatedInvalid, data.CPUAllocated, r.CPUAllocated[data.BillingModel])
		}
		return nil
	})

	// goroutine to validate the MemoryAllocated
	wg.Go(func() error {
		if !r.memoryAllocatedIsValid(data.MemoryAllocated) {
			return fmt.Errorf("%w: %d (Allowed values: %v)", ErrMemoryAllocatedInvalid, 0, r.MemoryAllocated)
		}
		return nil
	})

	// goroutine to validate the StorageBillingModel
	wg.Go(func() error {
		if !r.storageBillingModelIsValid(data.StorageBillingModel) {
			return fmt.Errorf("%w: %s (Allowed values: %v)", ErrStorageBillingModelNotFound, data.BillingModel, r.StorageBillingModel)
		}
		return nil
	})

	// goroutine to validate the DisponibilityClass
	wg.Go(func() error {
		if !r.disponibilityClassIsValid(data.DisponibilityClass) {
			return fmt.Errorf("%w: %s (Allowed values: %v)", ErrDisponibilityClassNotFound, data.DisponibilityClass, r.DisponibilityClasses)
		}
		return nil
	})

	// goroutine to validate the StorageProfileClass
	wg.Go(func() error {
		defaultStorageProfile := 0
		for c, sP := range data.StorageProfiles {
			if !r.storageProfileClassIsValid(c) {
				return fmt.Errorf("%w: %s (Allowed values: %v)", ErrStorageProfileClassNotFound, c, ALLStorageProfilesClass)
			}
			if sP.Limit < *r.StorageProfiles[c].Min || sP.Limit > *r.StorageProfiles[c].Max {
				return fmt.Errorf("%w: %d (Allowed values: %v)", ErrStorageProfileLimitInvalid, sP.Limit, r.StorageProfiles[c])
			}
			// // Limit is valid if modulo of 1024 is 0
			// if sP.Limit%1024 != 0 {
			// 	return fmt.Errorf("%w: %d. Value must be a multiple of 1024", ErrStorageProfileLimitNotIntegrer, sP.Limit)
			// }

			if sP.Default {
				defaultStorageProfile += 1
			}
		}

		if defaultStorageProfile > 1 {
			return ErrStorageProfileDefault
		}

		return nil
	})

	return wg.Wait()
}

// GetRulesDetails returns the RuleValues for the given BillingModel and DisponibilityClass.
// Return markdown table with the following columns:
// - BillingModel
// - StorageBillingModels
// - DisponibilityClass
// - VCPUInMhz
// - CPUAllocated
func (r Rule) GetRuleDetails() string {
	rules := [][]string{}

	for _, bm := range r.BillingModels {
		rules = append(rules, []string{
			string(bm), r.StorageBillingModel.String(), r.DisponibilityClasses.String(), r.VCPUInMhz[bm].String(), r.CPUAllocated[bm].String(), r.MemoryAllocated.String(),
		})
	}

	prettyPrintedTable, err := markdown.NewTableFormatterBuilder().
		WithPrettyPrint().
		Build("BillingModels", "StorageBillingModels", "DisponibilityClasses", "VCPUInMhz", "CPUAllocated", "MemoryAllocated").
		Format(rules)
	if err != nil {
		panic(err)
	}

	return prettyPrintedTable
}

// GetRuleDetails returns one markdown table for each ServiceClass.
func GetRulesDetails() string {
	x := ""
	x += "# Rules\n"
	x += "All fields with a ** are editable.\n\n"
	for _, sc := range ALLServiceClasses {
		r, err := GetRuleByServiceClass(sc)
		if err != nil {
			panic(err)
		}
		x += fmt.Sprintf("## ServiceClass %s\n", sc)
		x += r.GetRuleDetails()
		x += "\n"
	}
	return x
}
