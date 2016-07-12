# AStrings2ObjC
The converter of Android strings resources to iOS Localisable.strings

# Overview
The handy tool to convert Android strings resources to iOS Localisable.strings using ethalon Localisable.strings to match keys. Allows easy porting translations from one platform to another.

# Usage
**$ go run main.go inputFile locFile [outputFile]**

where:
* inputFile - the path to the Android String XML file
* locFile - the path to the Localizable.strings file which has ethalon keys
* outputFile - the path to the output file (optional). If missed than name of input will be used by appending .strings suffix
