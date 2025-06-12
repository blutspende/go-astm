# go-astm
Library for handling the ASTM protocol in Go.
###### Install
`go get github.com/blutspende/go-astm/v3`

# Features
  - Marshalling and unmarshalling of ASTM messages
  - Baseline message structures for LIS02-A2
  - Custom delimiters automatically identified (defaults are \^&)
  - Line breaks automatically identified (default is \n)
  - Encoding from-to raw bytes and automatic timezone conversions are included using blooblab-common

3 main functions and a utility is provided:
- `Marshal`: Converts a Go structure to an array of byte arrays
- `Unmarshal`: Converts a byte array to a Go structure
- `IdentifyMessage`: Identifies the type of message without decoding it
- `NewDefaultConfiguration`: Returns a copy of the default configuration
``` go
func Marshal(sourceStruct interface{}, configuration ...models.Configuration) (result [][]byte, err error) 
func Unmarshal(messageData []byte, targetStruct interface{}, configuration ...models.Configuration) (err error)
func IdentifyMessage(messageData []byte, configuration ...models.Configuration) (messageType messagetype.MessageType, err error)
func NewDefaultConfiguration() astmmodels.Configuration
```

# Setting up configuration
For all three functions a configuration structure can be provided to determine behaviour.

``` go
type Configuration struct {
	Encoding                   encoding.Encoding
	LineSeparator              string
	AutoDetectLineSeparator    bool
	TimeZone                   timezone.TimeZone
	EnforceSequenceNumberCheck bool
	Notation                   string
	DefaultDecimalPrecision    int
	RoundLastDecimal           bool
	KeepShortDateTimeZone      bool
	EscapeOutputStrings        bool
	Delimiters                 Delimiters
	TimeLocation               *time.Location
}
```
It can also be omitted, in case the default is used:
``` go
var DefaultConfiguration = Configuration{
	Encoding:                   encoding.ISO8859_1,
	LineSeparator:              lineseparator.LF,
	AutoDetectLineSeparator:    true,
	TimeZone:                   timezone.EuropeBerlin,
	EnforceSequenceNumberCheck: true,
	Notation:                   notation.Standard,
	DefaultDecimalPrecision:    3,
	RoundLastDecimal:           true,
	KeepShortDateTimeZone:      true,
	EscapeOutputStrings:        false,
	Delimiters:                 DefaultDelimiters,
	TimeLocation:               nil,
}
var DefaultDelimiters = Delimiters{
	Field:     `|`,
	Repeat:    `\`,
	Component: `^`,
	Escape:    `&`,
}
```
## Encoding
Character encoding for reading and writing bytes. Options are all enum constants from `github.com/blutspende/bloodlab-common/encoding`.
## LineSeparator
Line separator can be auto-detected, or set manually. If `AutoDetectLineSeparator` is set to true, this can be ignored. A few constants are provided for convenience, but any string is valid. This is only relevant for unmarshal.
``` go
lineseparator.LF
lineseparator.CR
lineseparator.LFCR
lineseparator.CRLF
```
## AutoDetectLineSeparator
If set to true, the line separator is detected automatically. If set to false, the line separator set in `LineSeparator` is used. This is only relevant for unmarshal.
## TimeZone
The timezone is used for date/time conversion. Options are all enum constants from `github.com/blutspende/bloodlab-common/timezone`.
## EnforceSequenceNumberCheck
In unmarshal, the sequence number (second field in every line) is checked for validity. If set to true, an error is returned if the sequence number is incorrect. If set to false, it is ignored. This is only relevant for unmarshal.
## Notation
The notation is only used marshal. The notation is set to one of the following:
``` go
notation.Standard
notation.Short
```
Standard notation will produce as many fields as there are in the source structure, while short notation will omit empty fields at the end of a line. This is only relevant for marshal.
## DefaultDecimalPrecision
The default decimal precision is used for floating point numbers. If a field is not annotated with `length:N`, the default decimal precision is used. `-1` can be set to allow numbers to take any length required. This is only relevant for marshal.
## RoundLastDecimal
If it is set to true, floating point numbers are rounded up or down (based on common rounding rules) at the last decimal place determined by either `DefaultDecimalPrecision` or `length:N` annotation. If it is set to false the excess decimals are truncated. This is only relevant for marshal.
## KeepShortDateTimeZone
As short dates are only year, month, day, representing the date as midnight of that day in the `time.Time`, time zone conversions can lead to the change of day (e.g. 23:00 the day before). Because logically the short date represents a whole day, it can be more important to preserve the actual date as is then to have the time in UTC. However `time.Time` (and most database representations) must have a time zone, so the only solution is to keep the original time zone unconverted.
If this flag is set to true, the timezone is kept in local time for the short date format. If set to false, the time is converted to UTC just like long dates. This applies both for marshal and unmarshal, so with the same configuration the string format of the date will be intact.
## EscapeOutputStrings
If set to true, the output strings are escaped according to the delimiters. Meaning that an escape character is put before each occurrence of the delimiters (including the escape character itself). If set to false, the output strings are not escaped, and will be output directly even if they contain delimiters. Default is false. This is only relevant for marshal.
## Delimiters
Used for building the protocol's record structure. When the configuration is provided for marshal the default is automatically used if any of the delimiter's fields are empty. If all fields are set, the default can be overridden. Each field should contain exactly one character. Unmarshal automatically detects the delimiters in the header record. This is only relevant for marshal.
``` go
type Delimiters struct {
	Field     string
	Repeat    string
	Component string
	Escape    string
}
```
## TimeLocation
For internal use only. Should be ignored.

