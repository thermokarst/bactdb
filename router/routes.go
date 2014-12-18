package router

const (
	User       = "users:get"
	CreateUser = "users:create"
	Users      = "users:list"
	GetToken   = "token:get"

	Genus       = "genus:get"
	CreateGenus = "genus:create"
	Genera      = "genus:list"
	UpdateGenus = "genus:update"
	DeleteGenus = "genus:delete"

	Species       = "species:get"
	CreateSpecies = "species:create"
	SpeciesList   = "species:list"
	UpdateSpecies = "species:update"
	DeleteSpecies = "species:delete"

	Strain       = "strain:get"
	CreateStrain = "strain:create"
	Strains      = "strain:list"
	UpdateStrain = "strain:update"
	DeleteStrain = "strain:delete"

	CharacteristicType       = "characteristic_type:get"
	CreateCharacteristicType = "characteristic_type:create"
	CharacteristicTypes      = "characteristic_type:list"
	UpdateCharacteristicType = "characteristic_type:update"
	DeleteCharacteristicType = "characteristic_type:delete"

	Characteristic       = "characteristic:get"
	CreateCharacteristic = "characteristic:create"
	Characteristics      = "characteristic:list"
	UpdateCharacteristic = "characteristic:update"
	DeleteCharacteristic = "characteristic:delete"

	TextMeasurementType       = "text_measurement_type:get"
	CreateTextMeasurementType = "text_measurement_type:create"
	TextMeasurementTypes      = "text_measurement_type:list"
	UpdateTextMeasurementType = "text_measurement_type:update"
	DeleteTextMeasurementType = "text_measurement_type:delete"

	UnitType       = "unit_type:get"
	CreateUnitType = "unit_type:create"
	UnitTypes      = "unit_type:list"
	UpdateUnitType = "unit_type:update"
	DeleteUnitType = "unit_type:delete"

	Measurement       = "measurements:get"
	CreateMeasurement = "measurements:create"
	Measurements      = "measurements:list"
	UpdateMeasurement = "measurements:update"
	DeleteMeasurement = "measurements:delete"
)
