# go-astm ![build status](https://travis-ci.org/78bit/uuid.svg?branch=master)

Golang library for handling ASTM lis2a2 Procotol

###### Install
`go get github.com/blutspende/go-astm`

## Features
  - Encoding 
    - UTF8 
    - ASCII
    - Windows1250 
    - Windows1251 
    - Windows1252 
    - DOS852 
    - DOS855 
    - DOS866 
	- ISO8859_1
  - Timezone conversion on marshal and unmarshal
  - Marshalling and Unmarshalling supported
  - Custom delimiters are recognized in the Header and appplied (defaults are \^&)
  - Supported Types : string, float32, float64, time.Time, enums, int

## Installation

Install the package with the following command.

``` shell
go get github.com/blutspende/go-astm/...
```

## Messaage Structure

In order to encode (marshal) or decode (unmarshal) a message from or to lis, you need to annotate a struct in go to identify the record-types 
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

## Reading ASTM

The following Go code decodes a ASTM provided as a string and stores all its information in the &message.

``` go
var message standardlis2a2.DefaultMessage

err := lis2a2.Unmarshal([]byte(textdata), &message,
		lis2a2.EncodingUTF8, lis2a2.TimezoneEuropeBerlin)
if err != nil {
  log.Fatal(err)		
}
```

## Reading ASTM with multiple message in one transmission
The same code, just use DefaultMultiMessage:

``` go
  var message standardlis2a2.DefaultMultiMessage

  lis2a2.Unmarshal([]byte(textdata), &message,
		lis2a2.EncodingUTF8, lis2a2.TimezoneEuropeBerlin)		

  for _, message := range message.Messages {
	fmt.Printf("%+v", message)
  }
  
```

## Writing ASTM

Converting an annotated Structure (see above) to an enocded bytestream. 

The bytestream is encoded by-row, lacking the CR code at the end. 

``` go
lines, err := lis2a2.Marshal(msg, lis2a2.EncodingASCII, lis2a2.TimezoneEuropeBerlin, lis2a2.ShortNotation)

// output on screen
for _, line := range lines {
		linestr := string(line)
		fmt.Println(linestr)
}
```

## Identifying a message
Identifying the type of a message without decoding it. There are 3 Types of messages 
  - MessageTypeQuery 
  - MessageTypeOrdersOnly
  - MessageTypeOrdersAndResults

``` go
messageType, _ := IdentifyMessage([]byte(astm), EncodingUTF8)

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

## Message Structure and Annotation

### Optional records
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

## Custom Record Structure

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