# Usage of the library functions

## Default configuration: NewDefaultConfiguration
The `NewDefaultConfiguration` function returns a copy of the default configuration. This can be used to then modify the configuration for specific use cases while leaving the rest as default. This is the safe way to get a configuration instance, as it removes the need to check changes of the configuration structure in projects that use go-astm after updating its version. Directly creating a new configuration instance can lead to unexpected behaviour if not all the values are set, especially if new values are added in future versions. The default aims to keep behaviour backwards compatible in case new functionalities are introduced.

## Identifying a message: IdentifyMessage
Identifies the type of message without decoding it. Return values are options from `github.com/blutspende/bloodlab-common/messagetype` enum definitions. Currently, there are 3 valid types and unidentified as possible return values.
It can be used for example as follows:
``` go
messageType, err := astm.IdentifyMessage([]byte(astm), config)
if err != nil {
    log.Fatal(err)
}
switch (messageType) {
	case messagetype.Unidentified:
	  ...
	case messagetype.Query:
	  ...
	case messagetype.Order:
	  ...
	case messagetype.Result:
	  ...
}
```

## Reading an ASTM message: Unmarshal
The following Go code decodes an ASTM message provided as a string and stores all its information in the message structure.
``` go
var message lis02a2.ResultMessage
err := astm.Unmarshal([]byte(textdata), &message, config)
if err != nil {
  log.Fatal(err)		
}
```
The Unmarshal can also be used for multiple messages, providing a multi-message structure like `ResultMultiMessage`:
``` go
  var message lis02a2.ResultMultiMessage
  astm.Unmarshal([]byte(textdata), &message, config)
  for _, message := range message.ResultMessages {
	fmt.Printf("%+v", message)
  }
```

## Writing an ASTM message: Marshal
Marshal converts an annotated structure to an encoded array of byte arrays. Each element represents a line of the message, and thus has no line break at the end.
``` go
lines, err := astm.Marshal(message, config)
if err != nil {
  log.Fatal(err)		
}
for _, line := range lines {
    fmt.Println(string(line))
}
```

# Annotated structures
In order to read or write an ASTM message, an annotated structure is required. The library uses the `astm` tag to identify the fields and their location in the message, as well as additional attributes.

There are two separate types of structures:
- Record structure: This is the structure that represents a single line in the message.
- Message structure: This is the structure that represents multiple lines, or the entire message.

For both of these cases, predefined structures are provided in the `lis02a2` package. These structures are based on the ASTM standard and can be used as is, or as building blocks for custom implementation.

### Attributes
After the mandatory part of the annotation a number of optional attributes can be added, separated by commas. The attributes are used to modify the behaviour of the field or record. Some attributes can also have values, which are separated from the attribute name by a colon.
``` go
FieldOrRecord Type `astm:"?,attribute1,attribute2:value"`
```

## Record structure
Example:
``` go
type Record struct {
    Field1 string `astm:"3"`
    Field2 string `astm:"4"`
    Field3 string `astm:"5"`
}
```
In this case, the 3, 4, 5 signifies the position of the field in the message.
The first two fields are reserved for the record name and the sequence number and can not be used in an annotated structure.

### Record field attributes
Additionally to the field position, there are a few attributes that can be used to modify the behaviour of the field:
``` go
type Record struct {
    Field1 string    `astm:"3,required"`
    Field2 float32   `astm:"4,length:3"`
    Field3 time.Time `astm:"5,longdate"`
}
```
- `required`: By default fields can be empty for unmarshal. However, a required field will produce an error if missing.
- `length:N`: This field is a fixed point number with N decimals. N has to be an integer >= -1. Excess decimals are either truncated or rounded during marshal.
- `longdate`: By default dates are converted in short format `YYYYMMDD` in marshal, but with this attribute it can be set to long format: `YYYYMMDDHHMMSS`.
These attributes can also be used in combination, listing them comma separated:
``` go
type Record struct {
    Field2 float32 `astm:"3,required,length:3"`
}
```

### Record field arrays
If a field is defined with an array type, it will be marshalled and unmarshalled as an array of repetitions within the field, with the repetition delimiter.
``` go
type Record struct {
    Field1 []string `astm:"3"`
}
```
```
R|1|value1\value2\value3
```

