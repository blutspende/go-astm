# go-astm

Library for handling the ASTM LIS2-A2 protocol in go.

###### Install
`go get github.com/blutspende/go-astm/v2`

## Features
  - Encoding UTF8, ASCII, Windows1250, Windows1251, Windows1252, DOS852, DOS855, DOS866, ISO8859_1
  - Timezone conversion
  - Custom delimiters automatically identified (defaults are \^&)

## Reading ASTM

The following Go code decodes a ASTM provided as a string and stores all its information in the message.

``` go
var message standardlis2a2.DefaultMessage

err := astm.Unmarshal([]byte(textdata), &message,
		astm.EncodingUTF8, astm.TimezoneEuropeBerlin)
if err != nil {
  log.Fatal(err)		
}
```

## Reading ASTM with multiple messages in one transmission
The same code, just use DefaultMultiMessage:

``` go
  var message standardlis2a2.DefaultMultiMessage

  astm.Unmarshal([]byte(textdata), &message,
		astm.EncodingUTF8, astm.TimezoneEuropeBerlin)		

  for _, message := range message.Messages {
	fmt.Printf("%+v", message)
  }
  
```

## Writing ASTM

Converting an annotated Structure (see above) to an encoded bytestream.

The bytestream is encoded by-row, lacking the CR code at the end. 

``` go
lines, err := astm.Marshal(msg, astm.EncodingASCII, astm.TimezoneEuropeBerlin, astm.ShortNotation)

// output on screen
for _, line := range lines {
		linestr := string(line)
		fmt.Println(linestr)
}
```

## Identifying a message
Identifying the type of a message without decoding it. There are 3 types of messages 
  - MessageTypeQuery 
  - MessageTypeOrdersOnly
  - MessageTypeOrdersAndResults

``` go
messageType, _ := astm.IdentifyMessage([]byte(astm), EncodingUTF8)

switch (messageType) {
	case MessageTypeUnkown :
	  ...
	case MessageTypeQuery :
	  ...
	case MessageTypeOrdersOnly :
	  ...
	case MessageTypeOrdersAndResults :
	  ...
}
```

# How the go-astm library works
In order to encode (marshal) or decode (unmarshal) a message from or to LIS, you need to annotate a struct in go to identify the record-types 
and within the record the field's location.

The Message does now  the information of what type of message is mapped by annotation. 

``` golang
type Message struct {
	Header struct {
		field1 string `astm:"1"`
		field2 string `astm:"2"`
	} `astm:"H"` // identify the Record-Type
	PatientOrder[] struct {
		Patient struct {
			field1 string `astm:"1"`
			field2 string `astm:"2"`
			...
		} `astm:"P"`
		Order struct {
			...
		} `astm:"O"
	} 
}
```

The lis2a2-default implementation provided with this library as a starting point, it should fit most instruments. Alter it as required. 

``` go
type CommentedResult struct {
	Result  Result    `astm:"R"`
	Comment []Comment `astm:"C,optional"`
}

type PORC struct {
	Patient         Patient   `astm:"P"`
	Comment         []Comment `astm:"C,optional"`
	Order           Order     `astm:"O"`
	CommentedResult []CommentedResult
}

type DefaultMessage struct {
	Header       Header       `astm:"H"`
	Manufacturer Manufacturer `astm:"M,optional"`
	OrderResults []PORC
	Terminator   Terminator `astm:"L"`
}
```

### Message Structure and Annotation

``` go
type SimpleMessage struct  {
	Header       standardlis2a2.Header       `astm:"H"`
	Manufacturer standardlis2a2.Manufacturer `astm:"M,optional"`
	Patient      standardlis2a2.Patient      `astm:"P"`
	Order        standardlis2a2.Order        `astm:"O"`
	Result       standardlis2a2.Result       `astm:"R"`
	Terminator   standardlis2a2.Terminator   `astm:"L"`
}
```

### Nested arrays
``` go
type MessagePORC struct {
	Header       standardlis2a2.Header       `astm:"H"`
	Manufacturer standardlis2a2.Manufacturer `astm:"M,optional"`
	OrderResults []struct {
		Patient         standardlis2a2.Patient `astm:"P"`
		Order           standardlis2a2.Order   `astm:"O"`
		CommentedResult []struct {
			Result  standardlis2a2.Result    `astm:"R"`
			Comment []standardlis2a2.Comment `astm:"C,optional"`
		}
	}
	Terminator standardlis2a2.Terminator `astm:"L"`
}
```

### Addressing fields 
Often the default is not enough. You can customize any record with annotation. 

#### ... by Field#
``` go
   ...
   Filed string `astm:"3"`  // Select 3rd field, start counting with 1
   ...
