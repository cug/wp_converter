# wp_converter

_iOVerlander CSV => OsmAnd GPX_

Quick and dirty waypoint conversion tool from an Overlander CSV download to OsmAnd favorites GPX.

You should be familiar with go, commandline, reading code to understand possible problems, and, last but not least, OsmAnd and its handling of favorite import and export.

Use this tool at your own risk. I'm using this tool also as my personl Go learning playground, so there might be weirdness in language constructs, misuse of language ideas and so on. File an issue if something actually bothers you.

# Usage

Clone the repository, compile wp_converter (or use "go run ."), then use like this:

```shell
wp_converter -i infile.csv -o outfile.gpx
```

or

```shell
go run . -i infile.csv -o outfile.gpx
```

or limit to a coordinate boundary box (simple rectangle), this example provides the min/max values for a boundary box that does _not_ limit the points:

```shell
go run .-i infile.csv -o outfile.gpx --lonMin=-180.0 --lonMax=180.0 --latMin=-90.0 --latMax=90.0
```

If you need to limit to a smaller area, check a maps app (e.g. Google Maps)

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
- Code cleanup ...
- One feature idea is to be able to provide a mapping file to map categories to OsmAnd categories, icons, and colors, right now this is hardcoded

# Feedback

I'm happy to take feedback and will continue to work on this tool since the combination of iOverlander data and OsmAnd has worked well for me on a recent trip. Please file feature requests or bug reports as needed.
