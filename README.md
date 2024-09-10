# wp_converter

Quick and dirty conversion tool from an Overlander CSV download to OsmAnd favorites GPX.

You should be familiar with go, commandline, reading code to understand possible problems, and, last but not least, OsmAnd and its handling of favorite import and export.

Use this tool at your own risk.

# Usage

Clone the repository, compile wp_converter (or use "go run ."), then use like this:

```shell
wp_converter -i infile.csv > outfile.gpx
```

or

```shell
go run . -i infile.csv > outfile.gpx
```

# Workflow

- Download an iOverlander csv file for a country
- Run wp_converter
- Make the outfile available on a device
- Import the gpx file as a favorites file into OsmAnd

# Supported Categories (colors and symbols, more can be added easily)

- Established Campground
- Informal Campsite
- Wild Camping
- Water
- Mechanic and Parts
- Shopping
- Laundromat
- Fuel Station

# Planned Work

- Tests ... I really need tests, I was just lazy so far
- Input validation, currently the infile is not validated at available
- At least some documentation
- Code cleanup ... I haven't written code in many years and it shows badly, sorry for the mess, hopefully it'll get better over time
- One feature idea is to be able to provide a mapping file to map categories to OsmAnd categories, icons, and colors, right now this is hardcoded

# Feedback

I'm happy to take feedback and will continue to work on this tool since the combination of iOverlander data and OsmAnd has worked well for me on a recent trip. Please file feature requests or bug reports as needed.