```
Example:
``` text
	X|field2|field3|field4|		Result: "field3"
	X|field2^1^2|field3^1^2|field4^5^6|		Result: "field3"	
	X|field2^1^2|field3_1^1_1^2_!\\field3_2^5_2^2_2|field4^6^2|		Result: "field3_1"
```

#### ... by Field#.Component#
``` go
   ...
   Filed string `astm:"3.2"`  // Select 3rd field, 2nd component, start counting with 1
   ...
```
Example:
``` text
	X|field2|field3|field4|		Result: ""	
	X|field2^1^2|field3^1^2|field4^5^6|		Result: "1"	
	X|field2^1^2|field3_1^1_1^2_!\\field3_2^1_2^2_2|field4^1^2|		Result: "1_1"
```
#### ... by Field#.Repeat#.Component#
``` go
   ...
   Filed string `astm:"3.2.2"`  // Select 3rd field, 2nd array index, 2nd component, start counting with 1
   ...
```
Example:
``` text
	X|field2|field3|field4|		Result: ""	
	X|field2^1^2|field3^1^2|field4^5^6|		Result: ""	
	X|field2^1^2|field3_1^1_1^2_!\\field3_2^1_2^2_2|field4^1^2|		Result: "1_2"
```
### Custom Record Format
``` go
type Result struct {
	SequenceNumber                           int       `astm:"2,sequence"`   // sequence generates numbers when value is 0 
	UniversalTestID                          string    `astm:"3.1"`         
	UniversalTestIDName                      string    `astm:"3.2"`         
	UniversalTestIDType                      string    `astm:"3.3"`         
	ManufacturersTestType                    string    `astm:"3.4"`         
	ManufacturersTestName                    string    `astm:"3.5"`         
	ManufacturersTestCode                    string    `astm:"3.6"`         
	TestCode                                 string    `astm:"3.7"`         
	DataMeasurementValue                     string    `astm:"4.1"`         
	InitialMeasurementValue                  string    `astm:"4.2"`         
	MeasurementValueOfDevice                 string    `astm:"4.3"`         
	Units                                    string    `astm:"5"`           
	ReferenceRange                           string    `astm:"6"`           
	ResultAbnormalFlag                       string    `astm:"7"`           
	NatureOfAbnormalTesting                  string    `astm:"8"`           
	ResultStatus                             string    `astm:"9"`           
	DateOfChangeInInstrumentNormativeTesting time.Time `astm:"10,longdate"` 
	OperatorIDPerformed                      string    `astm:"11.1"`        
	OperatorIDVerified                       string    `astm:"11.2"`        
	DateTimeTestStarted                      time.Time `astm:"12,longdate"` 
	DateTimeCompleted                        time.Time `astm:"13,longdate"` 
	IntstrumentIdentification                string    `astm:"14"`          
}
```

### Custom decimal length for marshal
The default decimal representation is 3 digits long, which can be overridden with the `length` annotation.
This annotation is ignored if the field is not of type `float32` or `float64`.
``` go
type CustomDecimal struct {
    DefaultFormat float32 `astm:"1"`
    CustomFormat  float64 `astm:"2,length:6"`
}

type Record struct {
    CustomDecimalRecord CustomDecimal `astm:"R"`
}

message := Record{
    CustomDecimalRecord: CustomDecimal {
        DefaultFormat: 1.222222222,
        CustomFormat:  2.111111111,
    }
}

astm.Marshal(message, astm.EncodingASCII, astm.TimezoneEuropeBerlin, astm.ShortNotation)
output:
    R|1.222|2.111111
```