### Record field components
A field can have multiple components, separated by component delimiter. These have to be defined by separate variables in the structure, with the proper annotation:
``` go
type Record struct {
    Field1Component1 string `astm:"3.1"`
    Field1Component2 string `astm:"3.2"`
    Field1Component3 string `astm:"3.3"`
}
```
```
R|1|component1^component2^component3
```

### Record field substructures
A field can contain a substructure, which is defined by a separate structure with proper annotation. In this case the substructure's variables will behave like components in the field.
``` go
type Substructure struct {
    Component1 string `astm:"1"`
    Component2 string `astm:"2"`
}
type Record struct {
    Field1 Substructure `astm:"3"`
}
```
```
R|1|component1^component2
```
This allows a componented field to be reused in multiple places (with different field position) with a reliable structure and an easier use in the code built on it.
A substructure can also be an array, in which case it will be marshalled and unmarshalled as an array of repetitions of the components:
``` go
type Record struct {
    Field1 []Substructure `astm:"3"`
}
```
```
R|1|comp1^comp2\comp1^comp2\comp1^comp2
```

### Pointers in record fields
Usually fields are direct values, however, this does not allow for numeric values to be empty, and will default to 0 in marshal. Pointer values allow nil to be used, which will produce an actual empty field as an output.
``` go
type Record struct {
    Field1 *int `astm:"3"`
    Field2 int  `astm:"4"`
}
```
```
R|1||0
```

### Enums in record fields
Enums are just strings with limited value sets, represented by a redefined string type. They are also supported.
``` go
type EnumType string
type Record struct {
    Field1 EnumType `astm:"3"`
}
```

## Message structure
Examples:
``` go
type Message struct {
    Record1 RecordType `astm:"R"`
}
```
``` go
type Lis02a2Message {
    MessageHeader lis02a2.Header     `astm:"H"`
    Record        RecordType         `astm:"R"`
    Terminator    lis02a2.Terminator `astm:"L"`
}
```
The letters (R, H, etc.) in the annotation mark the record line's record name (its first field), which has to match.
```
R|1|value1|value2
```
```
H|\^&||||
R|1|value1|value2
L|1|N
```

### Message structure attributes
Additionally to the record name, there is one attribute that can be used in message structures:
``` go
type Lis02a2Message {
    MessageHeader lis02a2.Header     `astm:"H"`
    Record        RecordType         `astm:"R,optional"`
    Record        RecordType         `astm:"R,subname:N"`
    Terminator    lis02a2.Terminator `astm:"L"`
}
```
- `optional`: By default all records are required for unmarshal, and will produce an error if missing. However, an optional record will just be skipped. This can also be used for composite and array structures (see below).
- `subname:N`: Some records (eg: 'M' manufacturer info) can be more free form record types. Therefore, it might have different structures with the same record name. To distinguish between them, the subname attribute can be used. Its string value will be matched with the 3rd field (the one after sequence number) to determine the type. This will be used to identify the record type in the message, or to signify the end of an array of that given type. This is only relevant for unmarshal.
These attributes can also be used in combination, listing them comma separated:
``` go
type Lis02a2Message {
    Record RecordType `astm:"R,subname:SUBNAME,optional"`
}
```

### Message structure arrays
Records can be defined as arrays, in which case they will be marshalled and unmarshalled as an array of repetitions of the records, each of them as a line in the message, with incrementing sequence numbers.
``` go
type Lis02a2Message {
    MessageHeader lis02a2.Header     `astm:"H"`
    Record        []RecordType       `astm:"R"`
    Terminator    lis02a2.Terminator `astm:"L"`
}
```
```
H|\^&||||
R|1|value1|value2
R|2|value1|value2
R|3|value1|value2
L|1|N
```

### Composite message structures
Message structures can also contain other nested message structures, in multiple layers of depth. This is defined by omitting the annotation of the nested structure.
``` go
type Nested {
    Record1 FirstRecordType  `astm:"F"`
    Record2 SecontRecordType `astm:"S"`
}
type Lis02a2Message {
    MessageHeader   lis02a2.Header     `astm:"H"`
    NestedStructure Nested
    Terminator      lis02a2.Terminator `astm:"L"`
}
```
```
H|\^&||||
F|1|value1|value2
S|1|value1|value2
L|1|N
```

### Composite array message structures
The nested structures can also be arrays, in which case it will be unmarshalled to an array as long as structure is present in the input, and will be marshalled from the array with the repetition in the source array.
``` go
type Lis02a2Message {
    MessageHeader   lis02a2.Header     `astm:"H"`
    NestedStructure []Nested
    Terminator      lis02a2.Terminator `astm:"L"`
}
```
```
H|\^&||||
F|1|value1|value2
S|1|value1|value2
F|2|value1|value2
S|1|value1|value2
F|3|value1|value2
S|1|value1|value2
L|1|N
```
Note that the sequence number is incremented for each instance of the nested structure, however only the first record of the nested structure takes the sequence number, and the rest is 1 (unless the nested structure has its own array inside).