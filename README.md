# wp_converter

Quick and dirty conversion tool from an Overlander CSV download to OsmAnd favorites GPX.

You should be familiar with go, commandline, reading code to understand possible problems, and, last but not least, OsmAnd and its handling of favorite import and export.

Use this tool at your own risk.

# Usage

Clone the repository, compile wp_converter (or use go run), then use like this:

```shell
wp_converter -i infile.csv > outfile.gpx
```

# Workflow

- Download an iOverlander csv file for a country
- Run wp_converter
- Make the outfile available on a device
- Import the gpx file as a favorites file into OsmAnd

# Supported Categories

- Established Campground
- Informal Campsite
- Wild Camping
- Water
- Mechanic and Parts
- Shopping
- Laundromat
- Fuel Station
