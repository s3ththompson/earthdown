![Example Image](northern-territory-australia-2260.jpg)

## Earthdown 

Simple CLI to download images from [Google Earth View](https://earthview.withgoogle.com/).

### Installation

```
$ go get -u github.com/s3ththompson/earthdown
```

### Usage

```
earthdown [options...] EARTH_VIEW_URL
Options:
	-o name of output file
```

### Example

```
$ earthdown https://g.co/ev/2260
Northern Territory, Australia – Earth View from Google
Lat: -11.992309, Lng: 131.807527, ©2014 CNES / Astrium, Cnes/Spot Image, DigitalGlobe, Landsat, Sinclair Knight Merz, Sinclair Knight Merz & Fugro
Downloaded https://www.gstatic.com/prettyearth/assets/full/2260.jpg to northern-territory-australia-2260.jpg (1 file, 261.3 kB)
```