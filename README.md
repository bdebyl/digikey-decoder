# DigiKey Decoder
Simple golang app to decode DigiKey part numbers from the 2D Datamatrix Barcode
on DigiKey labels. Spits out part numbers in `stdout` from a webcamera stream.

## TODO
Right now the webcam controls (exposure, brightness, focus, etc.) are manually
tuned for my personal, old, web camera. No flags currently set up, just forcibly
using `/dev/video0`