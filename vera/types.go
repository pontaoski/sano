package vera

import (
	"fmt"
	"image"
	"image/color"
	"io"
)

type BPP int

const (
	BPP1 BPP = 1
	BPP2 BPP = 2
	BPP4 BPP = 4
	BPP8 BPP = 8
)

func ToBPP(u int) (BPP, bool) {
	if u == 1 || u == 2 || u == 4 || u == 8 {
		return BPP(u), true
	}
	return BPP(0), false
}

func (b BPP) MaxColors() int {
	return 1 << b
}

type TileSize int

const (
	Eight     TileSize = 8
	Sixteen   TileSize = 16
	ThirtyTwo TileSize = 32
	SixtyFour TileSize = 64
)

func ToTileSize(u int) (TileSize, bool) {
	if u == 8 || u == 16 || u == 32 || u == 64 {
		return TileSize(u), true
	}
	return TileSize(0), false
}

func ExportBitmap(input image.PalettedImage, bpp BPP, output io.Writer) error {
	_, ok := input.ColorModel().(color.Palette)
	if !ok {
		return fmt.Errorf("image colour model is not indexed colour, but is instead %v", input.ColorModel())
	}
	bounds := input.Bounds()

	currentByte := byte(0)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixelNumber := ((bounds.Max.X - bounds.Min.X) * y) + x
			colour := input.ColorIndexAt(x, y)

			switch bpp {
			case BPP1:
				switch pixelNumber % 8 {
				case 0:
					currentByte = colour << 7
				case 1:
					currentByte |= colour << 6
				case 2:
					currentByte |= colour << 5
				case 3:
					currentByte |= colour << 4
				case 4:
					currentByte |= colour << 3
				case 5:
					currentByte |= colour << 2
				case 6:
					currentByte |= colour << 1
				case 7:
					currentByte |= colour
					_, err := output.Write([]byte{currentByte})
					if err != nil {
						return fmt.Errorf("failed to write byte to bitmap output: %w", err)
					}
				}
			case BPP2:
				switch pixelNumber % 4 {
				case 0:
					currentByte = colour << 6
				case 1:
					currentByte |= colour << 4
				case 2:
					currentByte |= colour << 2
				case 3:
					currentByte |= colour
					_, err := output.Write([]byte{currentByte})
					if err != nil {
						return fmt.Errorf("failed to write byte to bitmap output: %w", err)
					}
				}
			case BPP4:
				switch pixelNumber % 2 {
				case 0:
					currentByte = colour << 4
				case 1:
					currentByte |= colour
					_, err := output.Write([]byte{currentByte})
					if err != nil {
						return fmt.Errorf("failed to write byte to bitmap output: %w", err)
					}
				}
			case BPP8:
				_, err := output.Write([]byte{colour})
				if err != nil {
					return fmt.Errorf("failed to write byte to bitmap output: %w", err)
				}
			}
		}
	}

	return nil
}

func ExportTile(input image.PalettedImage, bpp BPP, tileWidth, tileHeight TileSize, output io.Writer) error {
	_, ok := input.ColorModel().(color.Palette)
	if !ok {
		return fmt.Errorf("image colour model is not indexed colour, but is instead %v", input.ColorModel())
	}

	bounds := input.Bounds()
	imageWidth := (bounds.Max.X - bounds.Min.X)
	imageHeight := (bounds.Max.Y - bounds.Min.Y)
	if imageWidth%int(tileWidth) != 0 {
		return fmt.Errorf("expected image width to be a multiple of %d, but it is %d", int(tileHeight), bounds.Max.X-bounds.Min.X)
	}
	if imageHeight%int(tileHeight) != 0 {
		return fmt.Errorf("expected image height to be a multiple of %d, but it is %d", int(tileHeight), bounds.Max.Y-bounds.Min.Y)
	}

	tilesX := (bounds.Max.X - bounds.Min.X) / int(tileWidth)
	tilesY := (bounds.Max.Y - bounds.Min.Y) / int(tileHeight)

	currentByte := byte(0)

	// write a header for loading
	_, err := output.Write([]byte{0, 0})
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// loop through tiles
	for tileY := 0; tileY < tilesY; tileY++ {
		for tileX := 0; tileX < tilesX; tileX++ {
			xOffset := tileX * int(tileWidth)

			// loop through a single tile
			for dy := 0; dy < int(tileHeight); dy++ {
				for dx := 0; dx < int(tileWidth); dx++ {
					ox := bounds.Min.X + xOffset + dx
					oy := bounds.Min.Y + (tileY*int(tileHeight) + dy)
					pixelNumber := ((dy + tileY) * imageWidth) + xOffset + dx

					colour := input.ColorIndexAt(ox, oy)

					switch bpp {
					case BPP1:
						switch pixelNumber % 8 {
						case 0:
							currentByte = colour << 7
						case 1:
							currentByte |= colour << 6
						case 2:
							currentByte |= colour << 5
						case 3:
							currentByte |= colour << 4
						case 4:
							currentByte |= colour << 3
						case 5:
							currentByte |= colour << 2
						case 6:
							currentByte |= colour << 1
						case 7:
							currentByte |= colour
							_, err := output.Write([]byte{currentByte})
							if err != nil {
								return fmt.Errorf("failed to write byte to tile output: %w", err)
							}
						}
					case BPP2:
						switch pixelNumber % 4 {
						case 0:
							currentByte = colour << 6
						case 1:
							currentByte |= colour << 4
						case 2:
							currentByte |= colour << 2
						case 3:
							currentByte |= colour
							_, err := output.Write([]byte{currentByte})
							if err != nil {
								return fmt.Errorf("failed to write byte to tile output: %w", err)
							}
						}
					case BPP4:
						switch pixelNumber % 2 {
						case 0:
							currentByte = colour << 4
						case 1:
							currentByte |= colour
							_, err := output.Write([]byte{currentByte})
							if err != nil {
								return fmt.Errorf("failed to write byte to tile output: %w", err)
							}
						}
					case BPP8:
						_, err := output.Write([]byte{colour})
						if err != nil {
							return fmt.Errorf("failed to write byte to tile output: %w", err)
						}
					}
				}
			}
		}
	}

	return nil
}
